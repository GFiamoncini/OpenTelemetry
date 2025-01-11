# Etapa de construção para server-a e server-b
FROM golang:1.23 as builder

WORKDIR /app

# Copiar go.mod e go.sum para garantir que as dependências sejam baixadas
COPY go.mod go.sum ./

RUN go mod tidy && go mod download

# Instalar as dependências do OpenTelemetry
RUN go get go.opentelemetry.io/otel && \
    go get go.opentelemetry.io/otel/exporters/zipkin && \
    go get go.opentelemetry.io/otel/sdk/trace

# Copiar o código-fonte dos dois servidores
COPY ./ServerA ./ServerA
COPY ./ServerB ./ServerB

# Construir ambos os servidores para a arquitetura amd64
RUN GOARCH=amd64 GOOS=linux go build -o /app/server-a ./ServerA
RUN GOARCH=amd64 GOOS=linux go build -o /app/server-b ./ServerB

# Etapa final: imagem mais leve (Alpine)
FROM alpine:3.18

# Instalar dependências necessárias (libc6-compat para compatibilidade com binários)
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copiar binários dos servidores
COPY --from=builder /app/server-a /app/server-a
COPY --from=builder /app/server-b /app/server-b

# Garantir permissões de execução
RUN chmod +x /app/server-a /app/server-b

# Expor as portas dos serviços
EXPOSE 8080
EXPOSE 8081

# Comando para rodar os servidores, dependendo da variável de ambiente
CMD ["/bin/sh", "-c", "echo 'Iniciando servidor $SERVICE_NAME'; if [ \"$SERVICE_NAME\" = \"server-a\" ]; then echo 'Executando server-a' && /app/server-a; elif [ \"$SERVICE_NAME\" = \"server-b\" ]; then echo 'Executando server-b' && /app/server-b; fi"]
