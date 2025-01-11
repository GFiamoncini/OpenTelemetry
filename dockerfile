# Build stage
FROM golang:1.23 as builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos de configuração do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código-fonte
COPY . .

# Usa um argumento para especificar o nome do binário durante a build
ARG SERVICE_NAME
RUN go build -o /app/${SERVICE_NAME} .

# Runtime stage
FROM debian:bullseye-slim

# Define o diretório de trabalho
WORKDIR /app

# Usa o argumento para copiar o binário correspondente
ARG SERVICE_NAME
COPY --from=builder /app/${SERVICE_NAME} /app/${SERVICE_NAME}

# Expõe a porta com base no serviço
ARG SERVICE_PORT
EXPOSE ${SERVICE_PORT}

# Define a variável de ambiente do Zipkin
ENV OTEL_EXPORTER_ZIPKIN_ENDPOINT=http://zipkin:9411/api/v2/spans

# Define o comando de execução
CMD ["/bin/bash", "-c", "./${SERVICE_NAME}"]
