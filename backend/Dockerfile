FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod tidy

COPY . .

RUN go build -o backend-server ./cmd/server/main.go


FROM alpine:latest  

RUN apk add --no-cache poppler-utils

WORKDIR /app

COPY credentials/mockServiceAccountKey.json credentials/mockServiceAccountKey.json

COPY migrations migrations

COPY --from=builder /app/backend-server .

EXPOSE 8080

CMD ["./backend-server"]
