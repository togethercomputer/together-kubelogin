package setup

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/togethercomputer/together-kubelogin/mocks/github.com/togethercomputer/together-kubelogin/pkg/usecases/authentication_mock"
	"github.com/togethercomputer/together-kubelogin/pkg/oidc"
	testingJWT "github.com/togethercomputer/together-kubelogin/pkg/testing/jwt"
	"github.com/togethercomputer/together-kubelogin/pkg/testing/logger"
	"github.com/togethercomputer/together-kubelogin/pkg/tlsclientconfig"
	"github.com/togethercomputer/together-kubelogin/pkg/usecases/authentication"
)

func TestSetup_Do(t *testing.T) {
	issuedIDToken := testingJWT.EncodeF(t, func(claims *testingJWT.Claims) {
		claims.Issuer = "https://issuer.example.com"
		claims.Subject = "YOUR_SUBJECT"
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	})
	dummyTLSClientConfig := tlsclientconfig.Config{
		CACertFilename: []string{"/path/to/cert"},
	}
	var grantOptionSet authentication.GrantOptionSet

	ctx := context.Background()
	in := Input{
		IssuerURL:       "https://accounts.google.com",
		ClientID:        "YOUR_CLIENT_ID",
		ClientSecret:    "YOUR_CLIENT_SECRET",
		ExtraScopes:     []string{"email"},
		GrantOptionSet:  grantOptionSet,
		TLSClientConfig: dummyTLSClientConfig,
		ChangedFlags: []string{
			"--oidc-issuer-url=https://accounts.google.com",
			"--oidc-client-id=YOUR_CLIENT_ID",
		},
	}
	mockAuthentication := authentication_mock.NewMockInterface(t)
	mockAuthentication.EXPECT().
		Do(ctx, authentication.Input{
			Provider: oidc.Provider{
				IssuerURL:    "https://accounts.google.com",
				ClientID:     "YOUR_CLIENT_ID",
				ClientSecret: "YOUR_CLIENT_SECRET",
				ExtraScopes:  []string{"email"},
			},
			GrantOptionSet:  grantOptionSet,
			TLSClientConfig: dummyTLSClientConfig,
		}).
		Return(&authentication.Output{
			TokenSet: oidc.TokenSet{
				IDToken:      issuedIDToken,
				RefreshToken: "YOUR_REFRESH_TOKEN",
			},
		}, nil)
	u := Setup{
		Authentication: mockAuthentication,
		Logger:         logger.New(t),
	}
	if err := u.Do(ctx, in); err != nil {
		t.Errorf("Do returned error: %+v", err)
	}
}
