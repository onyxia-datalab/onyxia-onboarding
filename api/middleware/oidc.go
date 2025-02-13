package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/coreos/go-oidc/v3/oidc"
	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
)

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (*oidc.IDToken, error)
}

type OIDCConfig struct {
	IssuerURI     string
	SkipTLSVerify bool
	ClientID      string
	Audience      string
	UsernameClaim string
	GroupsClaim   string
	RolesClaim    string
}

type oidcAuth struct {
	UsernameClaim     string
	GroupsClaim       string
	RolesClaim        string
	Verifier          TokenVerifier
	Audience          string
	userContextWriter usercontext.UserContextWriter
}

type noAuth struct {
	userContextWriter usercontext.UserContextWriter
}

var (
	_ api.SecurityHandler = (*oidcAuth)(nil)
	_ api.SecurityHandler = (*noAuth)(nil)
)

func OidcMiddleware(
	ctx context.Context,
	authenticationMode string,
	config OIDCConfig,
	userContextWriter usercontext.UserContextWriter,
) (api.SecurityHandler, error) {

	if authenticationMode == "none" {
		slog.Warn("🚀 Running in No-Auth Mode")
		return &noAuth{userContextWriter: userContextWriter}, nil
	}

	oidcProvider, err := oidc.NewProvider(ctx, config.IssuerURI)
	if err != nil {
		slog.Error("❌ Failed to initialize OIDC provider",
			slog.String("issuer", config.IssuerURI),
			slog.Any("error", err),
		)
		return nil, err
	}

	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:                   config.ClientID,
		InsecureSkipSignatureCheck: config.SkipTLSVerify,
	})

	if config.Audience == "" {
		slog.Warn("Skipping audience validation (empty)")

	}

	slog.Info("🔑 OIDC Middleware Initialized",
		slog.String("issuer", config.IssuerURI),
		slog.String("client_id", config.ClientID),
		slog.String("aud", config.Audience),
	)

	return &oidcAuth{
		UsernameClaim:     config.UsernameClaim,
		Verifier:          verifier,
		Audience:          config.Audience,
		GroupsClaim:       config.GroupsClaim,
		RolesClaim:        config.RolesClaim,
		userContextWriter: userContextWriter,
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
		slog.Error(
			"❌ OIDC Token Verification Failed",
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

	// ✅ Validate audience
	if err := a.validateAudience(claims); err != nil {
		return ctx, err
	}

	// ✅ Extract user
	userStr, err := a.extractClaim(claims, a.UsernameClaim)
	if err != nil {
		return ctx, err
	}

	groups := a.extractStringArray(claims, a.GroupsClaim)
	roles := a.extractStringArray(claims, a.RolesClaim)

	slog.Info("✅ OIDC Authentication Successful",
		slog.String("user", userStr),
		slog.String("operation", operation),
		slog.Any("groups", groups),
		slog.Any("roles", roles),
	)

	ctx = a.userContextWriter.WithUser(ctx, userStr)
	ctx = a.userContextWriter.WithGroups(ctx, groups)
	ctx = a.userContextWriter.WithRoles(ctx, roles)

	return ctx, nil
}

func (a *oidcAuth) validateAudience(claims map[string]any) error {
	if a.Audience == "" {
		return nil
	}

	aud, exists := claims["aud"]
	if !exists {
		slog.Error("❌ Missing audience claim")
		return fmt.Errorf("missing audience claim")
	}

	switch v := aud.(type) {
	case string:
		if v != a.Audience {
			slog.Error("❌ Invalid audience", slog.String("expected", a.Audience), slog.String("got", v))
			return fmt.Errorf("invalid audience: expected %q, got %q", a.Audience, v)
		}
	case []string:
		valid := slices.Contains(v, a.Audience)

		if !valid {
			slog.Error("❌ Invalid audience", slog.String("expected", a.Audience), slog.Any("got", v))
			return fmt.Errorf("invalid audience: expected %q, got %v", a.Audience, v)
		}
	default:
		slog.Error("❌ Unexpected audience format", slog.Any("aud", v))
		return fmt.Errorf("invalid audience format")
	}

	return nil
}

func (a *oidcAuth) extractClaim(claims map[string]any, claimName string) (string, error) {
	value, ok := claims[claimName]
	if !ok {
		slog.Error("❌ Missing required claim", slog.String("claim", claimName))
		return "", fmt.Errorf("missing %q claim", claimName)
	}

	strValue, ok := value.(string)
	if !ok {
		slog.Error("❌ Unexpected claim format", slog.String("claim", claimName))
		return "", fmt.Errorf("unknown format for claim %q", claimName)
	}

	return strValue, nil
}

func (a *oidcAuth) extractStringArray(claims map[string]any, claimName string) []string {
	if claimName == "" {
		return nil
	}

	value, exists := claims[claimName]
	if !exists {
		return nil
	}

	result, ok := value.([]string)
	if !ok {
		return nil
	}

	return result
}

func (n *noAuth) HandleOidc(
	ctx context.Context,
	operation string,
	req api.Oidc,
) (context.Context, error) {

	ctx = n.userContextWriter.WithUser(ctx, "anonymous")
	ctx = n.userContextWriter.WithGroups(ctx, []string{})
	ctx = n.userContextWriter.WithRoles(ctx, []string{})

	return ctx, nil
}
