package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesClient struct {
	Clientset *kubernetes.Clientset
}

func NewKubernetesClient() *KubernetesClient {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		log.Println("⚠️  Not running in a Kubernetes cluster, trying local kubeconfig...")

		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("❌ Failed to load kubeconfig: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Failed to create Kubernetes client: %v", err)
	}

	err = checkConnectivity(clientset)
	if err != nil {
		log.Fatalf("❌ Failed to connect to the APIServer : %v", err)
	}

	return &KubernetesClient{Clientset: clientset}
}

func checkConnectivity(clientSet *kubernetes.Clientset) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	version, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("❌ Kubernetes API unreachable: %v", err)
	}

	log.Printf("✅ Kubernetes API is reachable! Server version: %s", version.GitVersion)
	return nil
}
