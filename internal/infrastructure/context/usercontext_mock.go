package usercontext

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

// FakeUserContext is a test implementation of UserContextReader and UserContextWriter.
type FakeUserContext struct {
	User *domain.User
}

func (f *FakeUserContext) GetUser(ctx context.Context) (*domain.User, bool) {
	if f.User == nil {
		return nil, false
	}
	return f.User, true
}

func (f *FakeUserContext) GetUsername(ctx context.Context) (string, bool) {
	if f.User != nil {
		return f.User.Username, true
	}
	return "", false
}

func (f *FakeUserContext) GetGroups(ctx context.Context) ([]string, bool) {
	if f.User != nil {
		return f.User.Groups, true
	}
	return nil, false
}

func (f *FakeUserContext) GetRoles(ctx context.Context) ([]string, bool) {
	if f.User != nil {
		return f.User.Roles, true
	}
	return nil, false
}

func (f *FakeUserContext) GetAttributes(ctx context.Context) (map[string]any, bool) {
	if f.User != nil {
		return f.User.Attributes, true
	}
	return nil, false
}

func (f *FakeUserContext) WithUser(ctx context.Context, u *domain.User) context.Context {
	f.User = u
	return ctx
}

// NewFakeUserContext returns a test implementation of UserContextReader and UserContextWriter.
func NewFakeUserContext(
	user *domain.User,
) (interfaces.UserContextReader, interfaces.UserContextWriter) {
	return &FakeUserContext{User: user}, &FakeUserContext{User: user}
}
