FROM golang:1.25-alpine AS base

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev

RUN go install github.com/air-verse/air@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/google/wire/cmd/wire@latest

COPY . .

EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM base AS builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -ldflags="-s -w" -o /app/bin/kali-auth ./cmd/api

FROM scratch AS prod

WORKDIR /app

COPY --from=builder /app/bin/kali-auth /app/kali-auth
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /bin/busybox /bin/busybox

EXPOSE 8080
USER 1000:1000
ENTRYPOINT ["/app/kali-auth"]
