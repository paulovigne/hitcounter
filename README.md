# Hit Counter 🧮

Aplicação simples escrita em **Go** que expõe um contador de acessos via HTTP.  
O contador pode operar de forma **stateless** ou utilizando **Redis** para persistência e contagem compartilhada entre múltiplas instâncias.

O projeto foi criado como exemplo prático de:
- Aplicação cloud-native
- Uso de Docker Compose
- Conversão de Docker Compose para manifestos Kubernetes
- Separação entre aplicação stateless e serviço stateful

📁 Estrutura do Repositório

```
.
├── docker-compose.yml
├── Dockerfile
├── hitcounter.go
├── helm/
|   ├── Chart.yaml
|   ├── values.yaml
|   ├── templates/
│      ├── _helpers.tpl
│      ├── serviceaccount.yaml
│      ├── configmap.yaml
│      ├── hitcounter-deployment.yaml
│      ├── hitcounter-service.yaml
│      ├── redis-statefulset.yaml
│      ├── redis-headless-service.yaml
│      ├── route.yaml
│      └── ingress.yaml
├── k8s/
│   ├── hitcounter-deployment.yaml
│   ├── hitcounter-service.yaml
|   ├── configmap.yaml
|   ├── ingress.yaml
|   ├── serviceaccount.yaml
|   ├── redis-headless-service.yaml
│   └── redis-statefulset.yaml
└── README.md
```

---

## 📌 Visão Geral da Arquitetura

Componentes principais:

- **HitCounter**
  - Aplicação HTTP em Go
  - Porta padrão: `8080`
  - Endpoint de healthcheck
  - Pode operar com ou sem Redis

- **Redis**
  - Redis 7 (imagem alpine)
  - Armazena o contador global
  - Persistência via volume

---

## 🚀 Executando com Docker Compose

### Pré-requisitos
- Docker
- Docker Compose v2 ou superior

### Subir o ambiente

```
docker compose up -d
```
Verificar containers
```
docker compose ps
```
Acessar a aplicação
```
curl http://localhost:8080
```
Cada requisição incrementa e retorna o contador.

❤️ Healthcheck
A aplicação expõe o endpoint:
```
GET /healthz
```
Uso:
- Docker healthcheck
- Kubernetes livenessProbe / readinessProbe

Exemplo:
```
curl http://localhost:8080/healthz
```
Resposta esperada:
```
OK
```

⚙️ Variáveis de Ambiente

Aplicação (HitCounter)
| Variável         | Descrição                 | Valor padrão |
| ---------------- | ------------------------- | ------------ |
| `PORT`           | Porta HTTP da aplicação   | `8080`       |
| `ENABLE_REDIS`   | Habilita uso do Redis     | `true`       |
| `REDIS_HOST`     | Host do Redis             | `redis`      |
| `REDIS_PORT`     | Porta do Redis            | `6379`       |
| `REDIS_PASSWORD` | Senha do Redis (opcional) | vazio        |

Redis
| Variável         | Descrição                       |
| ---------------- | ------------------------------- |
| `REDIS_PASSWORD` | Senha do Redis (se configurada) |


🐳 Imagem Docker
A aplicação é distribuída como imagem Docker:
```
ghcr.io/paulovigne/hit-counter:main
```
Características:
- Build multi-stage
- Binário Go estático
- Imagem final enxuta
- Execução como usuário não-root

☸️ Kubernetes
Os manifestos Kubernetes deste repositório foram derivados diretamente do docker-compose.yml, mantendo a mesma lógica de dependências, portas, variáveis e healthchecks.

| Docker Compose | Kubernetes                         |
| -------------- | ---------------------------------- |
| `services`     | `Deployment`                       |
| `ports`        | `Service`                          |
| `environment`  | `ConfigMap`                        |
| `healthcheck`  | `livenessProbe` / `readinessProbe` |
| `volumes`      | `PersistentVolumeClaim`            |


Recursos Kubernetes Utilizados
* Deployment – hitcounter
* Deployment – redis
* Service (ClusterIP)
* ConfigMap – variáveis da aplicação
* PersistentVolumeClaim – persistência do Redis
* Ingress

Fluxo no cluster:
```
Ingress
   ↓
Service hitcounter
   ↓
Pod hitcounter
   ↓
Service redis
   ↓
Pod redis
```

🧪 Execução Local sem Redis
É possível executar a aplicação em modo totalmente stateless:
```
docker run -d \
  --name hitcounter \
  -p 8080:8080 \
  -e PORT=8080 \
  -e ENABLE_REDIS=false \
  ghcr.io/paulovigne/hit-counter:main
```
Nesse modo, o contador é mantido apenas em memória.

🧪 Execução HitCounter + Redis
```
docker network create hitcounter-net
docker volume create redis-data

docker run -d \
  --name redis \
  --network hitcounter-net \
  -v redis-data:/data \
  redis:7-alpine \
  --appendonly yes

docker run -d \
  --name hitcounter \
  --network hitcounter-net \
  -p 8080:8080 \
  -e PORT=8080 \
  -e ENABLE_REDIS=true \
  -e REDIS_HOST=redis \
  -e REDIS_PORT=6379 \
  ghcr.io/paulovigne/hit-counter:main

```
Nesse modo, o contador é mantido em disco no volume redis-data.

🎯 Objetivo do Projeto
Este projeto tem fins educacionais e demonstrativos, sendo útil para:
* Estudos de Kubernetes
* Conversão Docker Compose → Kubernetes
* Demonstração de healthchecks
* Testes de balanceamento e escalabilidade
* Exemplos de app stateless com backend stateful
