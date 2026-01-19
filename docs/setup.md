# Kubernetes OpenID Connect authentication with Together AI

This document guides how to set up Kubernetes OpenID Connect (OIDC) authentication using Together AI.
Follow these steps:

1. Set up Together AI OIDC provider
1. Verify authentication
1. Bind a cluster role
1. Set up the Kubernetes API server
1. Set up the kubeconfig
1. Verify cluster access

## 1. Set up Together AI OIDC Provider

Together-kubelogin supports authentication flows such as Device Authorization Grant or Authorization Code Flow.
For details of flows supported, see the [usage](usage.md).

**Together AI OIDC Configuration:**
- Issuer URL: `https://auth.together.ai` (placeholder - use your actual Together AI OIDC issuer)
- Client ID: Obtain from Together AI platform
- Client Secret: Optional, obtain if required by your configuration

## 2. Authenticate with Together AI

Run the following command to test the authentication setup:

```sh
kubectl together-login setup --oidc-issuer-url=https://auth.together.ai --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
```

Set the following flags:

- `--oidc-issuer-url`: Together AI OIDC issuer URL
- `--oidc-client-id`: Your Together AI client ID
- `--oidc-client-secret`: Your Together AI client secret (if required)

Together AI supports Authorization Code Flow (default):
- Starts a local server for authentication
- Launches the browser to Together AI's authentication page
- You log in with your Together AI credentials

For headless environments, use Device Authorization Grant:

```sh
kubectl together-login setup \
  --oidc-issuer-url=https://auth.together.ai \
  --oidc-client-id=YOUR_TOGETHER_CLIENT_ID \
  --grant-type=device-code
```

You can see all available options:

```sh
kubectl together-login setup --help
```

## 3. Bind a cluster role

Bind the `cluster-admin` role to Together AI users:

```sh
kubectl create clusterrolebinding together-ai-admin \
  --clusterrole=cluster-admin \
  --user='https://auth.together.ai#YOUR_SUBJECT'
```

Replace `YOUR_SUBJECT` with the user's subject claim from their Together AI token (typically email or user ID).

## 4. Set up the Kubernetes API server

Add the following flags to your kube-apiserver configuration:

```
--oidc-issuer-url=https://auth.together.ai
--oidc-client-id=YOUR_TOGETHER_CLIENT_ID
--oidc-username-claim=email
--oidc-groups-claim=groups
```

See [Kubernetes Authenticating: OpenID Connect Tokens](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) for all available flags.

## 5. Set up the kubeconfig

Add the `together-ai` user to your kubeconfig:

```sh
kubectl config set-credentials together-ai \
  --exec-interactive-mode=Never \
  --exec-api-version=client.authentication.k8s.io/v1 \
  --exec-command=kubectl \
  --exec-arg=together-login \
  --exec-arg=get-token \
  --exec-arg=--oidc-issuer-url=https://auth.together.ai \
  --exec-arg=--oidc-client-id=YOUR_TOGETHER_CLIENT_ID
```

If Together AI requires a client secret, add the secret argument:

```sh
--exec-arg=--oidc-client-secret=YOUR_TOGETHER_CLIENT_SECRET
```

For enhanced security, store the token cache in the system keyring:

```sh
--exec-arg=--token-cache-storage=keyring
```

If you encounter errors with keyring storage, see the [token cache](usage.md#token-cache) documentation.

## 6. Verify cluster access

Verify that you can access the Kubernetes cluster with Together AI authentication:

```sh
kubectl --user=together-ai cluster-info
```

Switch the current context to use Together AI authentication:

```sh
kubectl config set-context --current --user=together-ai
```

You can share the kubeconfig configuration with your team members for easy onboarding.
