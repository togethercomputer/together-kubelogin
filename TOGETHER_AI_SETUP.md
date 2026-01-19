# Together AI Kubernetes Authentication Setup

This guide shows how to set up Kubernetes authentication using Together AI's OIDC provider with together-kubelogin.

## Prerequisites

- Kubernetes cluster with OIDC authentication enabled
- Together AI OIDC credentials (Client ID and optionally Client Secret)
- kubectl installed

## Step 1: Install together-kubelogin

Build and install the plugin:

```bash
git clone https://github.com/pavneetahluwalia/together-kubelogin.git
cd together-kubelogin
make build

# Install as kubectl plugin
sudo cp bin/kubelogin /usr/local/bin/kubectl-together_login
chmod +x /usr/local/bin/kubectl-together_login
```

Verify installation:

```bash
kubectl together-login --help
```

## Step 2: Configure Kubernetes API Server

Your Kubernetes API server must be configured with Together AI OIDC settings. Add these flags to your `kube-apiserver`:

```bash
--oidc-issuer-url=https://auth.together.ai
--oidc-client-id=YOUR_TOGETHER_CLIENT_ID
--oidc-username-claim=email
--oidc-groups-claim=groups
```

**Note**: Replace `https://auth.together.ai` with the actual Together AI OIDC issuer URL provided by Together AI.

## Step 3: Configure kubeconfig

Add the Together AI authentication configuration to your kubeconfig:

### Basic Configuration

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
          - --oidc-issuer-url=https://auth.together.ai
          - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
```

### With Client Secret (if required)

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
          - --oidc-issuer-url=https://auth.together.ai
          - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
          - --oidc-client-secret=YOUR_TOGETHER_CLIENT_SECRET
```

### With Token Cache in Keyring (Recommended for Security)

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
          - --oidc-issuer-url=https://auth.together.ai
          - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
          - --token-cache-storage=keyring
```

### Using kubectl to Set Configuration

You can also use kubectl commands to configure the user:

```bash
kubectl config set-credentials together-ai \
  --exec-api-version=client.authentication.k8s.io/v1 \
  --exec-command=kubectl \
  --exec-arg=together-login \
  --exec-arg=get-token \
  --exec-arg=--oidc-issuer-url=https://auth.together.ai \
  --exec-arg=--oidc-client-id=YOUR_TOGETHER_CLIENT_ID
```

## Step 4: Create Cluster Role Binding

Bind appropriate roles to users based on their Together AI identity:

```bash
# Replace with your Together AI issuer URL and user subject
kubectl create clusterrolebinding together-ai-admin \
  --clusterrole=cluster-admin \
  --user='https://auth.together.ai#USER_SUBJECT'
```

The user subject is typically the `sub` claim from the OIDC token (usually an email or user ID).

## Step 5: Test Authentication

Switch to the Together AI user context:

```bash
kubectl config set-context --current --user=together-ai
```

Test access:

```bash
kubectl get pods
```

On first run, together-kubelogin will:
1. Open your browser to Together AI's authentication page
2. Prompt you to log in with your Together AI credentials
3. Store the authentication tokens locally
4. Return the token to kubectl

Subsequent kubectl commands will use the cached token until it expires.

## Troubleshooting

### View Token Claims

To see what claims are in your token:

```bash
kubectl together-login setup \
  --oidc-issuer-url=https://auth.together.ai \
  --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
```

### Clear Token Cache

If you need to re-authenticate:

```bash
kubectl together-login clean
```

### Enable Debug Logging

Add `-v1` flag to see detailed logs:

```yaml
users:
  - name: together-ai
    user:
      exec:
        args:
          - together-login
          - get-token
          - --oidc-issuer-url=https://auth.together.ai
          - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
          - -v1
```

## Together AI OIDC Configuration Reference

Replace these placeholder values with your actual Together AI configuration:

- **Issuer URL**: `https://auth.together.ai` (get actual URL from Together AI)
- **Client ID**: Obtain from Together AI platform console
- **Client Secret**: Optional, obtain from Together AI if required
- **Scopes**: Typically `openid email profile` (may vary)

## Advanced Options

### Device Code Flow

For headless servers or environments without a browser:

```yaml
args:
  - together-login
  - get-token
  - --oidc-issuer-url=https://auth.together.ai
  - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
  - --grant-type=device-code
```

### Custom Token Cache Directory

```yaml
args:
  - together-login
  - get-token
  - --oidc-issuer-url=https://auth.together.ai
  - --oidc-client-id=YOUR_TOGETHER_CLIENT_ID
  - --token-cache-dir=/custom/path/.kube/cache
```

## Security Best Practices

1. **Use Keyring Storage**: Store tokens in system keyring instead of filesystem
2. **Limit Scope**: Request only necessary OIDC scopes
3. **Role-Based Access**: Use Kubernetes RBAC to limit user permissions
4. **Regular Rotation**: Configure short token lifetimes and refresh tokens
5. **Audit Logging**: Enable Kubernetes audit logging to track authentication events
