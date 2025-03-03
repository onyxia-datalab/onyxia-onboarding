package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	usercontext "github.com/onyxia-datalab/onyxia-onboarding/internal/infrastructure/context"
	"github.com/stretchr/testify/assert"
)

func TestUserContextHandler_Enabled(t *testing.T) {
	userContextReader, _ := usercontext.NewFakeUserContext(&domain.User{
		Username: "test_user",
		Groups:   []string{"group1", "group2"},
		Roles:    []string{"role1"},
	})

	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(
		&buf,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	) // Set level threshold
	handler := NewUserContextLogger(baseHandler, userContextReader)

	assert.True(
		t,
		handler.Enabled(context.Background(), slog.LevelInfo),
		"LevelInfo should be enabled",
	)
	assert.False(
		t,
		handler.Enabled(context.Background(), slog.LevelDebug),
		"LevelDebug should be disabled",
	)
}

func TestUserContextHandler_Handle(t *testing.T) {
	// Create a fake user
	user := &domain.User{
		Username: "test_user",
		Groups:   []string{"group1", "group2"},
		Roles:    []string{"role1"},
	}

	// Get the fake user context
	userContextReader, _ := usercontext.NewFakeUserContext(user)

	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewUserContextLogger(baseHandler, userContextReader)

	// Simulate a log record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)

	err := handler.Handle(context.Background(), record)

	assert.NoError(t, err, "Handle should not return an error")

	var logEntry struct {
		Level    string   `json:"level"`
		Message  string   `json:"msg"`
		Username string   `json:"username"`
		Groups   []string `json:"groups"`
		Roles    []string `json:"roles"`
	}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err, "Log output should be valid JSON")

	assert.Equal(t, "test_user", logEntry.Username, "Username should match")
	assert.ElementsMatch(t, []string{"group1", "group2"}, logEntry.Groups, "Groups should match")
	assert.ElementsMatch(t, []string{"role1"}, logEntry.Roles, "Roles should match")
}

func TestUserContextHandler_WithAttrs(t *testing.T) {
	user := &domain.User{
		Username: "test_user",
		Groups:   []string{"group1", "group2"},
		Roles:    []string{"role1"},
	}

	userContextReader, _ := usercontext.NewFakeUserContext(user)

	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewUserContextLogger(baseHandler, userContextReader)

	// Add custom attributes
	attrs := []slog.Attr{slog.String("custom_key", "custom_value")}
	handlerWithAttrs := handler.WithAttrs(attrs)

	// Simulate a log record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	err := handlerWithAttrs.Handle(context.Background(), record)

	assert.NoError(t, err, "Handle should not return an error")

	var logEntry map[string]any
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err, "Log output should be valid JSON")

	assert.Equal(
		t,
		"custom_value",
		logEntry["custom_key"],
		"Custom attribute should be present in the log",
	)
}

func TestUserContextHandler_WithGroup(t *testing.T) {
	user := &domain.User{
		Username: "test_user",
		Groups:   []string{"group1", "group2"},
		Roles:    []string{"role1"},
	}

	userContextReader, _ := usercontext.NewFakeUserContext(user)

	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewUserContextLogger(baseHandler, userContextReader)

	// Apply grouping
	groupedHandler := handler.WithGroup("context")

	// Simulate a log record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	err := groupedHandler.Handle(context.Background(), record)

	assert.NoError(t, err, "Handle should not return an error")

	var logEntry map[string]any
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err, "Log output should be valid JSON")

	contextFields, exists := logEntry["context"].(map[string]any)
	assert.True(t, exists, "Grouped fields should exist under 'context'")
	assert.Equal(
		t,
		"test_user",
		contextFields["username"],
		"Username should be inside the 'context' group",
	)
}
