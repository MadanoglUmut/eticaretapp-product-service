FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .                

RUN go build -o product-service-app ./cmd/product/main.go

EXPOSE 5050

CMD ["./product-service-app"]

