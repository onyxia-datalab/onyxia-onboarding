package interfaces

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
)

type UserContextReader interface {
	GetUser(ctx context.Context) (*domain.User, bool)
	GetUsername(ctx context.Context) (string, bool)
	GetGroups(ctx context.Context) ([]string, bool)
	GetRoles(ctx context.Context) ([]string, bool)
	GetAttributes(ctx context.Context) (map[string]any, bool)
}

type UserContextWriter interface {
	WithUser(ctx context.Context, user *domain.User) context.Context
}
