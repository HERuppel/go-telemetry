FROM golang:1.23.3-alpine AS build

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/cmd/api

RUN go build -o consumer .

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/cmd/api/consumer /usr/local/bin/consumer

RUN apk --no-cache add ca-certificates

CMD ["consumer"]