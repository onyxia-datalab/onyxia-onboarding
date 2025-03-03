package usercontext

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

type userContextKey struct {
	name string
}

var ctxUserKey = &userContextKey{"user"}

type userContext struct{}

func (userContext) GetUser(ctx context.Context) (*domain.User, bool) {
	u, ok := ctx.Value(ctxUserKey).(*domain.User)
	return u, ok
}

func (uc userContext) GetUsername(ctx context.Context) (string, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Username, true
	}
	return "", false
}

func (uc userContext) GetGroups(ctx context.Context) ([]string, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Groups, true
	}
	return nil, false
}

func (uc userContext) GetRoles(ctx context.Context) ([]string, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Roles, true
	}
	return nil, false
}

func (uc userContext) GetAttributes(ctx context.Context) (map[string]any, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Attributes, true
	}
	return nil, false
}

func (userContext) WithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, ctxUserKey, u)
}

func NewUserContext() (interfaces.UserContextReader, interfaces.UserContextWriter) {
	return userContext{}, userContext{}
}
