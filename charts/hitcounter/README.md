# hitcounter

Helm chart to deploy **Hitcounter** on Kubernetes with multi-platform exposure support.

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
helm install hitcounter ./hitcounter
```

With custom values:

```bash
helm install hitcounter ./hitcounter -f values.yaml
```

---

## Introduction

This chart deploys the **Hitcounter** application on a Kubernetes cluster using the Helm package manager.

It provides a unified `exposure.type` abstraction allowing the application to be exposed through:

- `ingress`
- `gatewayapi`
- `route`
- `istio`

TLS configuration is centralized and automatically mapped to each provider.

---

## Prerequisites

- Kubernetes >= 1.24
- Helm >= 3.8

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

## Installing the Chart

```bash
helm install hitcounter ./hitcounter
```

Install with custom namespace:

```bash
helm install hitcounter ./hitcounter -n apps --create-namespace
```

---

## Uninstalling the Chart

```bash
helm uninstall hitcounter
```

---

# Configuration

The following table lists the configurable parameters of the chart and their default values.

---

## Global Parameters

| Name | Description | Default |
|------|------------|---------|
| `replicaCount` | Number of replicas (ignored if HPA enabled) | `1` |
| `nameOverride` | Override release name | `""` |
| `fullnameOverride` | Override full name | `""` |

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

# Example Configurations

## Ingress + cert-manager

```yaml
exposure:
  type: ingress
  host: app.example.com
  tls:
    enabled: true
    secret:
      create: false
      name: app-tls
      certManager:
        enabled: true
        issuerName: letsencrypt
        issuerKind: ClusterIssuer
```

---

## OpenShift Route

```yaml
exposure:
  type: route
  route:
    termination: edge
```

---

## Istio Gateway with TLS

```yaml
exposure:
  type: istio
  tls:
    enabled: true
    mode: SIMPLE
```

---

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