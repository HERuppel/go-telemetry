# Go Telemetry

Este projeto conta com duas aplicações para envio e recebimento de dados em tempo real com Kafka. A aplicação que dispara eventos é feita em Golang em formato CLI, enquanto a aplicação API também feita em Golang consome, processa e armazena os eventos em um banco Mongo para consultas futuras via requisição. Ao todo, o ecossistema do projeto sobe cinco contêineres Docker: Kafka, Zookeeper, Mongo, Producer CLI e Consumer API.

## Features
### Producer
A aplicação Golang CLI (simulando um código embarcado) que dispara os eventos de diferentes tipos pré-estabelecidos, com o tempo da ocorrência e um valor randômico que simula o valor de um sensor. O disparo de eventos ocorre em uma frequência de cinco em cinco segundos. Os tipos de eventos possíveis são:
- Velocidade do veículo
- RPM do motor
- Temperatura do motor
- Nível de combustível
- Quilometragem percorrida
- Localização GPS
- Status das luzes

### Consumer
API em Golang que consome os eventos do broker Kafka, os armazena no banco não-relacional Mongo e expõe rotas para consulta dos dados armazenados via requisição HTTP. Também possui uma rota pra visualizar algumas métricas por uma data específica, agregando os eventos pelo tipo e mostrando a quantidade de eventos e a média de valor para aquele evento, e outra rota pra visualizar as métricas armazenadas desde o consumo do primeiro evento. A API está documentada com o uso da ferramenta Swagger.

## Tech

As aplicações foram feitas utilizando:

- Golang 
- Kafka e Zookeeper
- Mongo DB
- Docker
- Swagger

## Instalação

Para rodar todo o ecossistema do projeto é necessário apenas Docker e Docker Compose instalados.

1 - Criar um arquivo .env na raiz do projeto e preencher com as variáveis do arquivo .env.example, Exemplo:
```sh
MONGO_USERNAME=root
MONGO_PASSWORD=root

#Producer and Consumer vars
BROKER_ADDRESS=kafka
BROKER_PORT=9092
TOPIC_NAME=vehicle.events
MONGO_URI=mongodb://${MONGO_USERNAME}:${MONGO_PASSWORD}@mongo:27017
MONGO_DB_NAME=telemetry
MONGO_DB_COLLECTION=events
MONGO_DB_METRICS_COLLECTION=metrics
```
2 - Rodar o comando na raiz do projeto:
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
