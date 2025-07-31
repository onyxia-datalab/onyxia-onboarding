package interfaces

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
)

type NamespaceCreationResult string
type QuotaApplicationResult string

const (
	NamespaceCreated            NamespaceCreationResult = "created"
	NamespaceAlreadyExists      NamespaceCreationResult = "already_exists"
	NamespaceAnnotationsUpdated NamespaceCreationResult = "annotations_updated"
)

const (
	QuotaCreated   QuotaApplicationResult = "created"
	QuotaUpdated   QuotaApplicationResult = "updated"
	QuotaUnchanged QuotaApplicationResult = "unchanged"
	QuotaIgnored   QuotaApplicationResult = "ignored"
)

type NamespaceService interface {
	CreateNamespace(
		ctx context.Context,
		name string,
		annotations map[string]string,
		labels map[string]string,
	) (NamespaceCreationResult, error)
	ApplyResourceQuotas(
		ctx context.Context,
		namespace string,
		quota *domain.Quota,
	) (QuotaApplicationResult, error)
}
