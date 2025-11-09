# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="nef1le"

WORKDIR /build

COPY gateway/go.* ./gateway/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./gateway && go mod download

COPY gateway/ ./gateway/
COPY common/ ./common/
COPY .env ./.env

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o api_gateway gateway/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/.env ../.env
COPY --from=builder /build/api_gateway .
COPY --from=builder /build/gateway/static ./static

CMD ["./api_gateway"]