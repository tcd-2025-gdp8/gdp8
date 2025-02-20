FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod tidy

COPY . .

RUN go build -o backend-server ./cmd/server/main.go


FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/backend-server .

EXPOSE 8080

CMD ["./backend-server"]
