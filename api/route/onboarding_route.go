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

	envQuotas := app.Env.Onboarding.Quotas

	rolesDomainQuotas := func() map[string]domain.Quota {
		result := make(map[string]domain.Quota)
		for key, q := range envQuotas.Roles {
			result[key] = convertBootstrapQuotaToDomain(q)
		}
		return result
	}()

	onboardingUsecase := usecase.NewOnboardingUsecase(
		namespaceCreator,
		app.Env.Onboarding.NamespacePrefix,
		app.Env.Onboarding.GroupNamespacePrefix,
		domain.Quotas{
			Enabled:      envQuotas.Enabled,
			Default:      convertBootstrapQuotaToDomain(envQuotas.Default),
			UserEnabled:  envQuotas.UserEnabled,
			User:         convertBootstrapQuotaToDomain(envQuotas.User),
			Roles:        rolesDomainQuotas,
			GroupEnabled: envQuotas.GroupEnabled,
			Group:        convertBootstrapQuotaToDomain(envQuotas.Group),
		},
		app.Env.Service.NamespaceAnnotations,
	)

	return controller.NewOnboardingController(onboardingUsecase, app.UserContextReader)

}

func convertBootstrapQuotaToDomain(q bootstrap.Quota) domain.Quota {
	return domain.Quota{
		MemoryRequest:           q.RequestsMemory,
		CPURequest:              q.RequestsCPU,
		MemoryLimit:             q.LimitsMemory,
		CPULimit:                q.LimitsCPU,
		StorageRequest:          q.RequestsStorage,
		MaxPods:                 q.CountPods,
		EphemeralStorageRequest: q.RequestsEphemeralStorage,
		EphemeralStorageLimit:   q.LimitsEphemeralStorage,
		GPURequest:              q.RequestsGPU,
		GPULimit:                q.LimitsGPU,
	}
}
