FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/bridge ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/bridge .
COPY config/config.yaml ./config/config.yaml
COPY api/openapi.yaml ./api/openapi.yaml
ENV CONFIG_PATH=/app/config/config.yaml
EXPOSE 8080
CMD ["./bridge"]
