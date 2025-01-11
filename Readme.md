
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
│
├── docker-compose.yml
|__ dockerfile
└── README.md
```

### Passo 1: Clonando o Repositório

Clone o repositório do projeto (ou crie a estrutura de diretórios conforme mostrado acima).

### Passo 2: Construindo e Executando os Contêineres

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

### Passo 3: Testando os Serviços

Você pode usar o `curl` para testar a interação entre o **Serviço A** e o **Serviço B**. Envie um request para o **Serviço A** passando um CEP válido:

```bash
curl -X POST http://localhost:8080/cep -d '{"cep":"89160222"}' -H "Content-Type: application/json"
```

Este request irá:

1. Validar o CEP.
2. Enviar o CEP para o **Serviço B** para buscar a temperatura.
3. Retornar a resposta ao cliente com os dados de temperatura e cidade.

### Passo 4: Verificando os Logs

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
