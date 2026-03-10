# -------- Stage 1: baixar e extrair chart --------
FROM alpine:3.20 AS charts

ARG CHART_VERSION=1.0.0

RUN apk add --no-cache curl tar

WORKDIR /charts

RUN curl -L https://paulovigne.github.io/hitcounter/hitcounter-${CHART_VERSION}.tgz \
    -o hitcounter.tgz \
 && tar -xzf hitcounter.tgz \
 && rm hitcounter.tgz


# -------- Stage 2: imagem final do operator --------
FROM quay.io/operator-framework/helm-operator:v1.42.0

ENV HOME=/opt/helm

WORKDIR ${HOME}

COPY watches.yaml ${HOME}/watches.yaml

# copiar chart extraído
COPY --from=charts /charts/hitcounter ${HOME}/charts/hitcounter
