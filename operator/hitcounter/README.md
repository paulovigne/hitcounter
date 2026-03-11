# hitcounter

Operator to deploy **Hitcounter** on Kubernetes with multi-platform exposure support.

Supports:

- Kubernetes Ingress
- Kubernetes Gateway API
- OpenShift Route
- Istio Gateway
- Optional Redis
- HPA
- TLS with cert-manager

---

## TL;DR

```bash
kubectl apply -f https://paulovigne.github.io/hitcounter/hitcounter-operator.yaml
```

---

## Introduction

This operator deploys the **Hitcounter** application on a Kubernetes cluster using the Helm package manager.

It provides a unified `exposure.type` abstraction allowing the application to be exposed through:

- `ingress`
- `gatewayapi`
- `route`
- `istio`

TLS configuration is centralized and automatically mapped to each provider.

---

## Prerequisites

- Kubernetes >= 1.24

Depending on `exposure.type`:

| Exposure Type | Requirement |
|---------------|------------|
| ingress | Ingress Controller (e.g. Traefik, NGINX) |
| gatewayapi | Gateway API CRDs installed |
| route | OpenShift cluster |
| istio | Istio installed |

For automatic TLS:

- cert-manager

---

## Installing the Operator

```bash
kubectl apply -f https://paulovigne.github.io/hitcounter/hitcounter-operator.yaml
```

## Uninstalling the Operator

```bash
kubectl apply -f https://paulovigne.github.io/hitcounter/hitcounter-operator.yaml
```

> [!NOTE]
> The operator will be deployed on `operators` namespace.

---

# Configuration

The following table lists the configurable parameters of the operator and their default values.

---

## Global Parameters

| Name | Description | Default |
|------|------------|---------|
| `replicaCount` | Number of replicas (ignored if HPA enabled) | `1` |
| `nameOverride` | Override release name | `""` |
| `fullnameOverride` | Override full name | `""` |

---

## Security Context Parameters

| Name | Description | Default |
|------|------------|---------|
| `securityContext.runAsUser` | Run container as non-root user. `65532` is commonly used for distroless images. On OpenShift, SCC may override UID range. | `65532` |
| `securityContext.runAsGroup` | Primary GID for the container process. On OpenShift, SCC may override GID range. | `65532` |
| `securityContext.fsGroup` | FSGroup applied to mounted volumes for shared file permissions. | `65532` |
| `securityContext.readOnlyRootFilesystem` | Mount container root filesystem as read-only for additional security hardening. | `true` |
| `securityContext.openshift` | Enable OpenShift-specific security handling (SCC-compatible UID/GID behavior). | `false` |

---

## Image Parameters

| Name | Description | Default |
|------|------------|---------|
| `image.repository` | Image repository | `ghcr.io/paulovigne/hitcounter` |
| `image.tag` | Image tag | `main` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `image.pullSecret` | Existing imagePullSecret | `null` |

---

## Service Parameters

| Name | Description | Default |
|------|------------|---------|
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `service.targetPort` | Container port | `8080` |

---

## Exposure Parameters

| Name | Description | Default |
|------|------------|---------|
| `exposure.type` | Exposure strategy (`ingress`, `gatewayapi`, `route`, `istio`) | `ingress` |
| `exposure.host` | Application hostname | `hitcounter.local` |
| `exposure.path` | Application path | `/` |

---

## Exposure TLS Parameters

| Name | Description | Default |
|------|------------|---------|
| `exposure.tls.enabled` | Enable TLS | `false` |
| `exposure.tls.mode` | TLS mode (provider dependent) | `terminate` |
| `exposure.tls.secret.create` | Create TLS secret | `false` |
| `exposure.tls.secret.name` | TLS secret name | `hitcounter-tls` |
| `exposure.tls.secret.crt` | TLS certificate (PEM) | `null` |
| `exposure.tls.secret.key` | TLS private key (PEM) | `null` |
| `exposure.tls.secret.ca` | Optional CA certificate | `null` |

---

## cert-manager Parameters

| Name | Description | Default |
|------|------------|---------|
| `exposure.tls.secret.certManager.enabled` | Enable cert-manager | `false` |
| `exposure.tls.secret.certManager.issuerName` | Issuer name | `null` |
| `exposure.tls.secret.certManager.issuerKind` | Issuer kind (`Issuer` or `ClusterIssuer`) | `null` |

---

## Ingress Parameters (if exposure.type=ingress)

| Name | Description | Default |
|------|------------|---------|
| `exposure.ingress.className` | Ingress class | `traefik` |

---

## Gateway API Parameters (if exposure.type=gatewayapi)

| Name | Description | Default |
|------|------------|---------|
| `exposure.gatewayapi.gatewayClassName` | GatewayClass name | `traefik` |
| `exposure.gatewayapi.namespace` | Gateway namespace | `default` |

---

## OpenShift Route Parameters (if exposure.type=route)

| Name | Description | Default |
|------|------------|---------|
| `exposure.route.termination` | TLS termination (`edge`, `passthrough`, `reencrypt`) | `edge` |

---

## Istio Parameters (if exposure.type=istio)

| Name | Description | Default |
|------|------------|---------|
| `exposure.istio.namespace` | Istio namespace | `istio-system` |
| `exposure.istio.ingressGatewayName` | Istio ingress gateway name | `istio-ingressgateway` |
| `exposure.istio.tlsMode` | TLS mode (`SIMPLE`, `PASSTHROUGH`, `MUTUAL`) | `SIMPLE` |

---

## HPA Parameters

| Name | Description | Default |
|------|------------|---------|
| `hpa.enabled` | Enable HPA | `false` |
| `hpa.minReplicas` | Minimum replicas | `2` |
| `hpa.maxReplicas` | Maximum replicas | `6` |
| `hpa.cpuUtilization` | Target CPU utilization (%) | `60` |

---

## Redis Parameters

| Name | Description | Default |
|------|------------|---------|
| `config.enableRedis` | Deploy embedded Redis | `true` |
| `config.redisHost` | External Redis host | `null` |
| `config.redisPort` | External Redis port | `"6379"` |
| `redis.image` | Redis image | `redis:7-alpine` |
| `redis.storage` | Redis PVC size | `1Gi` |

---

## Resource Parameters

| Name | Description | Default |
|------|------------|---------|
| `requests.cpu` | CPU requests | `50m` |
| `requests.memory` | Memory requests | `64Mi` |
| `limits.cpu` | CPU limits | `200m` |
| `limits.memory` | Memory limits | `128Mi` |

---

# TLS Mode Mapping

| exposure.type | TLS Mode Values |
|---------------|-----------------|
| ingress | `terminate` |
| gatewayapi | `terminate`, `passthrough` |
| route | `edge`, `passthrough`, `reencrypt` |
| istio | `SIMPLE`, `PASSTHROUGH`, `MUTUAL` |

---

# Examples of Application Deployment Configurations

## Ingress + Self Certificate

```bash
cat <<EOF | kubectl -n hitcounter apply -f -
apiVersion: apps.vigne.com.br/v1alpha1
kind: HitCounter
metadata:
  name: hitcounter
spec:

  exposure:
    type: ingress
    host: app.myhost.com

    ingress:
      className: traefik

    tls:
      enabled: true
      secret:
        name: hitcounter-tls
        create: true
        crt: |
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----

        key: |
          -----BEGIN PRIVATE KEY-----
          ...
          -----END PRIVATE KEY-----

        ca: |
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
EOF          
```

## Ingress + cert-manager

```bash
cat <<EOF | kubectl -n hitcounter apply -f -
apiVersion: apps.vigne.com.br/v1alpha1
kind: HitCounter
metadata:
  name: hitcounter
spec:

  exposure:
    type: ingress
    host: app.myhost.com

    ingress:
      className: traefik

    tls:
      enabled: true

      secret:
        name: hitcounter-tls

        certManager:
          enabled: true
          issuerName: selfsigned-cluster-issuer
          issuerKind: ClusterIssuer
EOF
```

---

## OpenShift Route + Auto hostname and certificate

```bash
cat <<EOF | kubectl -n hitcounter apply -f -
apiVersion: apps.vigne.com.br/v1alpha1
kind: HitCounter
metadata:
  name: hitcounter
spec:

  securityContext:
    openshift: true

  exposure:
    type: route

    host: null

    tls:
      enabled: true
EOF
```

---

## Istio Gateway with TLS + cert-manager

```bash
cat <<EOF | kubectl -n hitcounter apply -f -
apiVersion: apps.vigne.com.br/v1alpha1
kind: HitCounter
metadata:
  name: hitcounter
spec:

  exposure:
    type: istio
    host: app.myhost.com

    tls:
      enabled: true
      mode: SIMPLE

      secret:
        name: hitcounter-tls

        certManager:
          enabled: true
          issuerName: selfsigned-cluster-issuer
          issuerKind: ClusterIssuer
EOF
```

## K8S Gateway API with TLS + cert-manager

```bash
cat <<EOF | kubectl -n hitcounter apply -f -
apiVersion: apps.vigne.com.br/v1alpha1
kind: HitCounter
metadata:
  name: hitcounter
spec:

  exposure:
    type: gatewayapi
    host: app.myhost.com

    gatewayapi:
      gatewayClassName: traefik

    tls:
      enabled: true
      mode: terminate

      secret:
        name: hitcounter-tls

        certManager:
          enabled: true
          issuerName: selfsigned-cluster-issuer
          issuerKind: ClusterIssuer
EOF
```

## Uninstall HitCounter App

### Remove HitCounter Resources

```bash
kubectl delete hitcounter --all -A
```

or

```bash
kubectl delete hitcounters.apps.vigne.com.br --all -A
```

### Remove HitCounter CRD

```bash
kubectl delete crd hitcounters.apps.vigne.com.br
```

## CertManager for testing

### Helm Install

```bash
helm install \
  cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.19.4 \
  --namespace cert-manager \
  --create-namespace \
  --set crds.enabled=true
```
### Self Signed cluster-issuer

```bash
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-cluster-issuer
spec:
  selfSigned: {}
EOF
```

# Production Recommendations

- Avoid `main` tag — use immutable versions
- Enable HPA
- Define resource limits
- Use cert-manager for TLS automation
- Use external Redis in production
- Configure liveness/readiness probes

---

# Maintainers

- Paulo Vigne