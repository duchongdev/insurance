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

验证本服务**真实华安对接代码**（`huaan.Client`）能否正确调用华安 17 个接口。请求体始终包含 `key` 与 `sign`：默认 `HUAAN_KEY` 为空、`sign_enabled=false` 时 `sign` 也为空字符串。

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
| `HUAAN_TEST_SMS_CODE` | `sms/valid` | 验证码明文（映射为上游 body 字段 `code`） |
| `HUAAN_TEST_USER_ID` | 否 | 仅当 `proInsurance` 未返回 `userId` 时，`getSignUrl` 测试备用 |

`productCode`、`policyId` **无需配置**：集成测试分别从 `product/info`、`proInsurance` 响应解析（与生产一致——`policyId` 由下游请求携带，本服务只转发）。

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
- `getProductPricesByProductCode`、`proInsurance` 的 `productCode` **自动**从 `product/info` 获取（`TestHuaAnDirect_ProductFromProductInfo`）。
- 依赖 `policyId` 的接口（`getPolicyInfoByPolicyId`、`getProductPricesByPolicyId`、`getSignUrl`）在测试中 **先调 proInsurance** 取 `policyId`（`TestHuaAnDirect_PolicyFromProInsurance`、`TestHuaAnDirect_GetSignUrl`）。
- `proInsurance` 集成测试默认带 `hasSocialSecurity: 1`（整数）、`isUpgrade: 0`、`autoRenew: 1`；`getProductPricesByProductCode` 仍用 `hasSocialSecurity: "1"`（字符串）并需三要素；`getProductPricesByPolicyId` **仅需 `policyId`**。

### 3.4 覆盖的接口

与 `huaan.APIPaths` 一致，共 17 个：

| 路径 | 说明 |
|------|------|
| `/proInsurance` | 预投保（上游 `/proxy/upChannelApi/proInsurance`） |
| `/getSignUrl` | 获取签约链接（须 `payChannelId`，见 **3.6** 专项测试） |
| `/getProductPricesByProductCode` | 按产品编码报价 |
| `/getBankList` | 银行列表（上游 `/common/channel/api/getBankList`） |
| `/product/info` | 获取渠道产品信息（productCode 取自本接口） |
| `/product/info` | 获取渠道产品信息（上游 `/common/channel/api/product/info`） |
| `/sms/send` | 发送短信验证码（上游 `/common/channel/api/sms/send`） |
| `/sms/valid` | 短信验证码校验（上游 `/common/channel/api/sms/valid`） |
| `/sms/noValid` | 免短信验证码注册登录（上游 `/common/channel/api/sms/noValid`） |
| `/priceByUser` | 查询产品价格（上游 `/common/channel/api/priceByUser`） |
| `/policy/phone` | 查询用户投保情况（上游 `/common/channel/api/policy/phone`） |
| `/getUserInfoByPhoneNo` | 查询用户信息（上游 `/common/channel/api/getUserInfoByPhoneNo`，phoneNo 或 userId） |
| `/getLiabilitiesByProductId` | 查询可选责任列表（上游 `/common/channel/api/getLiabilitiesByProductId`） |
| `/getProductPricesByPolicyId` | 按保单 ID 查价格（上游 `/common/channel/api/getProductPricesByPolicyId`，仅需 policyId） |
| `/getPolicyInfoByPolicyId` | 按保单 ID 查详情（上游 `/common/channel/api/getPolicyInfoByPolicyId`，仅需 policyId） |
| `/getPhoneByToken` | 一键登录解密手机号（上游 `/common/channel/api/getPhoneByToken`，需 userInformation、token） |

测试代码：`internal/huaan/client_integration_test.go`。

### 3.5 建议先跑通的接口

无需额外业务参数，配置好 `HUAAN_BASE_URL`、`HUAAN_CHANNEL_CODE` 即可：

- `/getBankList`
- `/product/info`

含 PII 的接口需配置 `HUAAN_TEST_PHONE` / `NAME` / `ID_CARD`。

### 3.6 产品报价与投保（先 product/info）

`getProductPricesByProductCode`、`proInsurance` 依赖的 `productCode` 来自 **`product/info` 响应**（须三要素）。`getProductInfoByChannel` 已废弃，不再测试。

```bash
go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_ProductFromChannel
```

对渠道返回的**每个**产品各跑一遍报价与投保子测试，日志中会打印选用的 `productCode` / `productName`。

### 3.7 依赖 policyId 的接口（先 proInsurance）

`getPolicyInfoByPolicyId`、`getProductPricesByPolicyId` 的 `policyId` 来自 **`proInsurance` 成功响应的 `data.policyId`**，无需 `HUAAN_TEST_POLICY_ID`。

```bash
go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_PolicyFromProInsurance
```

流程：product/info → proInsurance → 用返回的 `policyId` 调上述三个接口。

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
POST http://localhost/upChannelApi/getBankList
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
| 华安真实环境 | `make test-huaan` | 17 接口可达、响应可解析 |
| 渠道全流程 | 手动 / 脚本 POST `/upChannelApi/*` | 验签、PII、转发 |
| 银行列表缓存 | 全流程 + 管理后台 | Redis / DB 命中不重复打华安 |

## 6. 新增接口时的测试清单

- [ ] `internal/huaan/paths.go` 增加路径
- [ ] `internal/huaan/client_integration_test.go` 增加用例
- [ ] `api/openapi.yaml` 与 README 同步
- [ ] 全流程手动验证一条渠道请求
- [ ] 样例写入 [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md)，并与 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) 对照
- [ ] 若暂未联调通过，登记至下文 [§7 待测试接口列表](#7-待测试接口列表)

## 7. 待测试接口列表

尚未完成**华安直连 + 渠道全流程**统一联调、或样例未采集的接口。批量联调时按本表逐项勾选；通过后移入对应样例文档并更新「采集状态 / 联调状态」。

| 接口 | 渠道路径 | 华安路径 | 三要素 | 计划环境 | 状态 | 备注 |
|------|----------|----------|--------|----------|------|------|
| 获取渠道产品信息 | `POST /upChannelApi/product/info` | `POST /common/channel/api/product/info` | `phoneNo` | `https://ins.api.hahealth.ink/` | **待测试** | 2026-06-25 直连试跑：`code=500`，`message=未授权的访问来源`；待确认该域名下 `channelCode`、IP 白名单及是否须开启签名后统一联调 |
| 获取用户短信验证码 | `POST /upChannelApi/sms/send` | `POST /common/channel/api/sms/send` | `phoneNo` | `https://ins.api.hahealth.ink/` | **待测试** | 新增接口，待与 product/info 一并联调 |
| 短信验证码校验 | `POST /upChannelApi/sms/valid` | `POST /common/channel/api/sms/valid` | `phoneNo` + `code` | `http://47.97.156.18:9040`（已测） / `https://ins.api.hahealth.ink/`（待测） | **部分完成** | 测试环境 `BLtJjF` 直连 200；上游签名为 `channelCode`+`phoneNo`+`code`+`timestamp`+`&key=`；需先 `sms/send` 取得验证码 |
| 免短信验证码注册登录 | `POST /upChannelApi/sms/noValid` | `POST /common/channel/api/sms/noValid` | `mobile` | `https://ins.api.hahealth.ink/` | **待测试** | 响应 `data` 含三要素须加密；无需 smsCode |
| 查询产品价格 | `POST /upChannelApi/priceByUser` | `POST /common/channel/api/priceByUser` | `idCard` + 业务字段 | `https://ins.api.hahealth.ink/` | **待测试** | `hasSocialSecurity` 为整数；`productPriceList` 可选 |
| 查询用户投保情况 | `POST /upChannelApi/policy/phone` | `POST /common/channel/api/policy/phone` | `phoneNo` 或 `userId` | `https://ins.api.hahealth.ink/` | **待测试** | 仅返回基础版；`phoneNo`/`userId` 二选一 |
| 查询可选责任列表 | `POST /upChannelApi/getLiabilitiesByProductId` | `POST /common/channel/api/getLiabilitiesByProductId` | `idCard` + 业务字段 | `https://ins.api.hahealth.ink/` | **待测试** | 配合 priceByUser；`productType` 1 体验版 2 正式版 |
| 查询用户信息 | `POST /upChannelApi/getUserInfoByPhoneNo` | `POST /common/channel/api/getUserInfoByPhoneNo` | `phoneNo` 或 `userId` | `https://ins.api.hahealth.ink/` | **待测试** | 响应 data 含三要素须加密 |

### 7.1 单接口复现命令（product/info）

华安直连（本服务 `huaan.Client`，PII 明文）：

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="13800138000" \
HUAAN_TEST_NAME="张三" \
HUAAN_TEST_ID_CARD="110101199001011234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/productInfo'
```

渠道全流程（验签 + PII 加解密）：

```bash
# 见 scripts/channel_sim/main.go，或 POST http://localhost/upChannelApi/product/info
```

### 7.2 单接口复现命令（sms/send）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="13800138000" \
HUAAN_TEST_NAME="张三" \
HUAAN_TEST_ID_CARD="110101199001011234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/smsSend'
```

### 7.3 单接口复现命令（sms/valid）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="13800138000" \
HUAAN_TEST_SMS_CODE="1234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/smsValid'
```

### 7.4 单接口复现命令（sms/noValid）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="13800138000" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/smsNoValid'
```

### 7.5 单接口复现命令（priceByUser）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PRODUCT_CODE="PROD2025001" \
HUAAN_TEST_ID_CARD="110101199003071234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/priceByUser'
```

### 7.6 单接口复现命令（policy/phone）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="13800138000" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/policyByPhone'
```

或使用 `userId`：

```bash
HUAAN_TEST_USER_ID="USER2025001" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/policyByPhone'
```

### 7.8 单接口复现命令（getUserInfoByPhoneNo）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="13800138000" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/getUserInfoByPhoneNo'
```

或使用 `userId`：

```bash
HUAAN_TEST_USER_ID="xxx" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/getUserInfoByPhoneNo'
```

### 7.7 单接口复现命令（getLiabilitiesByProductId）

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PRODUCT_CODE="PROD2025001" \
HUAAN_TEST_ID_CARD="110101199003071234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/liabilitiesByProductId'
```

联调通过后：

1. 更新 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) 与 [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md) 中对应接口节；
2. 将上表该行状态改为「已通过」，或从本表删除。

## 8. 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) — 分层与数据流
- [CHANNEL_API.md](./CHANNEL_API.md) — 渠道商接口文档
- [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md) — 渠道请求/响应 JSON 样例
- [DEPLOY.md](./DEPLOY.md) — 服务部署与配置项
