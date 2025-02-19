package middleware

import (
	"context"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserContextWriter struct {
	mock.Mock
}

func (m *MockUserContextWriter) WithUser(ctx context.Context, u *domain.User) context.Context {
	args := m.Called(ctx, u)
	return args.Get(0).(context.Context)
}

func TestValidateAudience(t *testing.T) {
	auth := &oidcAuth{Audience: "onyxia-onboarding"}

	tests := []struct {
		name      string
		claims    map[string]any
		expectErr bool
	}{
		{"Valid string audience", map[string]any{"aud": "onyxia-onboarding"}, false},
		{
			"Valid array audience",
			map[string]any{"aud": []string{"service1", "onyxia-onboarding"}},
			false,
		},
		{"Missing audience", map[string]any{}, true},
		{"Invalid string audience", map[string]any{"aud": "wrong-audience"}, true},
		{
			"Invalid array audience",
			map[string]any{"aud": []string{"service1", "other-service"}},
			true,
		},
		{"Unexpected format", map[string]any{"aud": 123}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.validateAudience(tt.claims)
			if tt.expectErr {
				assert.Error(t, err, "Expected error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestExtractClaim(t *testing.T) {
	auth := &oidcAuth{}

	tests := []struct {
		name      string
		claims    map[string]any
		claimName string
		expected  string
		expectErr bool
	}{
		{"Valid claim", map[string]any{"username": "test-user"}, "username", "test-user", false},
		{"Missing claim", map[string]any{}, "username", "", true},
		{"Wrong format", map[string]any{"username": 123}, "username", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := auth.extractClaim(tt.claims, tt.claimName)
			if tt.expectErr {
				assert.Error(t, err, "Expected error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tt.expected, value)
			}
		})
	}
}

// ✅ Test extractStringArray()
func TestExtractStringArray(t *testing.T) {
	auth := &oidcAuth{}

	tests := []struct {
		name      string
		claims    map[string]any
		claimName string
		expected  []string
	}{
		{
			"Valid array",
			map[string]any{"groups": []any{"group1", "group2"}},
			"groups",
			[]string{"group1", "group2"},
		},
		{"Missing claim", map[string]any{}, "groups", nil},
		{"Wrong format", map[string]any{"groups": "not-an-array"}, "groups", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := auth.extractStringArray(tt.claims, tt.claimName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoAuth(t *testing.T) {
	mockWriter := new(MockUserContextWriter)

	// ✅ Expect WithUser to be called with "anonymous" user, empty groups, and roles
	expectedUser := &domain.User{
		Username: "anonymous",
		Groups:   []string{},
		Roles:    []string{},
	}
	mockWriter.On("WithUser", mock.Anything, expectedUser).Return(context.Background())

	noAuthHandler := &noAuth{userContextWriter: mockWriter}
	req := api.Oidc{Token: "ignored-token"}
	ctx := context.Background()

	_, err := noAuthHandler.HandleOidc(ctx, "test-operation", req)
	assert.NoError(t, err, "Expected no error in NoAuth mode")

	// ✅ Only check `WithUser`, since groups and roles are included inside it
	mockWriter.AssertCalled(t, "WithUser", mock.Anything, expectedUser)
}
