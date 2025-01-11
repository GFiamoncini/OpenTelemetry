
# Passo a Passo para Rodar a Aplicação

Este documento descreve o passo a passo para rodar a aplicação **Serviço A** e **Serviço B**, além de configurar o rastreamento distribuído com **Zipkin** usando Docker.

## Requisitos

- Docker instalado na máquina local.
- Docker Compose instalado na máquina local.

## Estrutura do Projeto

A estrutura do projeto será semelhante à seguinte:

```
/OpenTelemetry
│
├── /servico-a
│   ├── main.go
│
├── /servico-b
│   ├── main.go
│   └── Dockerfile
│
├── docker-compose.yml
|__ dockerfile
└── README.md
```

### Passo 1: Clonando o Repositório

Clone o repositório do projeto (ou crie a estrutura de diretórios conforme mostrado acima).

### Passo 2: Criando o `Dockerfile` para os Serviços

#### Dockerfile para o **Serviço A** (Localizado em `/servico-a/Dockerfile`):

```Dockerfile
# Definir a imagem base do Go
FROM golang:1.20-alpine

# Configurar o diretório de trabalho no contêiner
WORKDIR /app

# Copiar o código fonte para o diretório de trabalho
COPY . .

# Baixar as dependências
RUN go mod tidy

# Compilar o aplicativo
RUN go build -o app .

# Expor a porta que o Serviço A usará
EXPOSE 8080

# Comando para rodar o Serviço A
CMD ["./app"]
```

#### Dockerfile para o **Serviço B** (Localizado em `/servico-b/Dockerfile`):

```Dockerfile
# Definir a imagem base do Go
FROM golang:1.20-alpine

# Configurar o diretório de trabalho no contêiner
WORKDIR /app

# Copiar o código fonte para o diretório de trabalho
COPY . .

# Baixar as dependências
RUN go mod tidy

# Compilar o aplicativo
RUN go build -o app .

# Expor a porta que o Serviço B usará
EXPOSE 8081

# Comando para rodar o Serviço B
CMD ["./app"]
```

### Passo 3: Criando o Arquivo `docker-compose.yml`

O arquivo `docker-compose.yml` orquestra os contêineres do **Serviço A**, **Serviço B** e **Zipkin**. Crie este arquivo na raiz do projeto (em `/meu-projeto/docker-compose.yml`):

```yaml
version: '3.8'

services:
  servico-a:
    build:
      context: ./servico-a  # Caminho para o diretório do Serviço A
    ports:
      - "8080:8080"          # Expor a porta 8080 para o Serviço A
    depends_on:
      - servico-b            # O Serviço A depende do Serviço B
    environment:
      - OTEL_EXPORTER_ZIPKIN_ENDPOINT=http://zipkin:9411/api/v2/spans  # Configuração do Zipkin

  servico-b:
    build:
      context: ./servico-b  # Caminho para o diretório do Serviço B
    ports:
      - "8081:8081"          # Expor a porta 8081 para o Serviço B
    environment:
      - OTEL_EXPORTER_ZIPKIN_ENDPOINT=http://zipkin:9411/api/v2/spans  # Configuração do Zipkin

  zipkin:
    image: openzipkin/zipkin:latest  # Imagem do Zipkin
    ports:
      - "9411:9411"  # Expor a interface do Zipkin
```

### Passo 4: Construindo e Executando os Contêineres

1. **Na raiz do projeto**, onde o arquivo `docker-compose.yml` está localizado, execute o comando para construir as imagens e iniciar os serviços:

   ```bash
   docker-compose up --build
   ```

   O `--build` força a reconstrução das imagens baseadas nos `Dockerfile` fornecidos.

2. **Aguarde até que o Docker Compose termine de construir e iniciar os contêineres**. Quando os contêineres estiverem rodando, você verá os logs de execução na terminal.

3. **Verifique se os serviços estão funcionando**:

   - O **Serviço A** estará disponível em `http://localhost:8080`.
   - O **Serviço B** estará disponível em `http://localhost:8081`.
   - O **Zipkin** estará disponível em `http://localhost:9411`.

### Passo 5: Testando os Serviços

Você pode usar o `curl` para testar a interação entre o **Serviço A** e o **Serviço B**. Envie um request para o **Serviço A** passando um CEP válido:

```bash
curl -X POST http://localhost:8080/cep -d '{"cep":"89160222"}' -H "Content-Type: application/json"
```

Este request irá:

1. Validar o CEP.
2. Enviar o CEP para o **Serviço B** para buscar a temperatura.
3. Retornar a resposta ao cliente com os dados de temperatura e cidade.

### Passo 6: Verificando os Logs

- Para ver os logs dos contêineres enquanto os serviços estão rodando, execute:

  ```bash
  docker-compose logs -f
  ```

- Você verá os logs do **Serviço A**, **Serviço B** e **Zipkin**. Verifique se os spans estão sendo enviados corretamente para o Zipkin, e os logs de cada serviço.

### Passo 7: Parando os Contêineres

Quando terminar de testar, você pode parar os contêineres com o comando:

```bash
docker-compose down
```

Esse comando também remove os contêineres, mas mantém as imagens, permitindo que você os reinicie posteriormente sem precisar reconstruir.

---

## Considerações Finais

- O **Serviço A** e **Serviço B** estão configurados para enviar spans para o **Zipkin**, onde você pode visualizar as métricas de rastreamento distribuído.
- O **Docker Compose** facilita a orquestração entre os serviços, garantindo que tudo esteja funcionando de maneira integrada.
- Para monitorar o Zipkin, acesse `http://localhost:9411` e consulte os spans para ver o tempo de execução de cada operação entre os serviços.
