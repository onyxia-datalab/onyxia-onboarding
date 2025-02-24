package middleware

import (
	"context"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	usercontext "github.com/onyxia-datalab/onyxia-onboarding/infrastructure/context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateAudience(t *testing.T) {
	tests := []struct {
		name      string
		auth      *oidcAuth // The OIDC Auth config
		claims    map[string]any
		expectErr bool
	}{
		{
			"Empty config audience",
			&oidcAuth{Audience: ""},
			map[string]any{"aud": "onyxia-onboarding"},
			false,
		},
		{
			"Valid string audience",
			&oidcAuth{Audience: "onyxia-onboarding"},
			map[string]any{"aud": "onyxia-onboarding"},
			false,
		},
		{
			"Valid array audience",
			&oidcAuth{Audience: "onyxia-onboarding"},
			map[string]any{"aud": []string{"service1", "onyxia-onboarding"}},
			false,
		},
		{
			"Missing audience in token",
			&oidcAuth{Audience: "onyxia-onboarding"},
			map[string]any{},
			true,
		},
		{
			"Invalid string audience",
			&oidcAuth{Audience: "onyxia-onboarding"},
			map[string]any{"aud": "wrong-audience"},
			true,
		},
		{
			"Invalid array audience",
			&oidcAuth{Audience: "onyxia-onboarding"},
			map[string]any{"aud": []string{"service1", "other-service"}},
			true,
		},
		{
			"Unexpected format",
			&oidcAuth{Audience: "onyxia-onboarding"},
			map[string]any{"aud": 123},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auth.validateAudience(tt.claims)
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
		{"Empty claim name", map[string]any{"groups": []any{"group1"}}, "", nil},
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

func TestOidcMiddleware_NoAuthMode(t *testing.T) {
	// ✅ Use real user context implementation
	userCtxReader, userCtxWriter := usercontext.NewUserContext()

	// ✅ Call OidcMiddleware with "none" mode
	securityHandler, err := OidcMiddleware(
		context.Background(),
		"none",
		OIDCConfig{},
		userCtxWriter,
	)

	// ✅ Assert that no error occurred
	assert.NoError(t, err, "Expected no error when using No-Auth mode")

	// ✅ Assert that the returned security handler is a *noAuth instance
	assert.IsType(t, &noAuth{}, securityHandler, "Expected securityHandler to be of type *noAuth")

	// ✅ Cast to *noAuth and validate that it has the correct userContextWriter
	noAuthHandler, ok := securityHandler.(*noAuth)
	assert.True(t, ok, "Expected securityHandler to be a *noAuth instance")
	assert.Equal(
		t,
		userCtxWriter,
		noAuthHandler.userContextWriter,
		"Expected userContextWriter to be the same as passed",
	)

	// ✅ Simulate an OIDC request and verify NoAuth behavior
	req := api.Oidc{Token: "ignored-token"}
	ctx := context.Background()
	ctx, err = noAuthHandler.HandleOidc(ctx, "test-operation", req)

	assert.NoError(t, err, "Expected no error when handling OIDC request in No-Auth mode")

	// ✅ Assert that the anonymous user was set
	expectedUser := &domain.User{
		Username:   "anonymous",
		Groups:     []string{},
		Roles:      []string{},
		Attributes: map[string]any{},
	}
	user, exists := userCtxReader.GetUser(ctx)
	assert.True(t, exists, "Expected user to exist in context")
	assert.Equal(t, expectedUser, user, "Expected user to be set in context")
}
