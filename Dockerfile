FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o subscriptionsService ./cmd/subscriptionsService/main.go

CMD ["./subscriptionsService"]