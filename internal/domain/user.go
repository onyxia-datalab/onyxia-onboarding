package domain

type User struct {
	Username   string
	Groups     []string
	Roles      []string
	Attributes map[string]any
}
