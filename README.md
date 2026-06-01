# 华安保险渠道对接服务 (insurance-bridge)

华安与下游渠道商之间的中间层：渠道调用本服务，本服务验签、解密用户三要素后转发华安，华安响应原样返回（三要素字段对渠道加密）。同时落库请求/响应及业务实体，供管理后台统计与对账。

## 功能概览

- **渠道 API**：完整实现 `ZF保险.md` 中 `/upChannelApi` 下 10 个接口
- **双密钥体系**：渠道 `channelKey`（本服务分配）↔ 华安 `huaAnKey`（华安分配），管理后台维护映射
- **三要素加解密**：渠道侧 AES-256-GCM；转发华安为明文（符合华安文档）
- **数据落库**：用户、保单、签约、产品、付费流水、银行信息（`bank_info_t`）；渠道 API 调用详情写入服务日志文件
- **管理后台**：Vue 3 + Element Plus（`web/admin/`），Nginx 托管 `/admin/`，REST API 在 `/admin/api`
- **运维**：健康检查（MySQL + Redis）、Nginx 反向代理、JSON 结构化服务日志（轮转、gzip 归档、90 天自动清理）

## 快速启动

### 服务器交付部署（推荐）

解压压缩包后：

```bash
chmod +x scripts/*.sh
./scripts/install.sh          # 生成 .env 与 config/config.yaml
# 编辑 .env 与 config/config.yaml 中的密码、华安域名、32 字节密钥
./scripts/start.sh
```

详细步骤见 **[DEPLOY.md](./DEPLOY.md)**。交付方打包容器：

```bash
make package   # 输出 dist/insurance-bridge-*.tar.gz
```

### 本地 Docker

```bash
cp .env.example .env
cp config/config.yaml.example config/config.yaml
# 编辑上述文件（本地可用默认测试密码）
./scripts/start.sh
```

- 渠道 API：`http://localhost:5051/upChannelApi/...`
- 管理后台：`http://localhost:5051/admin/`（Docker/Nginx）或 `cd web/admin && npm run dev`（本地开发，API 代理至 Go）
- 默认管理员：见 `config/config.yaml` 中 `admin` 段（首次启动自动创建）
- OpenAPI：`http://localhost:5051/openapi.yaml`

## 配置

主配置为 **`.env`**（Docker 部署）与 **`config/config.yaml`**；环境变量 `BRIDGE_*` 会覆盖 yaml 中同名项。模板见 `.env.example`、`config/config.yaml.example`。

| 变量 | 说明 |
|------|------|
| `BRIDGE_DATABASE_DSN` | MySQL/MariaDB 连接串 |
| `BRIDGE_REDIS_PASSWORD` | Redis 认证密码（Docker 部署由 `REDIS_PASSWORD` 自动同步） |
| `BRIDGE_REDIS_BANK_LIST_TTL` | 银行列表缓存 TTL，默认 `0`（不过期） |
| `BRIDGE_HUAAN_BASE_URL` | 华安域名 |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | 32 字节，三要素与库内敏感字段 |
| `BRIDGE_SECURITY_JWT_SECRET` | 管理后台 JWT |
| `BRIDGE_SERVER_UPSTREAM_TIMEOUT` | 华安调用超时，默认 25s |
| `BRIDGE_LOG_LEVEL` | 服务日志级别：`debug`（测试）/ `info`（生产） |
| `log.retention_days` | 日志文件保留天数，默认 90 |
| `log.archive_enabled` | 轮转后 gzip 压缩，默认 true |
| `log.max_size_mb` | 单文件上限（MB），达到后轮转，默认 100 |

### 服务日志

- 路径：`./logs/app.log`（Docker 映射至宿主机 `logs/`）
- 渠道 API 调用：`info` 记录摘要（traceId、渠道、路径、响应码、耗时）；`debug` 记录脱敏请求与华安响应全文
- 轮转：lumberjack 按文件大小轮转，旧文件 gzip 压缩；超过 `retention_days` 自动删除
- 每日 03:00 额外扫描 `logs/` 清理遗留过期文件

## 渠道接入说明

1. 在管理后台创建渠道，获得 `channelCode` 与 `channelKey`，并配置华安提供的 `huaAnKey`。
2. 请求体与 `ZF保险.md` 一致，但 `phoneNo`、`name`、`idCard` 需使用本服务配置的 `data_encryption_key` 做 AES-256-GCM 加密后 Base64 传输。
3. 签名规则与文档一致：参数按 key ASCII 排序拼接 `k=v&...`，MD5 32 位小写；`key` 字段填渠道密钥。

签名示例见单元测试 `internal/pkg/sign/sign_test.go`（与文档示例一致）。

## 本地开发

**环境要求：Go 1.24+**；管理后台前端另需 **Node.js 18+**。

```bash
go version   # 确认 >= go1.24
go mod tidy
go test ./...

# 后端（默认 :8080）
go run ./cmd/server

# 管理后台前端（:5173，/admin/api 代理至后端）
cd web/admin && npm install && npm run dev
```

本地跨域开发时，在 `config/config.yaml` 设置 `admin.cors_origins: "http://localhost:5173"`，或环境变量 `BRIDGE_ADMIN_CORS_ORIGINS=http://localhost:5173`。

生产构建：`bash scripts/build-admin.sh` 或 `cd web/admin && npm run build`。

## 项目结构

```
cmd/server/          # Go 主程序
web/admin/           # 管理后台 Vue 前端（Element Plus）
internal/
  handler/           # HTTP 处理器
  service/           # 代理、抽取、管理服务
  repository/        # 数据访问
  model/             # 数据模型
  pkg/sign,cipher,pii,redis
api/openapi.yaml
deploy/nginx/        # Nginx：/admin/ 静态 + API 反代
docker-compose.yml
```

## 数据模型

| 表 | 用途 |
|----|------|
| channels | 渠道与双密钥映射 |
| user_records | 用户（加密存储） |
| policy_records | 保单 |
| sign_records | 签约 |
| product_records | 产品 |
| payment_records | 价格/付费流水 |
