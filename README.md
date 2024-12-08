# Go Telemetry

Esse projeto conta com duas aplicações para envio e recebimento de dados em tempo real com Kafka. A aplicação que dispara eventos é feita em Golang em formato CLI, enquanto a aplicação que consome e processa os eventos é uma API também feita em Golang, que consome, processa e armazena os eventos em um banco MongoDB para consultas futuras via requisição. Ao todo, o ecossistema do projeto sobe cinco contêineres Docker: Kafka, Zookeeper, Mongo, Producer CLI e Consumer API.

## Features
### Producer
A aplicação Golang CLI que dispara os eventos de diferentes tipos preestabelecidos, com o tempo da ocorrência e um valor randômico que simula o valor de um sensor. O disparo de eventos ocorre numa frequência de cinco em cinco segundos.
### Consumer
API em Golang que consome os eventos do broker Kafka, os armazena no banco não-relacional MongoDB e expõe rotas para consulta dos dados armazenados via requisição HTTP. A API também possui uma rota pra visualizar algumas métricas por uma data específica, agregando os eventos pelo tipo e mostrando a quantidade de eventos e a média de valor para aquele tipo.

## Tech

As aplicações foram feitas utilizando:

- Golang 
- Kafka e Zookeeper
- Mongo DB
- Docker
- Swagger

## Instalação

Para rodar todo o ecossistema do projeto é necessário apenas Docker instalado.

1 - Preencher as variáveis de ambiente do arquivo .env que se encontra na raiz.

2 - Rodar o comando na raiz do projeto.
```sh
docker compose up --build
```
3 - Após a inicialização de todos os contêineres, a API ficará exposta em:

```sh
http://localhost:3333/
```
Ou, se preferir utilziar o Swagger:
```sh
http://localhost:3333/swagger/index.html
```
