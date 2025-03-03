package route

import (
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/bootstrap"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestConvertBootstrapQuotaToDomain(t *testing.T) {
	bootstrapQuota := bootstrap.Quota{
		RequestsMemory:           "512Mi",
		RequestsCPU:              "250m",
		LimitsMemory:             "1Gi",
		LimitsCPU:                "500m",
		RequestsStorage:          "10Gi",
		CountPods:                "20",
		RequestsEphemeralStorage: "5Gi",
		LimitsEphemeralStorage:   "10Gi",
		RequestsGPU:              "1",
		LimitsGPU:                "2",
	}

	expectedDomainQuota := domain.Quota{
		MemoryRequest:           bootstrapQuota.RequestsMemory,
		CPURequest:              bootstrapQuota.RequestsCPU,
		MemoryLimit:             bootstrapQuota.LimitsMemory,
		CPULimit:                bootstrapQuota.LimitsCPU,
		StorageRequest:          bootstrapQuota.RequestsStorage,
		MaxPods:                 bootstrapQuota.CountPods,
		EphemeralStorageRequest: bootstrapQuota.RequestsEphemeralStorage,
		EphemeralStorageLimit:   bootstrapQuota.LimitsEphemeralStorage,
		GPURequest:              bootstrapQuota.RequestsGPU,
		GPULimit:                bootstrapQuota.LimitsGPU,
	}

	result := convertBootstrapQuotaToDomain(bootstrapQuota)

	assert.Equal(t, expectedDomainQuota, result, "Quota conversion should correctly map all fields")
}
