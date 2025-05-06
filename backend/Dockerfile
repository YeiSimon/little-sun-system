FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/bin/backend ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/backend .
COPY --from=builder /app/pkg/configs ./pkg/configs
EXPOSE 8080
CMD ["./backend"]