# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="nef1le"

WORKDIR /build

COPY mapping_cleaner/go.* ./mapping_cleaner/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./mapping_cleaner && go mod download

COPY mapping_cleaner/ ./mapping_cleaner/
COPY common/ ./common/
COPY .env ./.env

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o mapping_cleaner_service mapping_cleaner/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/.env ../.env
COPY --from=builder /build/mapping_cleaner_service .

CMD ["./mapping_cleaner_service"]