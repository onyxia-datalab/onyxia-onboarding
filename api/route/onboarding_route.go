package route

import (
	"github.com/onyxia-datalab/onyxia-onboarding/api/controller"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure/kubernetes"
	"github.com/onyxia-datalab/onyxia-onboarding/usecase"
)

func SetupOnboardingController(
	app *bootstrap.Application,
) *controller.OnboardingController {
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

	return controller.NewOnboardingController(onboardingUsecase, app.UserContextReader)

}
