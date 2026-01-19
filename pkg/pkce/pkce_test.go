package pkce

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("MethodS256", func(t *testing.T) {
		params, err := New(MethodS256)
		require.NoError(t, err)
		assert.Equal(t, MethodS256, params.Method)
		assert.NotEmpty(t, params.Verifier, "verifier should be generated")
		// Verify verifier meets PKCE requirements (RFC 7636)
		assert.GreaterOrEqual(t, len(params.Verifier), 43, "verifier should be at least 43 characters")
		assert.LessOrEqual(t, len(params.Verifier), 128, "verifier should be at most 128 characters")
	})

	t.Run("NoMethod", func(t *testing.T) {
		params, err := New(NoMethod)
		require.NoError(t, err)
		assert.Equal(t, NoMethod, params.Method)
		assert.Empty(t, params.Verifier, "verifier should be empty for NoMethod")
	})

	t.Run("UniqueVerifiers", func(t *testing.T) {
		// Verify that multiple calls generate unique verifiers
		params1, err := New(MethodS256)
		require.NoError(t, err)
		params2, err := New(MethodS256)
		require.NoError(t, err)
		assert.NotEqual(t, params1.Verifier, params2.Verifier, "verifiers should be unique")
	})
}

func TestParams_AuthCodeOption(t *testing.T) {
	t.Run("MethodS256", func(t *testing.T) {
		params := Params{
			Method:   MethodS256,
			Verifier: "test-verifier-code-challenge",
		}
		option := params.AuthCodeOption()
		assert.NotNil(t, option, "should return S256 challenge option")
	})

	t.Run("NoMethod", func(t *testing.T) {
		params := Params{
			Method:   NoMethod,
			Verifier: "",
		}
		option := params.AuthCodeOption()
		assert.Nil(t, option, "should return nil for NoMethod")
	})

	t.Run("EmptyParams", func(t *testing.T) {
		params := Params{}
		option := params.AuthCodeOption()
		assert.Nil(t, option, "should return nil for empty params")
	})
}

func TestParams_TokenRequestOption(t *testing.T) {
	t.Run("MethodS256", func(t *testing.T) {
		params := Params{
			Method:   MethodS256,
			Verifier: "test-verifier-token-request",
		}
		option := params.TokenRequestOption()
		assert.NotNil(t, option, "should return verifier option")
	})

	t.Run("NoMethod", func(t *testing.T) {
		params := Params{
			Method:   NoMethod,
			Verifier: "",
		}
		option := params.TokenRequestOption()
		assert.Nil(t, option, "should return nil for NoMethod")
	})

	t.Run("EmptyParams", func(t *testing.T) {
		params := Params{}
		option := params.TokenRequestOption()
		assert.Nil(t, option, "should return nil for empty params")
	})
}

func TestMethod_Constants(t *testing.T) {
	t.Run("MethodValues", func(t *testing.T) {
		// Verify that method constants have expected values
		assert.Equal(t, Method(0), NoMethod)
		assert.Equal(t, Method(1), MethodS256)
		assert.NotEqual(t, NoMethod, MethodS256)
	})
}
