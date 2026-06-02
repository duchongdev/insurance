# 渠道接口请求与响应样例

本文档保存**下游渠道商调用本服务**（insurance-bridge）时的 HTTP 请求与响应 JSON，供平台内部联调归档。

**对外接口说明**（字段定义、签名、错误码、推荐调用顺序）见 **[CHANNEL_API.md](./CHANNEL_API.md)**。

与 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) 的分工（平台内部）：

| 文档 | 调用方向 | 典型用途 |
|------|----------|----------|
| **本文档** | 渠道 → 本服务 | 渠道侧实测 JSON |
| [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) | 本服务 → 核心业务系统（直连） | 上游原始 JSON 对照 |

## 通用约定

### 请求

```http
POST {平台域名}/upChannelApi/{接口路径}
Content-Type: application/json; charset=utf-8
```

| 字段 | 说明 |
|------|------|
| `timestamp` | 毫秒时间戳字符串 |
| `channelCode` | 管理后台分配的渠道编码 |
| `key` | 渠道 `channelKey`（**不是**华安 `huaAnKey`） |
| `sign` | 除 `sign` 外全部参数按 key 升序拼接后 MD5（算法见 [CHANNEL_API.md](./CHANNEL_API.md)） |
| `phoneNo` / `name` / `idCard` | 若接口涉及三要素，须 **AES-256-GCM + Base64** 加密后再提交 |

签名示例见 `internal/pkg/sign/sign_test.go`。

### 响应

- HTTP 状态码一般为 `200`；网关层错误（验签失败、解密失败等）也常以 JSON 返回，见 [CHANNEL_API.md](./CHANNEL_API.md) 第 5 节。
- 业务结构与原华安文档一致；响应中三要素字段为本服务加密后的 Base64 密文。
- 与华安原文的字段差异，在各接口「与华安原文对照」小节说明。

### 脱敏

文档中 `channelKey`、三要素明文、完整密文等敏感信息以 `***` 或截断展示；完整样例仅保存在内网或本地，勿提交 Git。

## 接口索引

| 接口 | 路径 | 三要素 | 采集状态 |
|------|------|--------|----------|
| getBankList | `/upChannelApi/getBankList` | 否 | 待采集 |
| getProductInfoByChannel | `/upChannelApi/getProductInfoByChannel` | 否 | 待采集 |
| proInsurance | `/upChannelApi/proInsurance` | 是 | 待采集 |
| upGradeIns | `/upChannelApi/upGradeIns` | 否 | 待采集 |
| getSignUrl | `/upChannelApi/getSignUrl` | 否 | 待采集 |
| verifyNoCode | `/upChannelApi/verifyNoCode` | 是 | 待采集 |
| getPolicyInfoByPhoneNo | `/upChannelApi/getPolicyInfoByPhoneNo` | 是 | 待采集 |
| getProductPricesByProductCode | `/upChannelApi/getProductPricesByProductCode` | 否 | 待采集 |
| getProductPricesByPolicyId | `/upChannelApi/getProductPricesByPolicyId` | 否 | 待采集 |
| getPolicyInfoByPolicyId | `/upChannelApi/getPolicyInfoByPolicyId` | 否 | 待采集 |

以下各节按上表顺序填写；**待采集**接口保留占位，联调通过后补充 JSON 并更新索引表「采集状态」。

---

## getBankList — 获取银行列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getBankList` |
| 采集时间 | — |
| 环境 | — |
| 渠道编码 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getbanklist--获取银行列表](./HUAAN_API_SAMPLES.md#getbanklist--获取银行列表) |

### 请求体（渠道 → 本服务）

> 待采集。本接口无三要素字段；`sign` 使用渠道 `channelKey` 计算。

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "***",
  "sign": "***"
}
```

### 响应体（本服务 → 渠道）

> 待采集。`data` 结构与华安一致；若命中 Redis / `bank_info_t` 缓存，渠道侧响应与直连华安可能相同（均无三要素）。

```json
{
  "code": 200,
  "message": "操作成功",
  "data": []
}
```

### 与华安原文对照

| 差异点 | 渠道侧（本文档） | 华安原文 |
|--------|------------------|----------|
| 签名 | 渠道 `channelKey` | 华安 `huaAnKey` 或关闭签名 |
| 三要素 | 无 | 无 |
| 数据来源 | 可能来自缓存，不必然实时打华安 | 直连华安 |

### 复现采集

```bash
# 1. 启动本服务（见 docs/TESTING.md §4）
# 2. 管理后台创建渠道，获得 channelCode、channelKey
# 3. 按 CHANNEL_API.md 计算 sign 后 POST
curl -sS -X POST 'http://localhost:5051/upChannelApi/getBankList' \
  -H 'Content-Type: application/json; charset=utf-8' \
  -d '{"timestamp":"1717300000000","channelCode":"YOUR_CHANNEL_CODE","key":"***","sign":"***"}'
```

---

## getProductInfoByChannel — 渠道查询产品

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getProductInfoByChannel` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getproductinfobychannel--渠道查询产品](./HUAAN_API_SAMPLES.md#getproductinfobychannel--渠道查询产品) |

### 请求体（渠道 → 本服务）

> 待采集。

### 响应体（本服务 → 渠道）

> 待采集。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 三要素 | 无 | 无 |

---

## proInsurance — 投保

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/proInsurance` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#proinsurance--投保](./HUAAN_API_SAMPLES.md#proinsurance--投保) |

### 请求体（渠道 → 本服务）

> 待采集。`phoneNo`、`name`、`idCard` 须为加密 Base64，非明文。

### 响应体（本服务 → 渠道）

> 待采集。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 三要素 | 请求密文、响应密文 | 明文 |
| 网关错误 | 可能出现 `401`/`400` 等本服务错误码 | 无 |

---

## upGradeIns — 保单升级

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/upGradeIns` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)（对应章节） |

### 请求体 / 响应体

> 待采集。

---

## getSignUrl — 获取签约链接

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getSignUrl` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)（对应章节） |

### 请求体 / 响应体

> 待采集。

---

## verifyNoCode — 无验证码实名

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/verifyNoCode` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)（对应章节） |

### 请求体 / 响应体

> 待采集。含三要素，请求须加密。

---

## getPolicyInfoByPhoneNo — 按手机号查保单

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getPolicyInfoByPhoneNo` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)（对应章节） |

### 请求体 / 响应体

> 待采集。`phoneNo` 请求须加密；响应 `data` 内三要素为密文。

---

## getProductPricesByProductCode — 按产品编码报价

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getProductPricesByProductCode` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)（对应章节） |

### 请求体 / 响应体

> 待采集。

---

## getProductPricesByPolicyId — 按保单 ID 报价

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getProductPricesByPolicyId` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)（对应章节） |

### 请求体 / 响应体

> 待采集。

---

## getPolicyInfoByPolicyId — 按保单 ID 查询

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getPolicyInfoByPolicyId` |
| 采集时间 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getpolicyinfobypolicyid--保单id查信息](./HUAAN_API_SAMPLES.md#getpolicyinfobypolicyid--保单id查信息) |

### 请求体 / 响应体

> 待采集。华安响应 `insuredList` 等含三要素；经本服务返回后为密文。

---

## 如何新增 / 更新样例

1. 按 [TESTING.md §4](./TESTING.md#4-全流程测试渠道--本服务--华安) 启动服务并准备渠道凭证。
2. 构造符合 [CHANNEL_API.md](./CHANNEL_API.md) 的请求，保存**实际**请求与响应 JSON（脱敏后写入本文档）。
3. 同一用例在 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) 中应有对应华安原文（或注明差异原因，如缓存、网关错误）。
4. 更新上文「接口索引」中的采集状态与采集时间。

## 相关文档

- [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) — 华安上游原始 JSON（对照用）
- [CHANNEL_API.md](./CHANNEL_API.md) — 渠道商接口文档（对外）
- [TESTING.md](./TESTING.md) — 全流程测试步骤
- [ARCHITECTURE.md](./ARCHITECTURE.md) — 验签、PII、转发与缓存
