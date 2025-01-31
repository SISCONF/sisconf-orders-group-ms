FROM golang:1.23.5-alpine

WORKDIR /app

COPY ../../. /app/

RUN go build ./cmd/main.go

CMD ["./main"]