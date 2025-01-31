// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/ogenerrors"
)

// SecurityHandler is handler for security parameters.
type SecurityHandler interface {
	// HandleOidc handles oidc security.
	HandleOidc(ctx context.Context, operationName OperationName, t Oidc) (context.Context, error)
}

func findAuthorization(h http.Header, prefix string) (string, bool) {
	v, ok := h["Authorization"]
	if !ok {
		return "", false
	}
	for _, vv := range v {
		scheme, value, ok := strings.Cut(vv, " ")
		if !ok || !strings.EqualFold(scheme, prefix) {
			continue
		}
		return value, true
	}
	return "", false
}

var oauth2ScopesOidc = map[string][]string{
	OnboardOperation: {},
}

func (s *Server) securityOidc(ctx context.Context, operationName OperationName, req *http.Request) (context.Context, bool, error) {
	var t Oidc
	token, ok := findAuthorization(req.Header, "Bearer")
	if !ok {
		return ctx, false, nil
	}
	t.Token = token
	t.Scopes = oauth2ScopesOidc[operationName]
	rctx, err := s.sec.HandleOidc(ctx, operationName, t)
	if errors.Is(err, ogenerrors.ErrSkipServerSecurity) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return rctx, true, err
}

// SecuritySource is provider of security values (tokens, passwords, etc.).
type SecuritySource interface {
	// Oidc provides oidc security value.
	Oidc(ctx context.Context, operationName OperationName) (Oidc, error)
}

func (s *Client) securityOidc(ctx context.Context, operationName OperationName, req *http.Request) error {
	t, err := s.sec.Oidc(ctx, operationName)
	if err != nil {
		return errors.Wrap(err, "security source \"Oidc\"")
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	return nil
}
