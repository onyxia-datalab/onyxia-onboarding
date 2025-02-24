package usercontext

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

// MockUserContext is a test implementation of UserContextReader and UserContextWriter
type MockUserContext struct {
	User *domain.User
}

func (m *MockUserContext) GetUser(ctx context.Context) (*domain.User, bool) {
	if m.User == nil {
		return nil, false
	}
	return m.User, true
}

func (m *MockUserContext) GetUsername(ctx context.Context) (string, bool) {
	if m.User != nil {
		return m.User.Username, true
	}
	return "", false
}

func (m *MockUserContext) GetGroups(ctx context.Context) ([]string, bool) {
	if m.User != nil {
		return m.User.Groups, true
	}
	return nil, false
}

func (m *MockUserContext) GetRoles(ctx context.Context) ([]string, bool) {
	if m.User != nil {
		return m.User.Roles, true
	}
	return nil, false
}

func (m *MockUserContext) GetAttributes(ctx context.Context) (map[string]any, bool) {
	if m.User != nil {
		return m.User.Attributes, true
	}
	return nil, false
}

func (m *MockUserContext) WithUser(ctx context.Context, u *domain.User) context.Context {
	m.User = u
	return ctx
}

// NewMockUserContext returns a mock implementation for tests.
func NewMockUserContext(
	user *domain.User,
) (interfaces.UserContextReader, interfaces.UserContextWriter) {
	return &MockUserContext{User: user}, &MockUserContext{User: user}
}
