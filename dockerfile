# Etapa base: Definir a imagem base do Go
FROM golang:1.23 as builder

# Configurar o diretório de trabalho no contêiner
WORKDIR /app

# Copiar os arquivos de dependências (mod e sum)
COPY go.mod go.sum ./

# Baixar as dependências do projeto
RUN go mod download

# Copiar o restante do código fonte para o diretório de trabalho
COPY . .

# Argumento para diferenciar os serviços durante o build
ARG SERVICE_NAME

# Compilar o aplicativo Go, gerando binário com o nome do serviço
RUN go build -o /app/${SERVICE_NAME} .

# Etapa final: Contêiner para rodar os serviços
FROM alpine:3.18

# Configurar o diretório de trabalho
WORKDIR /app

# Argumento para especificar o serviço durante o build
ARG SERVICE_NAME
ARG SERVICE_PORT

# Copiar o binário gerado na etapa de build
COPY --from=builder /app/${SERVICE_NAME} /app/app

# Expor a porta para o serviço
EXPOSE ${SERVICE_PORT}

# Comando para executar o serviço
CMD ["./app"]
