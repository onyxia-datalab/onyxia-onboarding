package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/api/controller"
	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure/kubernetes"
	"github.com/onyxia-datalab/onyxia-onboarding/usecase"
)

func SetupOnboardingRoutes(
	app *bootstrap.Application,
	router chi.Router,
	auth api.SecurityHandler,
	userContextReader usercontext.UserContextReader,
) {

	namespaceCreator := kubernetes.NewKubernetesNamespaceService(app.K8sClient.Clientset)

	onboardingUsecase := usecase.NewOnboardingUsecase(
		namespaceCreator,
		app.Env.Service.NamespacePrefix,
		app.Env.Service.GroupNamespacePrefix,
		domain.Quotas{
			Enabled: app.Env.Service.Quotas.Enabled,
			Default: domain.Quota{
				MemoryRequest:           app.Env.Service.Quotas.Default.RequestsMemory,
				CPURequest:              app.Env.Service.Quotas.Default.RequestsCPU,
				MemoryLimit:             app.Env.Service.Quotas.Default.LimitsMemory,
				CPULimit:                app.Env.Service.Quotas.Default.LimitsCPU,
				StorageRequest:          app.Env.Service.Quotas.Default.RequestsStorage,
				MaxPods:                 app.Env.Service.Quotas.Default.CountPods,
				EphemeralStorageRequest: app.Env.Service.Quotas.Default.RequestsEphemeralStorage,
				EphemeralStorageLimit:   app.Env.Service.Quotas.Default.LimitsEphemeralStorage,
				GPURequest:              app.Env.Service.Quotas.Default.RequestsGPU,
				GPULimit:                app.Env.Service.Quotas.Default.LimitsGPU,
			},
			UserEnabled: app.Env.Service.Quotas.UserEnabled,
			User: domain.Quota{
				MemoryRequest:           app.Env.Service.Quotas.User.RequestsMemory,
				CPURequest:              app.Env.Service.Quotas.User.RequestsCPU,
				MemoryLimit:             app.Env.Service.Quotas.User.LimitsMemory,
				CPULimit:                app.Env.Service.Quotas.User.LimitsCPU,
				StorageRequest:          app.Env.Service.Quotas.User.RequestsStorage,
				MaxPods:                 app.Env.Service.Quotas.User.CountPods,
				EphemeralStorageRequest: app.Env.Service.Quotas.User.RequestsEphemeralStorage,
				EphemeralStorageLimit:   app.Env.Service.Quotas.User.LimitsEphemeralStorage,
				GPURequest:              app.Env.Service.Quotas.User.RequestsGPU,
				GPULimit:                app.Env.Service.Quotas.User.LimitsGPU,
			},
			GroupEnabled: app.Env.Service.Quotas.GroupEnabled,
			Group: domain.Quota{
				MemoryRequest:           app.Env.Service.Quotas.Group.RequestsMemory,
				CPURequest:              app.Env.Service.Quotas.Group.RequestsCPU,
				MemoryLimit:             app.Env.Service.Quotas.Group.LimitsMemory,
				CPULimit:                app.Env.Service.Quotas.Group.LimitsCPU,
				StorageRequest:          app.Env.Service.Quotas.Group.RequestsStorage,
				MaxPods:                 app.Env.Service.Quotas.Group.CountPods,
				EphemeralStorageRequest: app.Env.Service.Quotas.Group.RequestsEphemeralStorage,
				EphemeralStorageLimit:   app.Env.Service.Quotas.Group.LimitsEphemeralStorage,
				GPURequest:              app.Env.Service.Quotas.Group.RequestsGPU,
				GPULimit:                app.Env.Service.Quotas.Group.LimitsGPU,
			},
		},
	)

	controller := controller.NewOnboardingController(onboardingUsecase, userContextReader)

	srv, err := api.NewServer(controller, auth)
	if err != nil {
		panic(err)
	}

	router.Mount("/", http.HandlerFunc(srv.ServeHTTP))
}
