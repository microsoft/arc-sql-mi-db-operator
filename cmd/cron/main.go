package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	ms "github.com/pplavetzki/arc-sql-mi/internal"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sqlmi "github.com/pplavetzki/arc-sql-mi/api/v1alpha1"
)

var (
	logger    logr.Logger
	clientset *kubernetes.Clientset
)

type DBResult struct {
	Result *string
	Error  error
}

func getEnvOrFail(key string) string {
	var val string
	if val = os.Getenv(key); val == "" {
		panic(fmt.Errorf("failed to get required env variable: %s", key))
	}
	return val
}

func getDatabaseID(ctx context.Context, msSQL *ms.MSSql, name string, result chan *DBResult) {
	dbID, err := msSQL.FindDatabaseID(ctx, name)
	result <- &DBResult{
		Result: dbID,
		Error:  err,
	}
}

func getDatabaseName(ctx context.Context, msSQL *ms.MSSql, id string, result chan *DBResult) {
	if id != "" {
		dbName, err := msSQL.FindDatabaseName(ctx, id)
		result <- &DBResult{
			Result: dbName,
			Error:  err,
		}
	} else {
		result <- &DBResult{
			Result: nil,
			Error:  nil,
		}
	}
}

func performSync(msSQL *ms.MSSql, db *sqlmi.Database) error {
	dbNameResult := make(chan *DBResult)
	dbIDResult := make(chan *DBResult)

	defer func() {
		if msSQL.DB != nil {
			msSQL.DB.Close()
		}
	}()

	go getDatabaseID(context.TODO(), msSQL, db.Spec.Name, dbIDResult)
	go getDatabaseName(context.TODO(), msSQL, db.Status.DatabaseID, dbNameResult)

	dbNameR := <-dbNameResult
	dbIDR := <-dbIDResult

	if dbNameR.Error != nil {
		panic(fmt.Errorf("failed to query name: %v", dbNameR.Error))
	}
	if dbIDR.Error != nil {
		panic(fmt.Errorf("failed to query db id: %v", dbIDR.Error))
	}
	if dbIDR.Result != nil {
		logger.V(1).Info("found database ID", "database-id", *dbIDR.Result)
	}
	if dbNameR.Result != nil {
		logger.V(1).Info("database name", "database-name", *dbNameR.Result)
	}

	if db.Status.DatabaseID == "" && dbIDR.Result == nil {
		logger.V(0).Info("database does not exist and is not managed by database controller -- serious error", "databaseName", db.Spec.Name)
	} else if db.Status.DatabaseID == "" && dbIDR.Result != nil {
		logger.V(0).Info("database exists on server but not managed by database controller", "databaseName", db.Spec.Name, "guid", *dbIDR.Result)
	} else if db.Status.DatabaseID != "" && (dbIDR.Result != nil && *dbIDR.Result != db.Status.DatabaseID) {
		logger.V(0).Info("database on server does not match what database controller is expecting", "databaseName", db.Spec.Name, "databaseGuid", *dbIDR.Result, "controllerGuid", db.Status.DatabaseID)
	}
	// Now let's sync
	params := &ms.DatabaseConfig{
		DatabaseName:               db.Spec.Name,
		DatabaseID:                 db.Status.DatabaseID,
		CompatibilityLevel:         db.Spec.CompatibilityLevel,
		Collation:                  db.Spec.Collation,
		Parameterization:           db.Spec.Parameterization,
		AllowSnapshotIsolation:     db.Spec.AllowSnapshotIsolation,
		AllowReadCommittedSnapshot: db.Spec.AllowReadCommittedSnapshot,
	}
	syncResponse, err := msSQL.SyncNeeded(context.TODO(), params, ms.Database)
	if err != nil {
		return err
	}
	if syncResponse != nil {
		logger.V(0).Info("database is out-of-sync with database controller", "database", syncResponse)
	} else {
		logger.V(0).Info("database sync not needed")
	}
	return nil
}

func connectionInfo(db *sqlmi.Database) (string, string, error) {
	/*******************************************************************************************************************
	* Quering the defined secret for the database connection
	*******************************************************************************************************************/
	mi, err := ms.QuerySQLManagedInstance(context.TODO(), db.Namespace, db.Spec.SQLManagedInstance)
	if err != nil {
		return "", "", err
	}
	logger.V(1).Info("successfully found managed instance", "sql-managed-instance", db.Spec.SQLManagedInstance)
	if mi.Status.State != "Ready" {
		return "", "", fmt.Errorf("the sql managed instance is not in a `Ready` state, current status is: %v", mi.Status)
	}

	sec, err := clientset.CoreV1().Secrets(mi.Spec.LoginRef.Namespace).Get(context.TODO(), mi.Spec.LoginRef.Name, v1.GetOptions{})
	if err != nil {
		logger.Error(err, "secrets credentials resource not found", "secret-name", mi.Spec.LoginRef.Name)
		return "", "", err
	}

	username := sec.Data["username"]
	password := sec.Data["password"]
	/******************************************************************************************************************/
	return string(username), string(password), nil
}

func main() {
	var config *rest.Config

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("failed starting logger (%v)?", err))
	}
	logger = zapr.NewLogger(zapLog)

	namespace := getEnvOrFail("NAMESPACE")
	databaseCRD := getEnvOrFail("DATABASE_CRD")
	password := getEnvOrFail("DATABASE_PASSWORD")
	user := getEnvOrFail("DATABASE_USER")
	port := getEnvOrFail("DATABASE_PORT")

	if os.Getenv("KUBECONFIG") != "" {
		path := os.Getenv("KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags("", path)
		if err != nil {
			panic(fmt.Errorf("could not load kubeconfig"))
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(fmt.Errorf("failed to create clientset"))
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	crScheme := runtime.NewScheme()
	sqlmi.AddToScheme(crScheme)

	cl, _ := client.New(config, client.Options{
		Scheme: crScheme,
	})

	list := &sqlmi.DatabaseList{}
	err = cl.List(context.TODO(), list, &client.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	db := &sqlmi.Database{}

	cl.Get(context.TODO(), client.ObjectKey{
		Namespace: namespace,
		Name:      databaseCRD,
	}, db)

	var server string
	if os.Getenv("MS_SERVER") != "" {
		server = os.Getenv("MS_SERVER")
	} else {
		server = fmt.Sprintf("%s-p-svc", db.Spec.SQLManagedInstance)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	msSQL := ms.NewMSSql(server, user, password, p)
	performSync(msSQL, db)
}
