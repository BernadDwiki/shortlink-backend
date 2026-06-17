FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/shortlink-backend ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/bin/shortlink-backend /usr/local/bin/shortlink-backend
WORKDIR /app
EXPOSE 8080
ENV APP_HOST=0.0.0.0
CMD ["/usr/local/bin/shortlink-backend"]
