FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download
COPY . .

RUN go build -o url-shortener ./cmd/url-shortener/main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/url-shortener .

CMD [ "./url-shortener" ]