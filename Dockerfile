FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/bridge ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /out/bridge .
COPY api/openapi.yaml ./api/openapi.yaml
# 运行时通过 volume 挂载 config/config.yaml；镜像内仅保留示例供参考
COPY config/config.yaml.example ./config/config.yaml.example
ENV CONFIG_PATH=/app/config/config.yaml
EXPOSE 8080
CMD ["./bridge"]
