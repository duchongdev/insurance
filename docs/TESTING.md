# 测试说明

本文说明如何验证 insurance-bridge 与华安的对接，以及渠道全流程。测试分层与 [ARCHITECTURE.md](./ARCHITECTURE.md) 中的渠道层 / 华安层一致。

## 1. 测试分层

| 层级 | 测什么 | 调用对象 | 是否需要渠道 sign / PII |
|------|--------|----------|-------------------------|
| **华安直连** | 本服务 → 华安是否正常 | `huaan.Client.Call`（生产代码） | 否；PII 明文发华安 |
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

验证本服务**真实华安对接代码**（`huaan.Client`）能否正确调用华安 10 个接口。

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
| `HUAAN_KEY` | 是 | 华安签名密钥（`HuaAnKey`） |
| `HUAAN_API_PATH` | 否 | 默认 `/upChannelApi` |
| `HUAAN_TEST_PHONE` | 含 PII 接口 | 明文手机号 |
| `HUAAN_TEST_NAME` | 含 PII 接口 | 明文姓名 |
| `HUAAN_TEST_ID_CARD` | 含 PII 接口 | 明文身份证号 |
| `HUAAN_TEST_PRODUCT_CODE` | 部分接口 | 产品编码 |
| `HUAAN_TEST_POLICY_ID` | 部分接口 | 保单 ID |

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
- 华安返回的明文 PII 在测试日志中**原样输出**（`t.Log`），不做重新加密。
- 未配置必填环境变量时，测试 `t.Skip` 跳过。
- 缺少业务参数（如 `policyId`）的用例单独 `Skip`，补全环境变量后可再跑。

### 3.4 覆盖的接口

与 `huaan.APIPaths` 一致，共 10 个：

| 路径 | 说明 |
|------|------|
| `/proInsurance` | 投保 |
| `/upGradeIns` | 升级险种 |
| `/getSignUrl` | 获取签约链接 |
| `/verifyNoCode` | 无验证码实名 |
| `/getPolicyInfoByPhoneNo` | 按手机号查保单 |
| `/getProductPricesByProductCode` | 按产品编码报价 |
| `/getBankList` | 银行列表 |
| `/getProductInfoByChannel` | 渠道产品列表 |
| `/getProductPricesByPolicyId` | 按保单报价 |
| `/getPolicyInfoByPolicyId` | 按保单 ID 查询 |

测试代码：`internal/huaan/client_integration_test.go`。

### 3.5 建议先跑通的接口

无需额外业务参数，配置好三个必填变量即可：

- `/getBankList`
- `/getProductInfoByChannel`

其余接口需补充 `HUAAN_TEST_*` 业务参数。

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

请求体规则见 [CHANNEL_INTEGRATION.md](./CHANNEL_INTEGRATION.md)：

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
| 华安 HTTP 客户端 | `go test ./internal/huaan/...` | URL、换 key、sign |
| 华安真实环境 | `make test-huaan` | 10 接口可达、响应可解析 |
| 渠道全流程 | 手动 / 脚本 POST `/upChannelApi/*` | 验签、PII、转发 |
| 银行列表缓存 | 全流程 + 管理后台 | Redis / DB 命中不重复打华安 |

## 6. 新增接口时的测试清单

- [ ] `internal/huaan/paths.go` 增加路径
- [ ] `internal/huaan/client_integration_test.go` 增加用例
- [ ] `api/openapi.yaml` 与 README 同步
- [ ] 全流程手动验证一条渠道请求

## 7. 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) — 分层与数据流
- [CHANNEL_INTEGRATION.md](./CHANNEL_INTEGRATION.md) — 渠道请求格式
- [DEPLOY.md](./DEPLOY.md) — 服务部署与配置项
