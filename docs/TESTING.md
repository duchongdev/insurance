# 测试说明

本文说明如何验证 insurance-bridge 与华安的对接，以及渠道全流程。测试分层与 [ARCHITECTURE.md](./ARCHITECTURE.md) 中的渠道层 / 华安层一致。

## 1. 测试分层

| 层级 | 测什么 | 调用对象 | 是否需要渠道 sign / PII |
|------|--------|----------|-------------------------|
| **华安直连** | 本服务 → 华安是否正常 | `huaan.Client.Call`（生产代码） | 否；body 含 `key`/`sign`（默认可为空）；PII 明文发华安 |
| **全流程** | 渠道 → 本服务 → 华安 | `POST /upChannelApi/*` | 是 |
| **单元测试** | 签名、加密、解析等 | `go test ./...` | 依用例而定 |

**原则**：华安链路测试必须走 `internal/huaan` 包的真实 `Client`，禁止在测试里复制 HTTP / 签名逻辑。

## 2. 日常单元测试

不含华安真实环境，不访问外网：

```bash
go test ./... -count=1
```

主要覆盖：

| 包 | 文件 | 内容 |
|----|------|------|
| `internal/pkg/sign` | `sign_test.go` | MD5 签名（与华安文档示例一致） |
| `internal/pkg/cipher` | `cipher_test.go` | AES-256-GCM 加解密 |
| `internal/huaan` | `client_test.go` | httptest 验证 URL、换 key、重签 |
| `internal/service` | `bank_sync_test.go` | 银行列表 JSON 解析 |

集成测试（见下节）带 `integration` build tag，**不会**被默认 `go test ./...` 执行。

## 3. 华安直连集成测试

验证本服务**真实华安对接代码**（`huaan.Client`）能否正确调用华安 10 个接口。请求体始终包含 `key` 与 `sign`：默认 `HUAAN_KEY` 为空、`sign_enabled=false` 时 `sign` 也为空字符串。

### 3.1 准备环境变量

复制模板并填入华安提供的真实参数：

```bash
cp .env.huaan.example .env.huaan
# 编辑 .env.huaan

set -a && source .env.huaan && set +a
```

| 变量 | 必填 | 说明 |
|------|------|------|
| `HUAAN_BASE_URL` | 是 | 华安域名 |
| `HUAAN_CHANNEL_CODE` | 是 | 华安侧渠道编码 |
| `HUAAN_KEY` | 否 | 对应配置 `huaan.key`，写入请求体 `key` 字段，默认空 |
| `HUAAN_API_PATH` | 否 | 默认 `/upChannelApi` |
| `HUAAN_TEST_PHONE` | 含 PII 接口 | 明文手机号 |
| `HUAAN_TEST_NAME` | 含 PII 接口 | 明文姓名 |
| `HUAAN_TEST_ID_CARD` | 含 PII 接口 | 明文身份证号 |
| `HUAAN_TEST_USER_ID` | 否 | 仅当 `proInsurance` 未返回 `userId` 时，`getSignUrl` 测试备用 |

`productCode`、`policyId` **无需配置**：集成测试分别从 `getProductInfoByChannel`、`proInsurance` 响应解析（与生产一致——`policyId` 由下游请求携带，本服务只转发）。

`.env.huaan` 已加入 `.gitignore`，**勿提交**。

### 3.2 运行

```bash
make test-huaan
```

等价于：

```bash
go test -tags=integration ./internal/huaan/ -v -count=1
```

### 3.3 行为说明

- 调用 `huaan.Client.Call`，**不经过** `ProxyService.Forward`、渠道验签、PII 加解密。
- `getBankList` **直接请求华安**，不走 Redis / `bank_info_t` 缓存。
- 华安原始响应 JSON 样例见 **[HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)**。
- 渠道调用本服务的请求/响应样例见 **[CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md)**（与华安原文成对对照）。
- 华安返回的明文 PII 在测试日志中**原样输出**（`t.Log`），不做重新加密。
- 未配置必填环境变量时，测试 `t.Skip` 跳过。
- `getProductPricesByProductCode`、`proInsurance` 的 `productCode` **自动**从 `getProductInfoByChannel` 获取（`TestHuaAnDirect_ProductFromChannel`）。
- 依赖 `policyId` 的接口（`getPolicyInfoByPolicyId`、`getProductPricesByPolicyId`、`upGradeIns`、`getSignUrl`）在测试中 **先调 proInsurance** 取 `policyId`（`TestHuaAnDirect_PolicyFromProInsurance`、`TestHuaAnDirect_GetSignUrl`）。
- `proInsurance`、`getProductPricesByProductCode`、`getProductPricesByPolicyId` 集成测试默认带 `hasSocialSecurity: "1"`（华安必填）；后两者还需三要素环境变量。

### 3.4 覆盖的接口

与 `huaan.APIPaths` 一致，共 10 个：

| 路径 | 说明 |
|------|------|
| `/proInsurance` | 投保 |
| `/upGradeIns` | 升级险种 |
| `/getSignUrl` | 获取签约链接（须 `payChannelId`，见 **3.6** 专项测试） |
| `/verifyNoCode` | 无验证码实名 |
| `/getPolicyInfoByPhoneNo` | 按手机号查保单 |
| `/getProductPricesByProductCode` | 按产品编码报价 |
| `/getBankList` | 银行列表 |
| `/getProductInfoByChannel` | 渠道产品列表 |
| `/getProductPricesByPolicyId` | 按保单报价 |
| `/getPolicyInfoByPolicyId` | 按保单 ID 查询 |

测试代码：`internal/huaan/client_integration_test.go`。

### 3.5 建议先跑通的接口

无需额外业务参数，配置好 `HUAAN_BASE_URL`、`HUAAN_CHANNEL_CODE` 即可：

- `/getBankList`
- `/getProductInfoByChannel`

含 PII 的接口需配置 `HUAAN_TEST_PHONE` / `NAME` / `ID_CARD`。

### 3.6 产品报价与投保（先 getProductInfoByChannel）

`getProductPricesByProductCode`、`proInsurance` 依赖的 `productCode` 来自 **`getProductInfoByChannel` 响应**（字段说明见 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)），与文档样例一致，例如 `ZFHLW1041001`、`ZFHLW1040003`。

```bash
go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_ProductFromChannel
```

对渠道返回的**每个**产品各跑一遍报价与投保子测试，日志中会打印选用的 `productCode` / `productName`。

### 3.7 依赖 policyId 的接口（先 proInsurance）

`getPolicyInfoByPolicyId`、`getProductPricesByPolicyId`、`upGradeIns` 的 `policyId` 来自 **`proInsurance` 成功响应的 `data.policyId`**，无需 `HUAAN_TEST_POLICY_ID`。

```bash
go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_PolicyFromProInsurance
```

流程：getProductInfoByChannel → proInsurance → 用返回的 `policyId` 调上述三个接口。

### 3.8 getSignUrl 专项测试

`policyId` 来自 **proInsurance**；`bankCode` + `payChannelId` 来自 **同一条 getBankList** 记录；`userId` 优先取自 proInsurance 响应。

| 变量 | 必填 | 说明 |
|------|------|------|
| `HUAAN_TEST_PHONE` / `NAME` / `ID_CARD` | 是 | 三要素明文 |
| `HUAAN_TEST_BANK_CODE` | 否 | 指定银行；未设置用列表第一条 |
| `HUAAN_TEST_CARD_TYPE` | 否 | 默认按银行能力推导 |
| `HUAAN_TEST_USER_ID` | 否 | 仅 proInsurance 未返回 userId 时使用 |

```bash
go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_GetSignUrl
```

实现见 `internal/huaan/client_integration_test.go`、`internal/huaan/banklist.go`、`internal/huaan/proinsurance.go`。

## 4. 全流程测试（渠道 → 本服务 → 华安）

验证渠道验签、PII 加解密与转发完整链路。

### 4.1 启动服务

```bash
./scripts/start.sh
# 或本地: go run ./cmd/server
```

### 4.2 准备渠道凭证

在管理后台创建渠道，获得：

- `channelCode`
- `channelKey`（本服务分配）
- 确认已配置 `huaAnKey`

三要素加密密钥与服务器 `security.data_encryption_key` 一致（32 字节）。

### 4.3 构造请求

```http
POST http://localhost:5051/upChannelApi/getBankList
Content-Type: application/json; charset=utf-8
```

请求体规则见 [CHANNEL_API.md](./CHANNEL_API.md)：

- `key` 填 **渠道** `channelKey`
- `sign` 用渠道密钥计算
- `phoneNo` / `name` / `idCard` 使用 AES-256-GCM 加密后 Base64

签名示例见 `internal/pkg/sign/sign_test.go`。

### 4.4 管理后台刷新银行列表

```
POST /admin/api/bank-list/refresh
Authorization: Bearer <token>
{"channelCode":"...", "huaAnKey":"..."}
```

走 `ProxyService.RefreshBankList` → `huaan.Client.Call`，成功后更新 Redis 与 `bank_info_t`。

## 5. 测试矩阵

| 场景 | 命令 / 方式 | 验证点 |
|------|-------------|--------|
| 签名算法 | `go test ./internal/pkg/sign/...` | 与华安文档示例一致 |
| 华安 HTTP 客户端 | `go test ./internal/huaan/...` | URL、key/sign 注入与 sign 开关 |
| 华安真实环境 | `make test-huaan` | 10 接口可达、响应可解析 |
| 渠道全流程 | 手动 / 脚本 POST `/upChannelApi/*` | 验签、PII、转发 |
| 银行列表缓存 | 全流程 + 管理后台 | Redis / DB 命中不重复打华安 |

## 6. 新增接口时的测试清单

- [ ] `internal/huaan/paths.go` 增加路径
- [ ] `internal/huaan/client_integration_test.go` 增加用例
- [ ] `api/openapi.yaml` 与 README 同步
- [ ] 全流程手动验证一条渠道请求
- [ ] 样例写入 [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md)，并与 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) 对照

## 7. 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) — 分层与数据流
- [CHANNEL_API.md](./CHANNEL_API.md) — 渠道商接口文档
- [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md) — 渠道请求/响应 JSON 样例
- [DEPLOY.md](./DEPLOY.md) — 服务部署与配置项
