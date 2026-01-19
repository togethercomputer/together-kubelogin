# together-kubelogin

**Kubernetes Authentication Plugin for Together AI**

This is a kubectl plugin for [Kubernetes OpenID Connect (OIDC) authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) customized for Together AI platform.

Together-kubelogin enables your customers to authenticate to Kubernetes clusters using Together AI identity provider. It's designed to run as a [client-go credential plugin](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins).

When you run kubectl, together-kubelogin opens the browser for Together AI authentication.
After successful login, together-kubelogin gets a token from Together AI and kubectl accesses Kubernetes APIs with the token.

![Diagram of the credential plugin](docs/credential-plugin-diagram.svg)

## Getting Started

### Installation

Build from source or download from releases:

```sh
# Build from source
git clone https://github.com/pahluwalia-tcloud/together-kubelogin.git
cd together-kubelogin
make build

# Copy binary to PATH as kubectl plugin
sudo cp bin/kubelogin /usr/local/bin/kubectl-together_login
```

The binary must be named `kubectl-together_login` on your PATH.
When you invoke `kubectl together-login`, kubectl finds it by the [naming convention of kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

### Configuration

You need to configure your kubeconfig with Together AI OIDC settings. Your kubeconfig should look like this:

```yaml
users:
  - name: together-ai
    user:
      exec:
        apiVersion: client.authentication.k8s.io/v1
        command: kubectl
        args:
          - together-login
          - get-token
          - --oidc-issuer-url=https://auth.together.ai    # Together AI OIDC issuer
          - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
          - --oidc-client-secret=YOUR_TOGETHER_CLIENT_SECRET  # Optional
```

**Together AI OIDC Configuration:**
- **Issuer URL**: `https://auth.together.ai` (placeholder - use your actual Together AI OIDC issuer)
- **Client ID**: Obtain from Together AI platform
- **Client Secret**: Optional, depending on your OIDC configuration

See the [setup guide](docs/setup.md) for complete details on:
- Setting up the Kubernetes API server OIDC flags
- Creating cluster role bindings
- Configuring Together AI OIDC provider

### Usage

Run kubectl commands as usual:

```sh
kubectl get pods
```

Kubectl executes together-kubelogin before calling the Kubernetes APIs.
Together-kubelogin automatically opens the browser for Together AI authentication.

After authentication, together-kubelogin returns the credentials to kubectl.
Kubectl then calls the Kubernetes APIs with the credentials.

```console
% kubectl get pods
Open http://localhost:8000 for authentication
NAME                          READY   STATUS    RESTARTS   AGE
my-app-86c78fdccd-nzmd5       1/1     Running   0          26d
```

Together-kubelogin stores the ID token and refresh token to the cache.
If the ID token is valid, it returns it immediately.
If the ID token has expired, it refreshes the token using the refresh token.
If the refresh token has expired, it performs re-authentication via the browser.

## Troubleshooting

### Token cache

Together-kubelogin stores the token cache to the file system by default.
For enhanced security, it is recommended to store it to the keyring.
See the [token cache](docs/usage.md#token-cache) for details.

You can log out by deleting the token cache:

```console
% kubectl together-login clean
Deleted the token cache at /home/user/.kube/cache/oidc-login
Deleted the token cache from the keyring
```

Together-kubelogin will ask you to log in via the browser again.
If the browser has a cookie for Together AI, you may need to log out from Together AI or clear the cookie.

### ID token claims

You can run the `setup` command to dump the claims of an ID token from Together AI:

```console
% kubectl together-login setup \
    --oidc-issuer-url=https://auth.together.ai \
    --oidc-client-id=YOUR_CLIENT_ID
...
You got a token with the following claims:

{
  "sub": "********",
  "iss": "https://auth.together.ai",
  "aud": "********",
  ...
}
```

You can set `-v1` option to increase the log level:

```yaml
users:
  - name: together-ai
    user:
      exec:
        apiVersion: client.authentication.k8s.io/v1
        command: kubectl
        args:
          - together-login
          - get-token
          - -v1
```

You can run the [acceptance test](acceptance_test) to verify if together-kubelogin works with Together AI.

## Docs

- [Setup guide](docs/setup.md)
- [Usage and options](docs/usage.md)
- [Standalone mode](docs/standalone-mode.md)
- [System test](system_test)
- [Acceptance_test for identity providers](acceptance_test)

## Contributing

This project is a fork of [int128/kubelogin](https://github.com/int128/kubelogin), customized for Together AI.

This is open source software licensed under Apache License 2.0.
Feel free to open issues and pull requests for improving code and documents.

## Credits

Based on [kubelogin](https://github.com/int128/kubelogin) by int128.
