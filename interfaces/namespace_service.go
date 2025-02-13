package interfaces

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type NamespaceCreationResult string
type QuotaApplicationResult string

const (
	NamespaceCreated       NamespaceCreationResult = "created"
	NamespaceAlreadyExists NamespaceCreationResult = "already_exists"
)

const (
	QuotaCreated   QuotaApplicationResult = "created"
	QuotaUpdated   QuotaApplicationResult = "updated"
	QuotaUnchanged QuotaApplicationResult = "unchanged"
	QuotaIgnored   QuotaApplicationResult = "ignored"
)

type NamespaceService interface {
	CreateNamespace(ctx context.Context, name string) (NamespaceCreationResult, error)
	ApplyResourceQuotas(
		ctx context.Context,
		namespace string,
		quota *domain.Quota,
	) (QuotaApplicationResult, error)
}
