FROM golang:1.25-alpine AS builder

WORKDIR /app

# install ca-certificates (HTTPS)
RUN apk add --no-cache ca-certificates

# copy go files
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o auth-service ./cmd

# =========================
# Runtime stage
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/auth-service /app/auth-service

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/auth-service"]
