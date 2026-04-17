FROM golang:1.23-alpine AS builder

WORKDIR /app

# install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY go.mod ./
RUN go mod download

COPY . .

RUN swag init
RUN go mod tidy

RUN go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/docs/swagger.json ./docs/swagger.json
EXPOSE 8080
CMD ["./server"]
