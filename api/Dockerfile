FROM public.ecr.aws/docker/library/golang:1.25.5-bookworm AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM public.ecr.aws/docker/library/debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN adduser --disabled-password --gecos '' appuser

COPY --from=go-builder /app/server /app/server
COPY --from=go-builder /app/internal/migrations /app/internal/migrations

RUN chown -R appuser:appuser /app
USER appuser
WORKDIR /app

EXPOSE 8080

CMD ["./server"]
