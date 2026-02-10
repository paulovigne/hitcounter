# =========================
# Stage 1 - Build
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Dependências básicas
RUN apk add --no-cache ca-certificates

# Copia arquivos de módulo primeiro (cache eficiente)
COPY go.mod hitcounter.go ./
RUN go mod tidy
RUN go mod download

# Build binário estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o hitcounter

# =========================
# Stage 2 - Runtime
# =========================
FROM gcr.io/distroless/base-debian13

WORKDIR /app

# Copia binário e certificados
COPY --from=builder /app/hitcounter /app/hitcounter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Porta da aplicação
EXPOSE 8080

# Variáveis default
ENV PORT=8080
ENV ENABLE_REDIS=false

# Usuário não-root (distroless já usa)
USER nonroot:nonroot

ENTRYPOINT ["/app/hitcounter"]
