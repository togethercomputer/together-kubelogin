package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/togethercomputer/together-kubelogin/pkg/infrastructure/clock"
	"github.com/togethercomputer/together-kubelogin/pkg/infrastructure/logger"
	"github.com/togethercomputer/together-kubelogin/pkg/oidc"
	"github.com/togethercomputer/together-kubelogin/pkg/pkce"
	"github.com/int128/oauth2dev"
	"golang.org/x/oauth2"
)

type Interface interface {
	GetAuthCodeURL(in AuthCodeURLInput) string
	ExchangeAuthCode(ctx context.Context, in ExchangeAuthCodeInput) (*oidc.TokenSet, error)
	GetTokenByAuthCode(ctx context.Context, in GetTokenByAuthCodeInput, localServerReadyChan chan<- string) (*oidc.TokenSet, error)
	NegotiatedPKCEMethod() pkce.Method
	GetTokenByROPC(ctx context.Context, username, password string) (*oidc.TokenSet, error)
	GetTokenByClientCredentials(ctx context.Context, in GetTokenByClientCredentialsInput) (*oidc.TokenSet, error)
	GetDeviceAuthorization(ctx context.Context) (*oauth2dev.AuthorizationResponse, error)
	ExchangeDeviceCode(ctx context.Context, authResponse *oauth2dev.AuthorizationResponse) (*oidc.TokenSet, error)
	Refresh(ctx context.Context, refreshToken string) (*oidc.TokenSet, error)
}

type client struct {
	httpClient           *http.Client
	provider             *gooidc.Provider
	oauth2Config         oauth2.Config
	clock                clock.Interface
	logger               logger.Interface
	negotiatedPKCEMethod pkce.Method
	useAccessToken       bool
	expectedAudience     []string
}

func (c *client) wrapContext(ctx context.Context) context.Context {
	if c.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)
	}
	return ctx
}

// Refresh sends a refresh token request and returns a token set.
func (c *client) Refresh(ctx context.Context, refreshToken string) (*oidc.TokenSet, error) {
	ctx = c.wrapContext(ctx)
	currentToken := &oauth2.Token{
		Expiry:       time.Now(),
		RefreshToken: refreshToken,
	}
	source := c.oauth2Config.TokenSource(ctx, currentToken)
	token, err := source.Token()
	if err != nil {
		return nil, fmt.Errorf("could not refresh the token: %w", err)
	}
	return c.verifyToken(ctx, token, "")
}

// verifyToken verifies the token with the certificates of the provider and the nonce.
// If the nonce is an empty string, it does not verify the nonce.
func (c *client) verifyToken(ctx context.Context, token *oauth2.Token, nonce string) (*oidc.TokenSet, error) {
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("id_token is missing in the token response: %#v", token)
	}
	verifier := c.provider.Verifier(&gooidc.Config{ClientID: c.oauth2Config.ClientID, Now: c.clock.Now})
	verifiedIDToken, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("could not verify the ID token: %w", err)
	}
	if nonce != "" && nonce != verifiedIDToken.Nonce {
		return nil, fmt.Errorf("nonce did not match (wants %s but got %s)", nonce, verifiedIDToken.Nonce)
	}

	if c.useAccessToken {
		accessToken, ok := token.Extra("access_token").(string)
		if !ok {
			return nil, fmt.Errorf("access_token is missing in the token response: %#v", accessToken)
		}

		// Verify the access token with audience validation.
		// If expectedAudience is set, verify against it; otherwise fall back to ClientID.
		// This prevents token confusion attacks where a token issued for one audience
		// is accepted by another service.
		var aud []string
		if len(c.expectedAudience) > 0 {
			aud = c.expectedAudience
		} else {
			aud = []string{c.oauth2Config.ClientID}
		}

		verifier = c.provider.Verifier(&gooidc.Config{
			ClientID:         aud[0],
			Now:              c.clock.Now,
			SkipClientIDCheck: len(aud) == 0,
		})

		_, err := verifier.Verify(ctx, accessToken)
		if err != nil {
			return nil, fmt.Errorf("could not verify the access token: %w", err)
		}

		// There is no `nonce` to check on the `access_token`. We rely on the
		// above `nonce` check on the `id_token`.

		return &oidc.TokenSet{
			IDToken:      accessToken,
			RefreshToken: token.RefreshToken,
		}, nil
	}
	return &oidc.TokenSet{
		IDToken:      idToken,
		RefreshToken: token.RefreshToken,
	}, nil
}