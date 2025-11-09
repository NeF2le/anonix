# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="nef1le"

WORKDIR /build

COPY mapping/go.* ./mapping/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./mapping && go mod download

COPY mapping/ ./mapping/
COPY common/ ./common/
COPY .env ./.env
COPY migrations/ ./migrations/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o mapping_service mapping/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/.env ../.env
COPY --from=builder /build/mapping_service .
COPY --from=builder /build/migrations ../migrations

CMD ["./mapping_service"]