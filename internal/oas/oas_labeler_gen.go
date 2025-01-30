// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// Labeler is used to allow adding custom attributes to the server request metrics.
type Labeler struct {
	attrs []attribute.KeyValue
}

// Add attributes to the Labeler.
func (l *Labeler) Add(attrs ...attribute.KeyValue) {
	l.attrs = append(l.attrs, attrs...)
}

// AttributeSet returns the attributes added to the Labeler as an attribute.Set.
func (l *Labeler) AttributeSet() attribute.Set {
	return attribute.NewSet(l.attrs...)
}

type labelerContextKey struct{}

// LabelerFromContext retrieves the Labeler from the provided context, if present.
//
// If no Labeler was found in the provided context a new, empty Labeler is returned and the second
// return value is false. In this case it is safe to use the Labeler but any attributes added to
// it will not be used.
func LabelerFromContext(ctx context.Context) (*Labeler, bool) {
	if l, ok := ctx.Value(labelerContextKey{}).(*Labeler); ok {
		return l, true
	}
	return &Labeler{}, false
}

func contextWithLabeler(ctx context.Context, l *Labeler) context.Context {
	return context.WithValue(ctx, labelerContextKey{}, l)
}
