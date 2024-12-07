FROM golang:1.23.3-alpine as build

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o producer .

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/producer /usr/local/bin/producer

RUN apk --no-cache add ca-certificates

CMD ["producer"]