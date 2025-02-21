package usercontext

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

type contextKey struct {
	name string
}

// Using a struct as a key instead of a plain string helps prevent accidental key conflicts.
var userKey = &contextKey{"user"}

type userContextImpl struct{}

func (userContextImpl) GetUser(ctx context.Context) (*domain.User, bool) {
	u, ok := ctx.Value(userKey).(*domain.User)
	return u, ok
}

func (uc userContextImpl) GetUsername(ctx context.Context) (string, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Username, true
	}
	return "", false
}

func (uc userContextImpl) GetGroups(ctx context.Context) ([]string, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Groups, true
	}
	return nil, false
}

func (uc userContextImpl) GetRoles(ctx context.Context) ([]string, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Roles, true
	}
	return nil, false
}

func (uc userContextImpl) GetAttributes(ctx context.Context) (map[string]any, bool) {
	if user, ok := uc.GetUser(ctx); ok {
		return user.Attributes, true
	}
	return nil, false
}

func (userContextImpl) WithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func NewUserContext() (interfaces.UserContextReader, interfaces.UserContextWriter) {
	return userContextImpl{}, userContextImpl{}
}
