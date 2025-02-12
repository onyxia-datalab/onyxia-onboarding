package middleware

import (
	"context"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/stretchr/testify/assert"
)

// ✅ Test context value retrieval (GetUser, GetGroups, GetRoles)
func TestContextValues(t *testing.T) {
	ctx := context.Background()
	expectedUser := "test-user"
	expectedGroups := []string{"group1", "group2"}
	expectedRoles := []string{"admin", "editor"}

	ctx = context.WithValue(ctx, userContextKey, expectedUser)
	ctx = context.WithValue(ctx, groupsContextKey, expectedGroups)
	ctx = context.WithValue(ctx, rolesContextKey, expectedRoles)

	// Test GetUser
	user, ok := GetUser(ctx)
	assert.True(t, ok, "Expected user to be found")
	assert.Equal(t, expectedUser, user, "Expected user does not match")

	// Test GetGroups
	groups, ok := GetGroups(ctx)
	assert.True(t, ok, "Expected groups to be found")
	assert.Equal(t, expectedGroups, groups, "Expected groups do not match")

	// Test GetRoles
	roles, ok := GetRoles(ctx)
	assert.True(t, ok, "Expected roles to be found")
	assert.Equal(t, expectedRoles, roles, "Expected roles do not match")
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
			map[string]any{"aud": []any{"service1", "onyxia-onboarding"}},
			false,
		},
		{"Missing audience", map[string]any{}, true},
		{"Invalid string audience", map[string]any{"aud": "wrong-audience"}, true},
		{"Invalid array audience", map[string]any{"aud": []any{"service1", "other-service"}}, true},
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

// ✅ Test NoAuth mode (should always return anonymous user)
func TestNoAuth(t *testing.T) {
	noAuthHandler := &noAuth{}
	req := api.Oidc{Token: "ignored-token"}
	ctx := context.Background()

	ctx, err := noAuthHandler.HandleOidc(ctx, "test-operation", req)
	assert.NoError(t, err, "Expected no error in NoAuth mode")

	// Validate anonymous user
	user, ok := GetUser(ctx)
	assert.True(t, ok, "Expected user in context")
	assert.Equal(t, "anonymous", user)

	groups, ok := GetGroups(ctx)
	assert.True(t, ok, "Expected groups in context")
	assert.Empty(t, groups, "Expected empty groups in NoAuth mode")

	roles, ok := GetRoles(ctx)
	assert.True(t, ok, "Expected roles in context")
	assert.Empty(t, roles, "Expected empty roles in NoAuth mode")
}
