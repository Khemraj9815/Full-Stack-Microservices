# api-gateway/Dockerfile

FROM golang:1.25.0-alpine AS builder
WORKDIR /build
# Copy the entire project root
COPY ../ .
RUN go mod download
WORKDIR /build/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["/app/server"]

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["/app/server"]