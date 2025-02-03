package middleware

import (
	"context"
	"fmt"
	"log"

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
	Verifier      TokenVerifier // ✅ Uses an interface now!
}

type noAuth struct{}

var (
	_ api.SecurityHandler = (*oidcAuth)(nil)
	_ api.SecurityHandler = (*noAuth)(nil)
)

func OidcMiddleware(ctx context.Context, authenticationMode string, issuerUri string, clientId string, usernameClaim string) (api.SecurityHandler, error) {

	if authenticationMode == "none" {
		log.Println("🚀 Running in No-Auth Mode")
		return &noAuth{}, nil
	}

	oidcProvider, err := oidc.NewProvider(ctx, issuerUri)
	if err != nil {
		return nil, err
	}

	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:                   clientId,
		SkipExpiryCheck:            false,
		SkipIssuerCheck:            false,
		SkipClientIDCheck:          false,
		InsecureSkipSignatureCheck: false,
	})

	return &oidcAuth{
		UsernameClaim: usernameClaim,
		Verifier:      verifier,
	}, nil
}

func (a *oidcAuth) HandleOidc(ctx context.Context, operation string, req api.Oidc) (context.Context, error) {
	token, err := a.Verifier.Verify(ctx, req.Token)
	if err != nil {
		return ctx, err
	}

	var claims map[string]any
	if err := token.Claims(&claims); err != nil {
		return ctx, err
	}

	user, ok := claims[a.UsernameClaim]
	if !ok {
		return ctx, fmt.Errorf("missing %q claim", a.UsernameClaim)
	}

	user, ok = user.(string)
	if !ok {
		return ctx, fmt.Errorf("unknown format for claim %q", a.UsernameClaim)
	}

	return context.WithValue(ctx, UserContextKey, user), nil
}

func (n *noAuth) HandleOidc(ctx context.Context, operation string, req api.Oidc) (context.Context, error) {
	log.Println("⚠️ No-Auth Mode: Skipping authentication.")
	return context.WithValue(ctx, UserContextKey, "anonymous"), nil
}
