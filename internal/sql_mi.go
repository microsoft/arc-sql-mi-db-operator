package internal

import "time"

type SQLManagedInstance struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			ManagementAzureComAPIVersion     string `json:"management.azure.com/apiVersion"`
			ManagementAzureComCorrelationID  string `json:"management.azure.com/correlationId"`
			ManagementAzureComCustomLocation string `json:"management.azure.com/customLocation"`
			ManagementAzureComLocation       string `json:"management.azure.com/location"`
			ManagementAzureComOperationID    string `json:"management.azure.com/operationId"`
			ManagementAzureComResourceID     string `json:"management.azure.com/resourceId"`
			ManagementAzureComTenantID       string `json:"management.azure.com/tenantId"`
			Traceparent                      string `json:"traceparent"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation        int       `json:"generation"`
		Name              string    `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Dev         bool   `json:"dev"`
		LicenseType string `json:"licenseType"`
		LoginRef    struct {
			Kind      string `json:"kind"`
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"loginRef"`
		Replicas   int `json:"replicas"`
		Scheduling struct {
			Default struct {
				Resources struct {
					Limits struct {
						CPU    string `json:"cpu"`
						Memory string `json:"memory"`
					} `json:"limits"`
					Requests struct {
						CPU    string `json:"cpu"`
						Memory string `json:"memory"`
					} `json:"requests"`
				} `json:"resources"`
			} `json:"default"`
		} `json:"scheduling"`
		Services struct {
			Primary struct {
				Type string `json:"type"`
			} `json:"primary"`
		} `json:"services"`
		Storage struct {
			Backups struct {
				Volumes []struct {
					ClassName string `json:"className"`
					Size      string `json:"size"`
				} `json:"volumes"`
			} `json:"backups"`
			Data struct {
				Volumes []struct {
					ClassName string `json:"className"`
					Size      string `json:"size"`
				} `json:"volumes"`
			} `json:"data"`
			Datalogs struct {
				Volumes []struct {
					ClassName string `json:"className"`
					Size      string `json:"size"`
				} `json:"volumes"`
			} `json:"datalogs"`
			Logs struct {
				Volumes []struct {
					ClassName string `json:"className"`
					Size      string `json:"size"`
				} `json:"volumes"`
			} `json:"logs"`
		} `json:"storage"`
		Tier string `json:"tier"`
	} `json:"spec"`
	Status struct {
		AGStatus           string `json:"AGStatus"`
		LogSearchDashboard string `json:"logSearchDashboard"`
		MetricsDashboard   string `json:"metricsDashboard"`
		ObservedGeneration int    `json:"observedGeneration"`
		PrimaryEndpoint    string `json:"primaryEndpoint"`
		ReadyReplicas      string `json:"readyReplicas"`
		SecondaryEndpoint  string `json:"secondaryEndpoint"`
		State              string `json:"state"`
	} `json:"status"`
}
