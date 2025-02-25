package logging

import (
	"context"
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

type userContextHandler struct {
	handler       slog.Handler
	userCtxReader interfaces.UserContextReader
}

func (h *userContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *userContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if user, ok := h.userCtxReader.GetUsername(ctx); ok {
		record.AddAttrs(slog.String("username", user))
	}

	if groups, ok := h.userCtxReader.GetGroups(ctx); ok {
		record.AddAttrs(slog.Any("groups", groups))
	}

	if roles, ok := h.userCtxReader.GetRoles(ctx); ok {
		record.AddAttrs(slog.Any("roles", roles))
	}

	return h.handler.Handle(ctx, record)
}

func (h *userContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &userContextHandler{handler: h.handler.WithAttrs(attrs), userCtxReader: h.userCtxReader}
}

func (h *userContextHandler) WithGroup(name string) slog.Handler {
	return &userContextHandler{handler: h.handler.WithGroup(name), userCtxReader: h.userCtxReader}
}

func NewUserContextLogger(
	baseHandler slog.Handler,
	userCtx interfaces.UserContextReader,
) slog.Handler {
	return &userContextHandler{handler: baseHandler, userCtxReader: userCtx}
}
