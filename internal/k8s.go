package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DatabaseKey struct {
	DatabaseID string `json:"database-id"`
}

func QuerySQLManagedInstance(ctx context.Context, namespace, name string) (*SQLManagedInstance, error) {
	_ = log.FromContext(ctx)
	logger := log.Log
	// uri := "http://localhost:8080/api"
	uri := fmt.Sprintf("http://localhost:9090/apis/sql.arcdata.microsoft.com/v1/namespaces/%s/sqlmanagedinstances/%s", namespace, name)
	logger.V(1).Info("uri of k8s api", "uri", uri)

	k8sClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := k8sClient.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get sqlmanagedinstance: %s, error: %s", name, string(body))
	}

	mi := &SQLManagedInstance{}
	jsonErr := json.Unmarshal(body, mi)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return mi, nil
}

func QueryJobPod(ctx context.Context, namespace, name string) (*string, error) {
	_ = log.FromContext(ctx)
	dbreg := regexp.MustCompile(`{"database-id": "(.*)"}`)

	logger := log.Log
	// uri := "http://localhost:8080/api"
	uri := fmt.Sprintf("http://localhost:9090/api/v1/namespaces/%s/pods?labelSelector=job-name=%s", namespace, name)
	logger.V(1).Info("uri of k8s api", "uri", uri)

	k8sClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := k8sClient.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get log: %s, error: %s", name, string(body))
	}
	var pl corev1.PodList
	jsonErr := json.Unmarshal(body, &pl)
	if jsonErr != nil {
		return nil, jsonErr
	}
	p := pl.Items[0]

	uriPod := fmt.Sprintf("http://localhost:9090/api/v1/namespaces/%s/pods/%s/log?pretty", namespace, p.Name)
	reqLog, err := http.NewRequest(http.MethodGet, uriPod, nil)
	if err != nil {
		return nil, err
	}
	resLog, getErr := k8sClient.Do(reqLog)
	if getErr != nil {
		return nil, err
	}
	if resLog.Body != nil {
		defer resLog.Body.Close()
	}
	bodyLog, readErr := ioutil.ReadAll(resLog.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resLog.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get log: %s, error: %s", name, string(body))
	}
	dbbb := dbreg.Find(bodyLog)
	dbbbbb := &DatabaseKey{}
	err = json.Unmarshal(dbbb, dbbbbb)
	if err != nil {
		return nil, err
	}
	return &dbbbbb.DatabaseID, nil
}
