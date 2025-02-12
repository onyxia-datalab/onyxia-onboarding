package middleware

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/coreos/go-oidc/v3/oidc"
	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
)

type contextKey int

const (
	UserContextKey contextKey = iota
)

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (*oidc.IDToken, error)
}

type oidcAuth struct {
	UsernameClaim string
	Verifier      TokenVerifier
}

type noAuth struct{}

var (
	_ api.SecurityHandler = (*oidcAuth)(nil)
	_ api.SecurityHandler = (*noAuth)(nil)
)

func OidcMiddleware(
	ctx context.Context,
	authenticationMode string,
	issuerUri string,
	clientId string,
	usernameClaim string,
) (api.SecurityHandler, error) {

	if authenticationMode == "none" {
		slog.Info("🚀 Running in No-Auth Mode")
		return &noAuth{}, nil
	}

	oidcProvider, err := oidc.NewProvider(ctx, issuerUri)
	if err != nil {
		slog.Error("❌ Failed to initialize OIDC provider",
			slog.String("issuer", issuerUri),
			slog.Any("error", err),
		)
		return nil, err
	}

	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:                   clientId,
		SkipExpiryCheck:            false,
		SkipIssuerCheck:            false,
		SkipClientIDCheck:          false,
		InsecureSkipSignatureCheck: false,
	})

	slog.Info("🔑 OIDC Middleware Initialized",
		slog.String("issuer", issuerUri),
		slog.String("client_id", clientId),
	)

	return &oidcAuth{
		UsernameClaim: usernameClaim,
		Verifier:      verifier,
	}, nil
}

func (a *oidcAuth) HandleOidc(
	ctx context.Context,
	operation string,
	req api.Oidc,
) (context.Context, error) {
	slog.Info("🔵 Verifying OIDC Token", slog.String("operation", operation))

	token, err := a.Verifier.Verify(ctx, req.Token)
	if err != nil {
		slog.Error("❌ OIDC Token Verification Failed",
			slog.String("operation", operation),
			slog.Any("error", err),
		)
		return ctx, err
	}

	var claims map[string]any
	if err := token.Claims(&claims); err != nil {
		slog.Error("❌ Failed to extract claims from token", slog.Any("error", err))
		return ctx, err
	}

	user, ok := claims[a.UsernameClaim]
	if !ok {
		slog.Error("❌ Missing required claim",
			slog.String("claim", a.UsernameClaim),
		)
		return ctx, fmt.Errorf("missing %q claim", a.UsernameClaim)
	}

	userStr, ok := user.(string)
	if !ok {
		slog.Error("❌ Unexpected claim format",
			slog.String("claim", a.UsernameClaim),
		)
		return ctx, fmt.Errorf("unknown format for claim %q", a.UsernameClaim)
	}

	slog.Info("✅ OIDC Authentication Successful",
		slog.String("user", userStr),
		slog.String("operation", operation),
	)

	return context.WithValue(ctx, UserContextKey, userStr), nil
}

func (n *noAuth) HandleOidc(
	ctx context.Context,
	operation string,
	req api.Oidc,
) (context.Context, error) {
	slog.Warn("⚠️ No-Auth Mode: Skipping authentication.", slog.String("operation", operation))
	return context.WithValue(ctx, UserContextKey, "anonymous"), nil
}
