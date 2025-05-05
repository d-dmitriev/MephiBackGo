FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /bank-api main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /bank-api .
CMD ["./bank-api"]