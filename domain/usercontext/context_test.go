package usercontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithUser(t *testing.T) {
	ctx := context.Background()
	_, writer := NewUserContext()
	ctx = writer.WithUser(ctx, "john.doe")

	reader, _ := NewUserContext()
	user, ok := reader.GetUser(ctx)
	assert.True(t, ok, "Expected user to be present in context")
	assert.Equal(t, "john.doe", user, "User should be 'john.doe'")
}

func TestWithGroups(t *testing.T) {
	ctx := context.Background()
	_, writer := NewUserContext()
	groups := []string{"admin", "developer"}
	ctx = writer.WithGroups(ctx, groups)

	reader, _ := NewUserContext()
	result, ok := reader.GetGroups(ctx)
	assert.True(t, ok, "Expected groups to be present in context")
	assert.Equal(t, groups, result, "Groups should match the input")
}

func TestWithRoles(t *testing.T) {
	ctx := context.Background()
	_, writer := NewUserContext()
	roles := []string{"reader", "writer"}
	ctx = writer.WithRoles(ctx, roles)

	reader, _ := NewUserContext()
	result, ok := reader.GetRoles(ctx)
	assert.True(t, ok, "Expected roles to be present in context")
	assert.Equal(t, roles, result, "Roles should match the input")
}

func TestMissingValues(t *testing.T) {
	ctx := context.Background()
	reader, _ := NewUserContext()

	user, ok := reader.GetUser(ctx)
	assert.False(t, ok, "Expected user to be missing in context")
	assert.Equal(t, "", user, "User should be an empty string")

	groups, ok := reader.GetGroups(ctx)
	assert.False(t, ok, "Expected groups to be missing in context")
	assert.Nil(t, groups, "Groups should be nil")

	roles, ok := reader.GetRoles(ctx)
	assert.False(t, ok, "Expected roles to be missing in context")
	assert.Nil(t, roles, "Roles should be nil")
}
