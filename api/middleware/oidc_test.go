package middleware

import (
	"context"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserContextWriter struct {
	mock.Mock
}

func (m *MockUserContextWriter) WithUser(ctx context.Context, user string) context.Context {
	m.Called(ctx, user)
	return ctx
}

func (m *MockUserContextWriter) WithGroups(ctx context.Context, groups []string) context.Context {
	m.Called(ctx, groups)
	return ctx
}

func (m *MockUserContextWriter) WithRoles(ctx context.Context, roles []string) context.Context {
	m.Called(ctx, roles)
	return ctx
}

func TestNoAuth(t *testing.T) {
	mockWriter := new(MockUserContextWriter)
	mockWriter.On("WithUser", mock.Anything, "anonymous").Return(context.Background())
	mockWriter.On("WithGroups", mock.Anything, []string{}).Return(context.Background())
	mockWriter.On("WithRoles", mock.Anything, []string{}).Return(context.Background())

	noAuthHandler := &noAuth{userContextWriter: mockWriter}
	req := api.Oidc{Token: "ignored-token"}
	ctx := context.Background()

	_, err := noAuthHandler.HandleOidc(ctx, "test-operation", req)
	assert.NoError(t, err, "Expected no error in NoAuth mode")

	mockWriter.AssertCalled(t, "WithUser", mock.Anything, "anonymous")
	mockWriter.AssertCalled(t, "WithGroups", mock.Anything, []string{})
	mockWriter.AssertCalled(t, "WithRoles", mock.Anything, []string{})
}
