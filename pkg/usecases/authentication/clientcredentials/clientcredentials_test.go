package clientcredentials

import (
	"context"
	"errors"
	"testing"

	"github.com/togethercomputer/together-kubelogin/mocks/github.com/togethercomputer/together-kubelogin/pkg/oidc/client_mock"
	"github.com/togethercomputer/together-kubelogin/pkg/oidc"
	"github.com/togethercomputer/together-kubelogin/pkg/oidc/client"
	"github.com/togethercomputer/together-kubelogin/pkg/testing/logger"
)

func TestClientCredentials(t *testing.T) {
	ctx := context.TODO()

	t.Run("Authorization error", func(t *testing.T) {
		mockClient := client_mock.NewMockInterface(t)
		cc := &ClientCredentials{
			Logger: logger.New(t),
		}
		errTest := errors.New("make sure Do() passes err")
		in := client.GetTokenByClientCredentialsInput{}
		mockClient.EXPECT().GetTokenByClientCredentials(ctx, in).Return(nil, errTest).Once()
		opts := client.GetTokenByClientCredentialsInput{}
		_, err := cc.Do(ctx, &opts, mockClient)
		if !errors.Is(err, errTest) {
			t.Errorf("returned error is not the test error: %v", err)
		}
	})

	t.Run("Server returns access_token", func(t *testing.T) {
		mockClient := client_mock.NewMockInterface(t)
		cc := &ClientCredentials{
			Logger: logger.New(t),
		}
		in := client.GetTokenByClientCredentialsInput{}
		testToken := &oidc.TokenSet{}
		mockClient.EXPECT().GetTokenByClientCredentials(ctx, in).Return(testToken, nil).Once()
		opts := client.GetTokenByClientCredentialsInput{}
		ts, err := cc.Do(ctx, &opts, mockClient)
		if err != nil {
			t.Errorf("returned unexpected error: %v", err)
		}
		if ts != testToken {
			t.Errorf("returned unexpected token set: %v", ts)
		}
	})

}
