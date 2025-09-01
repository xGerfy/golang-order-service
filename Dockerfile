FROM golang:1.25
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/order-service ./cmd/server/

EXPOSE 8080

CMD ["./order-service"]