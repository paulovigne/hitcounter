# Hit Counter 🧮

Aplicação simples escrita em **Go** que expõe um contador de acessos via HTTP.  
O contador pode operar de forma **stateless** ou utilizando **Redis** para persistência e contagem compartilhada entre múltiplas instâncias.

O projeto foi criado como exemplo prático de:
- Aplicação cloud-native
- Uso de Docker Compose
- Conversão de Docker Compose para manifestos Kubernetes
- Separação entre aplicação stateless e serviço stateful

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
| `depends_on`   | `readinessProbe`                   |
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

📁 Estrutura Sugerida do Repositório
.
├── docker-compose.yml
├── Dockerfile
├── hitcounter.go
├── k8s/
│   ├── hitcounter-deployment.yaml
│   ├── hitcounter-service.yaml
│   ├── redis-deployment.yaml
│   ├── redis-service.yaml
│   ├── redis-pvc.yaml
│   └── configmap.yaml
└── README.md

🧪 Execução Local sem Redis
É possível executar a aplicação em modo totalmente stateless:

```
ENABLE_REDIS=false go run hitcounter.go
```
Nesse modo, o contador é mantido apenas em memória.

🎯 Objetivo do Projeto
Este projeto tem fins educacionais e demonstrativos, sendo útil para:
* Estudos de Kubernetes
* Conversão Docker Compose → Kubernetes
* Demonstração de healthchecks
* Testes de balanceamento e escalabilidade
* Exemplos de app stateless com backend stateful
