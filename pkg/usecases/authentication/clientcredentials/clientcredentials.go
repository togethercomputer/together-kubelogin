package clientcredentials

import (
	"context"
	"fmt"

	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/infrastructure/logger"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/oidc"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/oidc/client"
)

// DeviceCode provides the oauth2 device code flow.
type ClientCredentials struct {
	Logger logger.Interface
}

func (u *ClientCredentials) Do(ctx context.Context, in *client.GetTokenByClientCredentialsInput, oidcClient client.Interface) (*oidc.TokenSet, error) {
	u.Logger.V(1).Infof("starting the oauth2 client credentials code flow")
	if in == nil {
		return nil, fmt.Errorf("nil input")
	}
	tokenSet, err := oidcClient.GetTokenByClientCredentials(ctx, *in)
	if err != nil {
		return nil, fmt.Errorf("authorization error: %w", err)
	}
	u.Logger.V(1).Infof("finished the oauth2 client credentials code flow")
	return tokenSet, nil
}
