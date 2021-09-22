package v1alpha1

import corev1 "k8s.io/api/core/v1"

// SecretsFromSource key/value pair from secret data
type SecretsFromSource struct {
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}
