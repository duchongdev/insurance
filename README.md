# 华安保险渠道对接服务 (insurance-bridge)

> **必读 — 全局开发标准**  
> 本文档是项目的**全局要求与标准**。所有参与本项目的开发（含人工与自动化 Agent）在**动手前须先阅读并严格遵守**本文档。  
> 细则亦见 `.cursor/rules/`（应与本文档保持一致）；详细设计见 `docs/` 目录。

---

## 目录

- [第一部分：开发与项目规范（强制）](#第一部分开发与项目规范强制)
  - [1. 项目定位](#1-项目定位)
  - [2. 技术栈与环境](#2-技术栈与环境)
  - [3. 目录与模块边界](#3-目录与模块边界)
  - [4. 分层架构与依赖方向](#4-分层架构与依赖方向)
  - [5. 渠道层与华安层分离](#5-渠道层与华安层分离)
  - [6. 新增代码放哪里](#6-新增代码放哪里)
  - [7. Go 编码规范](#7-go-编码规范)
  - [8. 管理后台前端规范](#8-管理后台前端规范)
  - [9. API 与文档同步（强制）](#9-api-与文档同步强制)
  - [10. 测试规范](#10-测试规范)
  - [11. Git 提交规范](#11-git-提交规范)
  - [12. 部署与交付](#12-部署与交付)
  - [13. 安全与敏感信息](#13-安全与敏感信息)
  - [14. 改动原则](#14-改动原则)
- [第二部分：项目说明](#第二部分项目说明)
  - [功能概览](#功能概览)
  - [快速启动](#快速启动)
  - [配置](#配置)
  - [渠道接入说明](#渠道接入说明)
  - [本地开发](#本地开发)
  - [项目结构](#项目结构)
  - [文档索引](#文档索引)
  - [数据模型](#数据模型)

---

# 第一部分：开发与项目规范（强制）

## 1. 项目定位

本服务是**华安与下游渠道商之间的中间层**：

- 渠道调用本服务 `/upChannelApi/*`，使用本服务分配的 `channelKey` 验签；三要素默认 AES 加密（可按渠道配置明文）。
- 本服务验签、按渠道 PII 开关加解密、密钥替换、转发华安；当前阶段以**透明转发**为主。
- 华安上游使用华安分配的 `huaAnKey` 验签，三要素为明文。

**核心设计原则**：渠道侧契约与华安上游契约在代码层面**隔离**（见 [5. 渠道层与华安层分离](#5-渠道层与华安层分离)），便于分别测试与演进。详见 [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)。

## 2. 技术栈与环境

| 项 | 要求 |
|----|------|
| 模块名 | `github.com/huaan/insurance-bridge` |
| Go | **1.24+**（见 `go.mod`） |
| Web 框架 | Gin |
| ORM | GORM + MySQL/MariaDB |
| 缓存 | Redis |
| 管理后台 | Vue 3 + Element Plus + Vite + TypeScript（`web/admin/`） |
| API 契约 | OpenAPI 3（`api/openapi.yaml`） |
| 运行时配置 | `config/config.yaml` 或 `BRIDGE_*` 环境变量 |

**环境要求（本地开发）**：Go 1.24+；管理后台另需 Node.js 18+。

## 3. 目录与模块边界

```
cmd/server/          # 入口：依赖注入、路由注册、定时任务
internal/            # 全部业务代码（禁止在模块根建可 import 的 pkg/）
  handler/           # HTTP 处理器
  service/           # 业务编排
  repository/        # 数据访问
  model/             # GORM 实体
  huaan/             # 华安上游 HTTP 客户端（与渠道层解耦）
  middleware/        # 全站 Gin 中间件
  job/ bootstrap/    # 定时任务、首次 Seed
  pkg/               # 可复用工具：sign、cipher、pii、redis、logger
web/admin/           # 管理后台 Vue 前端
api/openapi.yaml     # HTTP 契约（单一事实来源之一）
config/              # 运行时配置模板与实例
docs/                # 架构、测试、部署、渠道接入等专题文档
scripts/             # 构建、部署、安装脚本
deploy/nginx/        # Nginx 配置
```

- **入口与组装**：仅 `cmd/server/main.go` 负责 wiring，不在 handler 中隐式初始化全局依赖。
- **业务代码均在 `internal/`**，禁止在仓库根目录创建可被外部 import 的 `pkg/`。

## 4. 分层架构与依赖方向

```
handler → service → repository → model
                ↘ internal/pkg（sign / cipher / pii / logger / redis）
                ↘ internal/huaan（华安 HTTP，由 service 编排调用）
```

| 包 | 职责 | 禁止 |
|----|------|------|
| `handler` | 读 HTTP body、调 service、写响应 | 不写业务规则、不直接操作 GORM（健康检查等例外） |
| `service` | 验签、代理编排、管理后台逻辑 | — |
| `repository` | GORM 数据访问 | 不含业务判断 |
| `model` | 持久化实体（GORM struct） | 不依赖 handler/service |
| `middleware` | 全站 HTTP 约束（trace、CORS 等） | 在 `main.go` 的 `r.Use` 注册 |
| `huaan` | 华安 HTTP 调用、路径常量 | 不含渠道验签、PII 加解密 |
| `job` / `bootstrap` | 定时任务、Seed | — |

**硬性约束**：

- `repository` / `model` **不得**依赖 `handler` / `service`。
- **禁止循环依赖**。
- `handler` 为薄层，业务逻辑放在 `service`；不在 handler 重复 service 已有逻辑。

## 5. 渠道层与华安层分离

这是本项目的**架构红线**，违反将导致测试困难与职责混乱。

| 层级 | 包 / 入口 | 职责 |
|------|-----------|------|
| **渠道层** | `handler/channel.go` → `service/proxy.go` | 校验 `channelCode` / `channelKey` / `sign`；从 Redis 读 `piiEncrypted`；`getBankList` 多级缓存；PII 加解密 |
| **华安层** | `internal/huaan` | 写入 `huaAnKey` 与 `sign`；POST 华安；返回**原始 JSON**（PII 明文） |
| **编排** | `ProxyService.Forward` | 验签 → 读渠道缓存 → [缓存] → [PII 解密] → `huaan.Call` → [写缓存] → [PII 加密] → 返回 |

**接口路径单一维护**：所有渠道 POST 路径定义在 `internal/huaan/paths.go` 的 `APIPaths`；`ChannelHandler.Register` 遍历该列表注册路由；华安上游路径差异通过 `upstreamPathOverrides` 映射。

**测试分层**（详见 [10. 测试规范](#10-测试规范)）：

- 华安链路测试必须走 `huaan.Client`，**禁止**在测试中复制 HTTP / 签名逻辑。
- 华安集成测试**不经过** `ProxyService.Forward`、渠道验签、PII 加解密。

**渠道配置缓存**（Redis 键 `bridge:channel:{channelCode}`）：

- 新建渠道默认 `piiEncrypted = true`（须加密）。
- 管理后台渠道 CRUD 后**同步**更新/删除 Redis 缓存。
- 服务启动时从 `channels` 表预热；`piiEncrypted` **以缓存为准**。

## 6. 新增代码放哪里

| 需求 | 放置位置 |
|------|----------|
| 全站 HTTP 约束 | `internal/middleware` + `cmd/server/main.go` |
| 仅管理后台鉴权/约束 | `handler/admin.go` 路由组中间件 |
| 渠道 API 共用规则 | `service/proxy.go`（或 `ProxyService` 私有方法） |
| 新增华安接口路径 | `internal/huaan/paths.go` + handler 自动注册 |
| 华安 HTTP 细节 | `internal/huaan/` |
| 表级约束 | `model` GORM tag 或迁移 SQL |
| 可复用工具（签名、加密） | `internal/pkg/<name>/` |
| 管理后台页面/API 调用 | `web/admin/src/` |

## 7. Go 编码规范

### 7.1 注释（必要且克制）

| 必须注释 | 说明 |
|---------|------|
| 每个 `package` | 一行中文说明包职责 |
| 导出类型 / 函数 / 方法 | 中文 godoc：做什么、关键参数/返回值 |
| 非显而易见的业务规则 | 验签顺序、三要素加解密、key 替换、错误码含义 |
| 魔法值、状态码、表字段语义 | 行尾或块注释说明「为什么」 |
| 复杂分支或算法 | 简要说明意图 |

**不必**注释：getter/setter、一眼能看懂的 CRUD、标准库用法。

### 7.2 风格

- 构造函数：`NewXxx` / `NewXxxRepo` / `NewXxxHandler`。
- 仓储错误：`var ErrXxx = errors.New("...")` 放在对应 `repository` 包。
- 日志：`go.uber.org/zap`；渠道/代理链路保留 `trace_id` 字段。
- import 顺序：标准库 → 第三方 → 本项目模块。
- 渠道 API 响应格式遵循华安文档（多数 `HTTP 200` + body.code）。
- 管理后台 JSON 字段与 `model` 的 json tag 保持一致。

### 7.3 各层写法示例

```go
// handler：薄层，委托 service
func (h *ChannelHandler) handle(c *gin.Context, apiPath string) {
    out, status, _ := h.proxy.Forward(c.Request.Context(), apiPath, raw)
    c.Data(status, "application/json; charset=utf-8", out)
}

// service：业务规则与编排
// repository：仅 GORM 操作，返回 model 或 error
```

### 7.4 避免

- 在 `handler` 重复 `service` 已有逻辑。
- 为「可能以后用到」提前抽 interface（除非单测需要 mock）。
- 修改 `go.mod` 中 Go 版本低于 1.24。

## 8. 管理后台前端规范

| 项 | 约定 |
|----|------|
| 目录 | `web/admin/` |
| 技术栈 | Vue 3 + Element Plus + Pinia + Vue Router + TypeScript |
| API 调用 | `web/admin/src/api/`，统一经 `request.ts` 封装 axios |
| 路由 | `web/admin/src/router/index.ts` |
| 页面 | `web/admin/src/views/`；共用组件放 `components/` |
| 本地开发 | `cd web/admin && npm run dev`（`:5173`，API 代理至后端） |
| 生产构建 | `bash scripts/build-admin.sh` 或 `npm run build`，产物由 Nginx 托管 |

**约束**：

- 管理后台 API 变更须同步 `api/openapi.yaml` 与后端 `handler/admin.go`。
- JSON 字段命名与后端 `model` json tag 一致。
- 不引入与 Element Plus 功能重复的重型 UI 库。

## 9. API 与文档同步（强制）

### 9.1 单一事实来源

| 类型 | 文件 |
|------|------|
| HTTP 契约 | `api/openapi.yaml` |
| 全局标准与使用说明 | **本文档 `README.md`** |
| 下游渠道商接口 | `docs/CHANNEL_API.md` |
| 实现 | `internal/handler/channel.go`、`internal/handler/admin.go` |

### 9.2 规则：改 OpenAPI → 必须同步 README

修改 `api/openapi.yaml` 时，**同一任务内**更新 README 对应内容：

| OpenAPI 变更 | README 需同步 |
|-------------|--------------|
| 新增/删除/改名路径 | 「功能概览」「快速启动」URL 说明 |
| 管理后台 API 变动 | 管理后台相关描述 |
| 鉴权方式变更 | 「渠道接入说明」 |
| 请求体/三要素/签名规则 | 「渠道接入说明」 |
| `info.version` 或重大行为变更 | 「功能概览」 |

**禁止**只改 OpenAPI 而不更新 README（纯 typo/格式修正且不影响语义除外）。

### 9.3 规则：接口实现变动 → 必须同步 OpenAPI

以下任一改动，**同一任务内**更新 `api/openapi.yaml`：

- `handler/channel.go` 增删改路径或 HTTP 方法
- `handler/admin.go` 增删改路由、参数、响应结构
- `cmd/server/main.go` 新增/变更对外 HTTP 路径
- `service` 层导致对外 JSON 字段、错误码语义、鉴权方式变化

同步检查项：

- [ ] `paths` 与代码一致
- [ ] `requestBody` / 参数与 handler 绑定一致
- [ ] `components/schemas` 反映实际 JSON 字段
- [ ] 渠道接口标注 `tags: [渠道接口]`；管理接口标注管理相关 tag
- [ ] 破坏性变更 bump `info.version`

### 9.4 同一 commit 内完成

OpenAPI、README、handler 的接口相关改动应放在**同一提交**，便于 review。

**对照关系**：

```
handler/channel.go apiRoutes  ↔  openapi paths /upChannelApi/*
handler/admin.go Register     ↔  openapi paths /admin/api/*
main.go 路由                  ↔  openapi /health/* 等
README 快速启动 URL           ↔  openapi servers + paths
huaan.APIPaths                ↔  渠道路由注册列表
```

## 10. 测试规范

### 10.1 测试分层

| 层级 | 测什么 | 调用对象 | 渠道 sign / PII |
|------|--------|----------|-----------------|
| **单元测试** | 签名、加密、解析等 | `go test ./...` | 依用例 |
| **华安直连** | 本服务 → 华安 | `huaan.Client.Call` | 否；PII 明文 |
| **全流程** | 渠道 → 本服务 → 华安 | `POST /upChannelApi/*` | 是 |

### 10.2 日常命令

```bash
go test ./... -count=1    # 或 make test
make test-huaan           # 华安直连集成测试（需 .env.huaan）
```

- 集成测试带 `integration` build tag，**不会**被默认 `go test ./...` 执行。
- 单测文件：`<pkg>_test.go`，与实现同包；签名/加密等纯函数优先补测试。
- **提交前**须通过 `go test ./...`。
- 不新增无意义的测试（仅验证编译通过、恒真断言等）。

### 10.3 华安集成测试

- 环境变量见 `docs/TESTING.md` 与 `.env.huaan.example`。
- `.env.huaan` 已加入 `.gitignore`，**勿提交**。
- 华安原始响应样例见 `docs/HUAAN_API_SAMPLES.md`。

## 11. Git 提交规范

### 11.1 提交前 Code Review（强制）

执行 `git commit` 前须走 **pre-commit-review**（`.cursor/skills/pre-commit-review/SKILL.md`）：

- 满分 100，**及格 60**。
- 不及格**禁止提交**，须按报告修复后重新 review。

### 11.2 原子性：一功能 / 一问题一提交

同一 commit 只包含解决**同一个 bug** 或完成**同一项功能**的改动；无关文件不要混在同一提交。

示例：

- `fix(proxy): 修复渠道验签失败时未记录 trace`
- `feat(admin): 管理后台支持按渠道筛选`
- `feat(api): 新增 xxx 接口并更新 openapi 与 README`
- `chore: 升级 gin 至 v1.10.0`

### 11.3 提交信息

- **标题与说明优先中文**；模块名、路径、命令可中英混写。
- 改 OpenAPI 时同步的 README 属于同一功能，可放在同一 commit。

### 11.4 提交前检查

```bash
git status          # 确认仅含本次原子变更
git diff            # 确认 diff 范围
go test ./...       # 测试通过
# pre-commit-review 通过
```

## 12. 部署与交付

### 12.1 测试环境

| 项 | 值 |
|----|-----|
| 地址 | `10.41.61.41` |
| 用户 | `root` |
| 部署目录 | `/home/ins/` |

### 12.2 部署流程

**任何部署前必须先停服**：

```bash
ssh root@10.41.61.41 'cd /home/ins && ./scripts/stop.sh'
```

推荐流程：

1. 本机打包：`make package`
2. 停服（见上）
3. 上传压缩包至 `/home/ins/`，保留服务器 `.env` 与 `config/config.yaml`
4. 解压、合并配置
5. 启动：`cd /home/ins && ./scripts/start.sh`
6. 验证：`curl http://127.0.0.1/health/ready`

快捷命令：`make deploy-test`（见 `scripts/deploy-test.sh`）。

### 12.3 可交付改动后的默认动作

完成功能实现、bug 修复或管理后台/接口改造后，**同一任务内**执行 `make deploy-test`，并确认：

- 健康检查 `ready OK`
- 管理后台可访问：`http://10.41.61.41/admin/`

**可跳过部署**：仅改文档/注释且无运行时影响；用户明确说「先不要部署」。

详细步骤见 [docs/DEPLOY.md](./docs/DEPLOY.md)。

## 13. 安全与敏感信息

| 规则 | 说明 |
|------|------|
| 密钥不入库 | `.env`、`.env.huaan`、真实密码/密钥**不得**提交 Git |
| 日志不记 body | 渠道 API 调用仅记摘要（traceId、渠道、路径、响应码、耗时），**不记录**请求/响应 body |
| PII 默认加密 | 新建渠道默认 `piiEncrypted=true` |
| 配置模板 | 仅提交 `.env.example`、`config/config.yaml.example` |
| 密钥长度 | `data_encryption_key` 须 32 字节；JWT secret 须足够随机 |

## 14. 改动原则

1. **最小 diff**：只改与任务相关的文件，不顺手重构。
2. **不过度抽象**：不为「可能以后用到」提前抽层或 interface。
3. **遵循现有约定**：新增代码应像同一作者所写，匹配命名、结构与文档风格。
4. **不新增无请求的测试、文档或抽象**（接口相关文档同步除外）。
5. **注释解释「为什么」**，不重复代码字面意思。
6. **敏感配置与密钥不得写入仓库**。

---

# 第二部分：项目说明

## 功能概览

- **渠道 API**：完整实现华安 `/upChannelApi` 下 10 个接口，并新增 `/common/channel/api/*` 扩展接口（产品信息、短信、查价、投保情况等）
- **双密钥体系**：渠道 `channelKey`（本服务分配）↔ 华安 `huaAnKey`（华安分配），管理后台维护映射
- **三要素加解密**：默认渠道侧 AES-256-GCM；可按渠道关闭（明文对接）。是否加密由 Redis 渠道配置缓存决定，管理后台创建/修改/删除渠道时同步更新缓存；**新建渠道默认须加密**
- **银行列表**：`getBankList` 转发华安 `/common/channel/api/getBankList`；支持 Redis → `bank_info_t` → 华安 多级缓存；管理后台可查询与手动刷新
- **管理后台**：Vue 3 + Element Plus（`web/admin/`），渠道配置（含回调地址）、银行列表、管理员登录
- **投保结果回调**：华安 `POST /huaan/callback/insureNotify` → 本服务转发至渠道 `callbackUrl`
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

详细步骤见 **[docs/DEPLOY.md](./docs/DEPLOY.md)**。更多文档见 **[docs/README.md](./docs/README.md)**。交付方打包容器：

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

- 渠道 API：`http://localhost/upChannelApi/...`（或配置 SSL 后 `https://localhost/...`）
- 管理后台：`http://localhost/`（Docker/Nginx）或 `cd web/admin && npm run dev`（本地开发）
- 默认管理员：见 `config/config.yaml` 中 `admin` 段（首次启动自动创建）
- OpenAPI：`http://localhost/openapi.yaml`

## 配置

主配置为 **`.env`**（Docker 部署）与 **`config/config.yaml`**；环境变量 `BRIDGE_*` 会覆盖 yaml 中同名项。模板见 `.env.example`、`config/config.yaml.example`。

| 变量 | 说明 |
|------|------|
| `BRIDGE_DATABASE_DSN` | MySQL/MariaDB 连接串 |
| `BRIDGE_REDIS_PASSWORD` | Redis 认证密码 |
| `BRIDGE_REDIS_BANK_LIST_TTL` | 银行列表缓存 TTL，默认 `0`（不过期） |
| `BRIDGE_HUAAN_BASE_URL` | 华安域名 |
| `BRIDGE_HUAAN_KEY` | 请求华安时 body 中的 `key` 字段 |
| `BRIDGE_HUAAN_SIGN_ENABLED` | 是否按华安规则生成 `sign`；默认 `false` |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | 32 字节，渠道三要素 AES 加解密 |
| `BRIDGE_SECURITY_JWT_SECRET` | 管理后台 JWT |
| `BRIDGE_SERVER_UPSTREAM_TIMEOUT` | 华安调用超时，默认 25s |
| `BRIDGE_LOG_LEVEL` | 服务日志级别：`debug` / `info` |
| `log.retention_days` | 日志文件保留天数，默认 90 |
| `log.archive_enabled` | 轮转后 gzip 压缩，默认 true |
| `log.max_size_mb` | 单文件上限（MB），默认 100 |

### 服务日志

- 路径：`./logs/app.log`（Docker 映射至宿主机 `logs/`）
- 渠道 API：`info` 记录摘要，**不记录**请求/响应 body
- 轮转：lumberjack 按文件大小轮转，旧文件 gzip 压缩；超过 `retention_days` 自动删除
- 每日 03:00 额外扫描 `logs/` 清理遗留过期文件

## 渠道接入说明

1. 在管理后台创建渠道，获得 `channelCode` 与 `channelKey`，并配置华安提供的 `huaAnKey` 及**投保结果回调地址** `callbackUrl`（必填）。**新建渠道默认开启三要素加密**（`piiEncrypted=true`）。
2. 向运营确认该渠道的 `piiEncrypted` 配置：为 `true` 时，`phoneNo`、`name`、`idCard` 须使用 `data_encryption_key` 做 AES-256-GCM 加密后 Base64 传输；为 `false` 时可传明文（须与平台约定，通常仅联调/内网）。
3. 签名规则：参数按 key ASCII 排序拼接 `k=v&...`，MD5 32 位小写；`key` 字段填渠道密钥。签名参与计算的 PII 值须与请求 JSON 中实际提交一致（密文或明文）。

方案细节见 **[docs/PII_CHANNEL_ENCRYPTION.md](./docs/PII_CHANNEL_ENCRYPTION.md)**；渠道商接口文档见 **[docs/CHANNEL_API.md](./docs/CHANNEL_API.md)**。

## 本地开发

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
web/admin/           # 管理后台 Vue 前端
internal/
  handler/           # HTTP 处理器
  service/           # 代理、管理服务
  huaan/             # 华安上游 HTTP 客户端
  repository/        # 数据访问
  model/             # 数据模型
  pkg/sign,cipher,pii,redis,logger
docs/                # 专题文档
api/openapi.yaml
deploy/nginx/        # Nginx：/admin/ 静态 + API 反代
docker-compose.yml
```

## 文档索引

| 文档 | 说明 |
|------|------|
| **README.md（本文档）** | **全局开发标准 + 项目说明** |
| [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) | 渠道层 / 华安层架构、数据流 |
| [docs/TESTING.md](./docs/TESTING.md) | 单元测试、华安直连集成测试 |
| [docs/DEPLOY.md](./docs/DEPLOY.md) | Docker 交付部署 |
| [docs/CHANNEL_API.md](./docs/CHANNEL_API.md) | 下游渠道商接口文档 |
| [docs/PII_CHANNEL_ENCRYPTION.md](./docs/PII_CHANNEL_ENCRYPTION.md) | 三要素加解密方案 |
| [docs/README.md](./docs/README.md) | 文档目录 |

## 数据模型

| 表 | 用途 |
|----|------|
| `channels` | 渠道与双密钥映射 |
| `admin_users` | 管理后台登录账号 |
| `bank_info_t` | 银行列表（getBankList 同步） |

已有部署若仍存在历史表（如 `user_records`、`policy_records`、`api_request_logs` 等），代码已不再使用，可按需手动 `DROP TABLE` 清理。
