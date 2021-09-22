/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CredentialsSecret is the credentials of the secret to use for the sql server login
type CredentialsSecret struct {
	// Name is the Database name.
	Name        string `json:"name"`
	PasswordKey string `json:"passwordKey"`
	UsernameKey string `json:"usernameKey"`
}

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the Database name.
	Name string `json:"name"`
	// Server is the sql server (fqdn/ip addresss)
	Server string `json:"server,omitempty"`
	// CredentialsSecret is the name of the secret to use for the sql server login credentials
	Credentials CredentialsSecret `json:"credentials,omitempty"`
	// Port where Sql Server is listening
	Port int `json:"port,omitempty"`
	// CollationName
	Collation string `json:"collation,omitempty"`
	// AllowSnapshotIsolation
	AllowSnapshotIsolation     bool   `json:"allowSnapshotIsolation,omitempty"`
	AllowReadCommittedSnapshot bool   `json:"allowReadCommittedSnapshot,omitempty"`
	Parameterization           string `json:"parameterization,omitempty"`
	CompatibilityLevel         int    `json:"compatibilityLevel,omitempty"`
	// SQLManagedInstance name of the managed instance to create database in
	// this is used to query for the status of the instance as well as
	// primary endpoint and connection info
	SQLManagedInstance string `json:"sqlManagedInstance"`
	// Schedule how often the database to k8s state should occur in cron format
	Schedule string `json:"schedule,omitempty"`
}

// DatabaseStatus defines the observed state of Database
type DatabaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status"`
	// DatabaseID guid of the database
	DatabaseID string `json:"databaseID,omitempty"`
	// Conditions the array of conditions of the object
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Database ID",type="string",JSONPath=`.status.databaseID`,description="MSSql Database ID"
//+kubebuilder:printcolumn:name="Database Name",type=string,JSONPath=`.spec.name`,description="Name of Database"
//+kubebuilder:printcolumn:name="Database Status",type=string,JSONPath=`.status.status`,description="Status of Database"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Database is the Schema for the databases API
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
