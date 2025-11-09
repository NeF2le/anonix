# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="nef1le"

WORKDIR /build

COPY auth_service/go.* ./auth_service/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./auth_service && go mod download

COPY auth_service/ ./auth_service/
COPY common/ ./common/
COPY .env ./.env
COPY migrations/ ./migrations/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o auth-service auth_service/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/.env ../.env
COPY --from=builder /build/auth-service .
COPY --from=builder /build/migrations ../migrations

CMD ["./auth-service"]