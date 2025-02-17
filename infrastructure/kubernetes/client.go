package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
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

func NewKubernetesClient() (*KubernetesClient, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		slog.Warn("⚠️  Not running in a Kubernetes cluster, trying local kubeconfig...")

		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	err = checkConnectivity(clientset)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubernetes API server: %w", err)
	}

	return &KubernetesClient{Clientset: clientset}, nil
}

func checkConnectivity(clientSet *kubernetes.Clientset) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	version, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("❌ Kubernetes API unreachable: %w", err)
	}

	slog.Info("✅ Kubernetes API is reachable!",
		slog.String("server_version", version.GitVersion),
	)
	return nil
}
