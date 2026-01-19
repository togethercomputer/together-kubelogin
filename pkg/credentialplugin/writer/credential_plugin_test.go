package writer

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/credentialplugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	clientauthenticationv1beta1 "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

func TestWriter_Write(t *testing.T) {
	expiryTime := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("V1Beta1", func(t *testing.T) {
		var stdout bytes.Buffer
		w := Writer{Stdout: &stdout}
		out := credentialplugin.Output{
			Token:                          "test-token-v1beta1",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1beta1",
		}

		err := w.Write(out)
		require.NoError(t, err)

		var execCred clientauthenticationv1beta1.ExecCredential
		err = json.Unmarshal(stdout.Bytes(), &execCred)
		require.NoError(t, err)

		assert.Equal(t, "client.authentication.k8s.io/v1beta1", execCred.APIVersion)
		assert.Equal(t, "ExecCredential", execCred.Kind)
		assert.Equal(t, "test-token-v1beta1", execCred.Status.Token)
		assert.Equal(t, expiryTime.Unix(), execCred.Status.ExpirationTimestamp.Unix())
	})

	t.Run("V1", func(t *testing.T) {
		var stdout bytes.Buffer
		w := Writer{Stdout: &stdout}
		out := credentialplugin.Output{
			Token:                          "test-token-v1",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1",
		}

		err := w.Write(out)
		require.NoError(t, err)

		var execCred clientauthenticationv1.ExecCredential
		err = json.Unmarshal(stdout.Bytes(), &execCred)
		require.NoError(t, err)

		assert.Equal(t, "client.authentication.k8s.io/v1", execCred.APIVersion)
		assert.Equal(t, "ExecCredential", execCred.Kind)
		assert.Equal(t, "test-token-v1", execCred.Status.Token)
		assert.Equal(t, expiryTime.Unix(), execCred.Status.ExpirationTimestamp.Unix())
	})

	t.Run("EmptyAPIVersion_DefaultsToV1Beta1", func(t *testing.T) {
		var stdout bytes.Buffer
		w := Writer{Stdout: &stdout}
		out := credentialplugin.Output{
			Token:                          "test-token-default",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "",
		}

		err := w.Write(out)
		require.NoError(t, err)

		var execCred clientauthenticationv1beta1.ExecCredential
		err = json.Unmarshal(stdout.Bytes(), &execCred)
		require.NoError(t, err)

		assert.Equal(t, "client.authentication.k8s.io/v1beta1", execCred.APIVersion)
		assert.Equal(t, "ExecCredential", execCred.Kind)
		assert.Equal(t, "test-token-default", execCred.Status.Token)
	})

	t.Run("InvalidAPIVersion", func(t *testing.T) {
		var stdout bytes.Buffer
		w := Writer{Stdout: &stdout}
		out := credentialplugin.Output{
			Token:                          "test-token-invalid",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "invalid/v999",
		}

		err := w.Write(out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown apiVersion")
		assert.Contains(t, err.Error(), "invalid/v999")
	})

	t.Run("EmptyToken", func(t *testing.T) {
		var stdout bytes.Buffer
		w := Writer{Stdout: &stdout}
		out := credentialplugin.Output{
			Token:                          "",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1",
		}

		err := w.Write(out)
		require.NoError(t, err)

		var execCred clientauthenticationv1.ExecCredential
		err = json.Unmarshal(stdout.Bytes(), &execCred)
		require.NoError(t, err)
		assert.Equal(t, "", execCred.Status.Token)
	})

	t.Run("ZeroExpiry", func(t *testing.T) {
		var stdout bytes.Buffer
		w := Writer{Stdout: &stdout}
		out := credentialplugin.Output{
			Token:                          "test-token-zero-expiry",
			Expiry:                         time.Time{},
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1",
		}

		err := w.Write(out)
		require.NoError(t, err)

		var execCred clientauthenticationv1.ExecCredential
		err = json.Unmarshal(stdout.Bytes(), &execCred)
		require.NoError(t, err)
		// When expiry is zero, it might be serialized as null or zero time
		if execCred.Status.ExpirationTimestamp != nil {
			assert.True(t, execCred.Status.ExpirationTimestamp.Time.IsZero())
		}
	})
}

func TestGenerateExecCredential(t *testing.T) {
	expiryTime := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("V1Beta1_StructureValidation", func(t *testing.T) {
		out := credentialplugin.Output{
			Token:                          "test-token",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1beta1",
		}

		execCred, err := generateExecCredential(out)
		require.NoError(t, err)

		v1beta1Cred, ok := execCred.(*clientauthenticationv1beta1.ExecCredential)
		require.True(t, ok, "should be v1beta1.ExecCredential")
		assert.Equal(t, "client.authentication.k8s.io/v1beta1", v1beta1Cred.APIVersion)
		assert.Equal(t, "ExecCredential", v1beta1Cred.Kind)
		assert.NotNil(t, v1beta1Cred.Status)
		assert.Equal(t, "test-token", v1beta1Cred.Status.Token)
		assert.Equal(t, expiryTime, v1beta1Cred.Status.ExpirationTimestamp.Time)
	})

	t.Run("V1_StructureValidation", func(t *testing.T) {
		out := credentialplugin.Output{
			Token:                          "test-token",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1",
		}

		execCred, err := generateExecCredential(out)
		require.NoError(t, err)

		v1Cred, ok := execCred.(*clientauthenticationv1.ExecCredential)
		require.True(t, ok, "should be v1.ExecCredential")
		assert.Equal(t, "client.authentication.k8s.io/v1", v1Cred.APIVersion)
		assert.Equal(t, "ExecCredential", v1Cred.Kind)
		assert.NotNil(t, v1Cred.Status)
		assert.Equal(t, "test-token", v1Cred.Status.Token)
		assert.Equal(t, expiryTime, v1Cred.Status.ExpirationTimestamp.Time)
	})

	t.Run("EmptyAPIVersion_ReturnsV1Beta1", func(t *testing.T) {
		out := credentialplugin.Output{
			Token:                          "test-token",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "",
		}

		execCred, err := generateExecCredential(out)
		require.NoError(t, err)

		v1beta1Cred, ok := execCred.(*clientauthenticationv1beta1.ExecCredential)
		require.True(t, ok, "should default to v1beta1.ExecCredential")
		assert.Equal(t, "client.authentication.k8s.io/v1beta1", v1beta1Cred.APIVersion)
	})

	t.Run("UnsupportedAPIVersion", func(t *testing.T) {
		out := credentialplugin.Output{
			Token:                          "test-token",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v2alpha1",
		}

		execCred, err := generateExecCredential(out)
		assert.Error(t, err)
		assert.Nil(t, execCred)
		assert.Contains(t, err.Error(), "unknown apiVersion")
		assert.Contains(t, err.Error(), "v2alpha1")
	})

	t.Run("TypeMetaFieldsSet", func(t *testing.T) {
		out := credentialplugin.Output{
			Token:                          "test-token",
			Expiry:                         expiryTime,
			ClientAuthenticationAPIVersion: "client.authentication.k8s.io/v1",
		}

		execCred, err := generateExecCredential(out)
		require.NoError(t, err)

		v1Cred := execCred.(*clientauthenticationv1.ExecCredential)
		assert.NotEmpty(t, v1Cred.APIVersion, "APIVersion should be set")
		assert.Equal(t, "ExecCredential", v1Cred.Kind, "Kind should be ExecCredential")
	})
}
