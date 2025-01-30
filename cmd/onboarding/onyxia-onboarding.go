package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/onboarding"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/security"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type config struct {
	OidcIssuerUri     string `env:"OIDC_ISSUER_URI,required,notEmpty"`
	OidcUsernameClaim string `env:"OIDC_USERNAME_CLAIM" envDefault:"preferred_username"`
	OidcClientId      string `env:"OIDC_CLIENT_ID" envDefault:"onyxia-api"`
	NamespacePrefix   string `env:"NAMESPACE_PREFIX" envDefault:"onyxia-user-"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	cfg, err := env.ParseAs[config]()
	if err != nil {
		panic(err)
	}

	auth, err := security.NewOidcAuth(context.Background(),
		cfg.OidcIssuerUri,
		cfg.OidcClientId,
		cfg.OidcUsernameClaim,
	)
	if err != nil {
		panic(err)
	}

	// create the clientset

	clientset, err := getKubernetesClientSet()
	if err != nil {
		panic(err)
	}

	srv, err := oas.NewServer(onboarding.NewKubernetesOnboarder(clientset.CoreV1().Namespaces(), cfg.NamespacePrefix), auth)
	if err != nil {
		panic(err)
	}
	r.Mount("/", http.HandlerFunc(srv.ServeHTTP))

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getKubernetesClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(config)
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
