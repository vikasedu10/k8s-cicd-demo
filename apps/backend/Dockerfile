FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY main.go .
# Disable CGO for static binary
ENV CGO_ENABLED=0
RUN go build -o backend main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/backend .
EXPOSE 8080
CMD ["./backend"]
