package services

import (
	"k8s.io/client-go/kubernetes"
)

var (
	// K8sClient holds the connection to the Kubernetes API
	K8sClient *kubernetes.Clientset
)

func GetK8sClient() *kubernetes.Clientset {
	return K8sClient
}
