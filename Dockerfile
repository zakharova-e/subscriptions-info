FROM golang:1.24 as builder 

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o subscriptionsService ./cmd/subscriptionsService/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /build/subscriptionsService .

CMD ["./subscriptionsService"]