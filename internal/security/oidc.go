package security

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/oas"
)

type contextKey int

const (
	UserContextKey contextKey = iota
)

type oidcAuth struct {
	UsernameClaim string
	Verifier      *oidc.IDTokenVerifier
}

func NewOidcAuth(ctx context.Context, issuerUri string, clientId string, usernameClaim string) (*oidcAuth, error) {
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

var _ oas.SecurityHandler = (*oidcAuth)(nil)

func (a *oidcAuth) HandleOidc(ctx context.Context, operation string, req oas.Oidc) (context.Context, error) {
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
