# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="nef1le"

WORKDIR /build

COPY tokenizer/go.* ./tokenizer/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./tokenizer && go mod download

COPY tokenizer/ ./tokenizer/
COPY common/ ./common/
COPY .env ./.env
COPY certs/ ./certs/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o tokenizer_service tokenizer/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/.env ../.env
COPY --from=builder /build/certs/ ../certs
COPY --from=builder /build/tokenizer_service .

CMD ["./tokenizer_service"]