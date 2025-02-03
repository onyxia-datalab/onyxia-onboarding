package bootstrap

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesClient struct {
	Clientset *kubernetes.Clientset
}

func NewKubernetesClient() *KubernetesClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("❌ Failed to get in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Failed to create Kubernetes client: %v", err)
	}

	log.Println("✅ Successfully connected to Kubernetes API")

	return &KubernetesClient{Clientset: clientset}
}
