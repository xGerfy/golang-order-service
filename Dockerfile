FROM golang:1.25
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o /app/order-service ./cmd/server/

EXPOSE 8080

CMD ["./order-service"]