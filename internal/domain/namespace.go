package domain

type Annotation struct {
	Enabled bool
	Static  map[string]string
	Dynamic struct {
		LastLoginTimestamp bool
		UserAttributes     []string
	}
}
type Namespace struct {
	NamespacePrefix      string
	GroupNamespacePrefix string
	Annotation           Annotation
	NamespaceLabels      map[string]string
}
