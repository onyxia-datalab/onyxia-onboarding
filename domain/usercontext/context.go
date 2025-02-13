package usercontext

import "context"

type contextKey struct {
	name string
}

var (
	userKey   = &contextKey{"user"}
	groupsKey = &contextKey{"groups"}
	rolesKey  = &contextKey{"roles"}
)

type UserContextReader interface {
	GetUser(ctx context.Context) (string, bool)
	GetGroups(ctx context.Context) ([]string, bool)
	GetRoles(ctx context.Context) ([]string, bool)
}

type UserContextWriter interface {
	WithUser(ctx context.Context, user string) context.Context
	WithGroups(ctx context.Context, groups []string) context.Context
	WithRoles(ctx context.Context, roles []string) context.Context
}

type userContext struct{}

func (userContext) GetUser(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(userKey).(string)
	return user, ok
}

func (userContext) GetGroups(ctx context.Context) ([]string, bool) {
	groups, ok := ctx.Value(groupsKey).([]string)
	return groups, ok
}

func (userContext) GetRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(rolesKey).([]string)
	return roles, ok
}

func (userContext) WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func (userContext) WithGroups(ctx context.Context, groups []string) context.Context {
	return context.WithValue(ctx, groupsKey, groups)
}

func (userContext) WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesKey, roles)
}

func NewUserContext() (UserContextReader, UserContextWriter) {
	return userContext{}, userContext{}
}
