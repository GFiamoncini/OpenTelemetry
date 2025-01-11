# Etapa de construção para server-a
FROM golang:1.23 as server-a-builder

WORKDIR /app

# Copiar go.mod e go.sum para garantir que as dependências sejam baixadas
COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod download
RUN go get go.opentelemetry.io/otel
RUN go get go.opentelemetry.io/otel/exporters/zipkin
RUN go get go.opentelemetry.io/otel/sdk/trace

# Copiar o código fonte de ServerA
COPY ./ServerA ./ServerA

# Construir o servidor A para arquitetura amd64 (especificando a plataforma)
RUN GOARCH=amd64 GOOS=linux go build -o /app/server-a ./ServerA

# Garantir permissões de execução
RUN chmod +x /app/server-a

# Verificar se o binário foi criado corretamente
RUN ls -l /app

# Etapa de construção para server-b
FROM golang:1.23 as server-b-builder

WORKDIR /app

# Copiar go.mod e go.sum para garantir que as dependências sejam baixadas
COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod download
RUN go get go.opentelemetry.io/otel
RUN go get go.opentelemetry.io/otel/exporters/zipkin
RUN go get go.opentelemetry.io/otel/sdk/trace

# Copiar o código fonte de ServerB
COPY ./ServerB ./ServerB

# Construir o servidor B para arquitetura amd64 (especificando a plataforma)
RUN GOARCH=amd64 GOOS=linux go build -o /app/server-b ./ServerB

# Garantir permissões de execução
RUN chmod +x /app/server-b

# Verificar se o binário foi criado corretamente
RUN ls -l /app

# Etapa final: imagem mais leve (alpine)
FROM alpine:3.18

RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copiar binários do server-a e server-b
COPY --from=server-a-builder /app/server-a /app/server-a
COPY --from=server-b-builder /app/server-b /app/server-b

# Verificar se os binários estão sendo copiados corretamente para a imagem final
RUN ls -l /app

# Garantir permissões de execução
RUN chmod +x /app/server-a /app/server-b

# Expor as portas dos serviços
EXPOSE 8080
EXPOSE 8081

# Comando para rodar os servidores, dependendo da variável de ambiente
CMD ["/bin/sh", "-c", "echo 'Iniciando servidor $SERVICE_NAME'; if [ \"$SERVICE_NAME\" = \"server-a\" ]; then echo 'Executando server-a' && /app/server-a; elif [ \"$SERVICE_NAME\" = \"server-b\" ]; then echo 'Executando server-b' && /app/server-b; fi"]
