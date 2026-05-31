# 华安保险渠道对接服务 (insurance-bridge)

华安与下游渠道商之间的中间层：渠道调用本服务，本服务验签、解密用户三要素后转发华安，华安响应原样返回（三要素字段对渠道加密）。同时落库请求/响应及业务实体，供管理后台统计与对账。

## 功能概览

- **渠道 API**：完整实现 `ZF保险.md` 中 `/upChannelApi` 下 10 个接口
- **双密钥体系**：渠道 `channelKey`（本服务分配）↔ 华安 `huaAnKey`（华安分配），管理后台维护映射
- **三要素加解密**：渠道侧 AES-256-GCM；转发华安为明文（符合华安文档）
- **数据落库**：接口日志、用户、保单、签约、产品、付费流水
- **管理后台**：`/admin/` Web 界面 + `/admin/api` REST
- **运维**：健康检查、JSON 结构化日志、90 天日志清理与归档

## 快速启动

```bash
cd insurance-bridge
docker compose up -d --build
```

- 渠道 API：`http://localhost:8080/upChannelApi/...`
- 管理后台：`http://localhost:8080/admin/`
- 默认管理员：`admin` / `admin123`（首次启动自动创建）
- OpenAPI：`api/openapi.yaml`

## 配置

`config/config.yaml` 或通过环境变量 `BRIDGE_*` 覆盖，例如：

| 变量 | 说明 |
|------|------|
| `BRIDGE_DATABASE_DSN` | MySQL/MariaDB 连接串 |
| `BRIDGE_HUAAN_BASE_URL` | 华安域名 |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | 32 字节，三要素与库内敏感字段 |
| `BRIDGE_SECURITY_JWT_SECRET` | 管理后台 JWT |
| `BRIDGE_SERVER_UPSTREAM_TIMEOUT` | 华安调用超时，默认 25s |

## 渠道接入说明

1. 在管理后台创建渠道，获得 `channelCode` 与 `channelKey`，并配置华安提供的 `huaAnKey`。
2. 请求体与 `ZF保险.md` 一致，但 `phoneNo`、`name`、`idCard` 需使用本服务配置的 `data_encryption_key` 做 AES-256-GCM 加密后 Base64 传输。
3. 签名规则与文档一致：参数按 key ASCII 排序拼接 `k=v&...`，MD5 32 位小写；`key` 字段填渠道密钥。

签名示例见单元测试 `internal/pkg/sign/sign_test.go`（与文档示例一致）。

## 本地开发

**环境要求：Go 1.22 及以上**（`go.mod` 已声明；低于 1.22 时 `go mod tidy` 会报 `log/slog`、`slices` 不在 GOROOT）。

```bash
go version   # 确认 >= go1.22
go mod tidy
go test ./...
go run ./cmd/server
```

## 项目结构

```
cmd/server/          # 主程序与管理后台静态资源
internal/
  handler/           # HTTP 处理器
  service/           # 代理、抽取、管理服务
  repository/        # 数据访问
  model/             # 数据模型
  pkg/sign,cipher,pii
api/openapi.yaml
docker-compose.yml
```

## 数据模型

| 表 | 用途 |
|----|------|
| channels | 渠道与双密钥映射 |
| api_request_logs | 全量请求/响应（脱敏请求） |
| user_records | 用户（加密存储） |
| policy_records | 保单 |
| sign_records | 签约 |
| product_records | 产品 |
| payment_records | 价格/付费流水 |
# insurance
