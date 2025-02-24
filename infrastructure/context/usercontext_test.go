package usercontext

import (
	"context"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
)

func TestWithUser(t *testing.T) {
	ctx := context.Background()
	_, writer := NewUserContext()

	expectedUser := &domain.User{
		Username:   "john.doe",
		Groups:     []string{"admin", "developer"},
		Roles:      []string{"reader", "writer"},
		Attributes: map[string]any{"key1": "value1", "key2": 42}, // ✅ Added attributes
	}

	// ✅ Store the full user in context
	ctx = writer.WithUser(ctx, expectedUser)

	reader, _ := NewUserContext()

	// ✅ Retrieve the full user
	retrievedUser, ok := reader.GetUser(ctx)
	assert.True(t, ok, "Expected user to be present in context")
	assert.Equal(t, expectedUser, retrievedUser, "Retrieved user should match stored user")

	// ✅ Test retrieving individual attributes
	username, ok := reader.GetUsername(ctx)
	assert.True(t, ok, "Expected username to be present in context")
	assert.Equal(t, expectedUser.Username, username, "Username should match")

	groups, ok := reader.GetGroups(ctx)
	assert.True(t, ok, "Expected groups to be present in context")
	assert.Equal(t, expectedUser.Groups, groups, "Groups should match")

	roles, ok := reader.GetRoles(ctx)
	assert.True(t, ok, "Expected roles to be present in context")
	assert.Equal(t, expectedUser.Roles, roles, "Roles should match")

	attributes, ok := reader.GetAttributes(ctx)
	assert.True(t, ok, "Expected attributes to be present in context")
	assert.NotNil(t, attributes, "Attributes should not be nil")
	assert.Equal(t, expectedUser.Attributes, attributes, "Attributes should match")
}

func TestMissingValues(t *testing.T) {
	ctx := context.Background()
	reader, _ := NewUserContext()

	// ✅ Ensure GetUser returns nil when no user is set
	retrievedUser, ok := reader.GetUser(ctx)
	assert.False(t, ok, "Expected user to be missing in context")
	assert.Nil(t, retrievedUser, "User should be nil")

	// ✅ Ensure GetUsername returns empty string
	username, ok := reader.GetUsername(ctx)
	assert.False(t, ok, "Expected username to be missing in context")
	assert.Equal(t, "", username, "Username should be an empty string")

	// ✅ Ensure GetGroups returns nil
	groups, ok := reader.GetGroups(ctx)
	assert.False(t, ok, "Expected groups to be missing in context")
	assert.Nil(t, groups, "Groups should be nil")

	// ✅ Ensure GetRoles returns nil
	roles, ok := reader.GetRoles(ctx)
	assert.False(t, ok, "Expected roles to be missing in context")
	assert.Nil(t, roles, "Roles should be nil")

	attributes, ok := reader.GetAttributes(ctx)
	assert.False(t, ok, "Expected attributes to be missing in context")
	assert.Nil(t, attributes, "Attributes should be nil when missing")

}
