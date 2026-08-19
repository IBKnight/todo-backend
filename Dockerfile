# --- build stage ---
FROM golang:1.26.5-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/todo ./cmd/app

# --- run stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /app/todo .
COPY configs ./configs
COPY migrations ./migrations

USER appuser

EXPOSE 8000

CMD ["./todo"]