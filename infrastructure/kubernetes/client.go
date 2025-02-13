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

func NewKubernetesClient() *KubernetesClient {
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
			slog.Error("❌ Failed to load kubeconfig", slog.Any("error", err))
			panic(err) // Still exit on failure
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("❌ Failed to create Kubernetes client", slog.Any("error", err))
		panic(err) // Still exit on failure
	}

	err = checkConnectivity(clientset)
	if err != nil {
		slog.Error("❌ Failed to connect to the APIServer", slog.Any("error", err))
		panic(err) // Still exit on failure
	}

	return &KubernetesClient{Clientset: clientset}
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
