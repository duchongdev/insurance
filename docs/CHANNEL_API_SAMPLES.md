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
| getBankList | `/upChannelApi/getBankList` | 是 | 2026-06-25（路径/字段对齐规范） |
| getProductInfoByChannel | `/upChannelApi/getProductInfoByChannel` | 否 | 待采集 |
| product/info | `/upChannelApi/product/info` | 是 | **待测试** |
| sms/send | `/upChannelApi/sms/send` | 是 | **待测试** |
| sms/valid | `/upChannelApi/sms/valid` | 是 | **待测试** |
| sms/noValid | `/upChannelApi/sms/noValid` | 是 | **待测试** |
| priceByUser | `/upChannelApi/priceByUser` | 是 | **待测试** |
| policy/phone | `/upChannelApi/policy/phone` | 是 | **待测试** |
| getUserInfoByPhoneNo | `/upChannelApi/getUserInfoByPhoneNo` | 是 | **待测试** |
| getLiabilitiesByProductId | `/upChannelApi/getLiabilitiesByProductId` | 是 | **待测试** |
| proInsurance | `/upChannelApi/proInsurance` | 是 | 2026-06-02（华安样例已对齐） |
| upGradeIns | `/upChannelApi/upGradeIns` | 否 | 待采集 |
| getSignUrl | `/upChannelApi/getSignUrl` | 否 | 待采集 |
| verifyNoCode | `/upChannelApi/verifyNoCode` | 是 | 待采集 |
| getPolicyInfoByPhoneNo | `/upChannelApi/getPolicyInfoByPhoneNo` | 是 | 待采集 |
| getProductPricesByProductCode | `/upChannelApi/getProductPricesByProductCode` | 是 | 2026-06-02（华安样例已对齐） |
| getProductPricesByPolicyId | `/upChannelApi/getProductPricesByPolicyId` | 否 | 2026-06-25（路径/字段对齐规范） |
| getPolicyInfoByPolicyId | `/upChannelApi/getPolicyInfoByPolicyId` | 否（响应含 PII） | 2026-06-25（路径/字段对齐规范） |

以下各节按上表顺序填写；**待采集**接口保留占位，联调通过后补充 JSON 并更新索引表「采集状态」。

---

## getBankList — 获取支持的银行列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getBankList`（转发华安 `/common/channel/api/getBankList`） |
| 采集时间 | — |
| 环境 | — |
| 渠道编码 | — |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getbanklist--获取支持的银行列表](./HUAAN_API_SAMPLES.md#getbanklist--获取支持的银行列表) |

### 请求体（渠道 → 本服务）

本接口无三要素；请求体含公共签名字段，华安侧仅需 `channelCode`（本服务转发时补齐 timestamp/key/sign）。

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "***",
  "sign": "***"
}
```

### 响应体（本服务 → 渠道）

```json
{
  "code": 200,
  "message": "成功",
  "data": [
    {
      "bankName": "中国工商银行",
      "bankCode": "ICBC",
      "debitCard": "1",
      "creditCard": "1",
      "payChannelId": "xxx",
      "isActBank": 1
    }
  ]
}
```

若命中 Redis，为华安原始 JSON（可能含历史字段如 `sortOrder`）；若仅命中 `bank_info_t`，返回上述规范字段。

### 与华安原文对照

| 差异点 | 渠道侧（本文档） | 华安原文 |
|--------|------------------|----------|
| 华安 URL | 本服务拼 `/common/channel/api/getBankList` | 同左 |
| 签名 | 渠道 `channelKey` | 华安 `huaAnKey` 或关闭签名 |
| 三要素 | 无 | 无 |
| 数据来源 | 可能来自缓存，不必然实时打华安 | 直连华安 |

### 复现采集

```bash
# 1. 启动本服务（见 docs/TESTING.md §4）
# 2. 管理后台创建渠道，获得 channelCode、channelKey
# 3. 按 CHANNEL_API.md 计算 sign 后 POST
curl -sS -X POST 'http://localhost/upChannelApi/getBankList' \
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

## product/info — 获取渠道产品信息

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/product/info` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/product/info` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#productinfo--获取渠道产品信息](./HUAAN_API_SAMPLES.md#productinfo--获取渠道产品信息) |

### 请求体（渠道 → 本服务）

> 待测试。须带公共字段 + `phoneNo`（三要素密文）。

### 响应体（本服务 → 渠道）

> 待测试。预期 `data` 含 `productCode`、`productName`、`companyName`、`productType`、`channelCode`。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 三要素 | `phoneNo` 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/product/info` | `/common/channel/api/product/info`（非 `/upChannelApi` 前缀） |

---

## sms/send — 获取用户短信验证码

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/sms/send` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/sms/send` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#smssend--获取用户短信验证码](./HUAAN_API_SAMPLES.md#smssend--获取用户短信验证码) |

### 请求体（渠道 → 本服务）

> 待测试。须带公共字段 + `phoneNo`（三要素密文）。

### 响应体（本服务 → 渠道）

> 待测试。预期 `data` 为字符串，如 `"验证码发送成功"`。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 三要素 | `phoneNo` 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/sms/send` | `/common/channel/api/sms/send` |

---

## sms/valid — 短信验证码校验

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/sms/valid` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/sms/valid` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#smsvalid--短信验证码校验](./HUAAN_API_SAMPLES.md#smsvalid--短信验证码校验) |

### 请求体（渠道 → 本服务）

> 待测试。须带公共字段 + `mobile`（加密）+ `smsCode`（明文）。

### 响应体（本服务 → 渠道）

> 待测试。预期 `data` 含 `userId`、`idCard`、`name`、`phoneNo`、`channelCode`（三要素字段为密文）。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 手机号字段 | 请求 `mobile` 须加密 | 明文 `mobile` |
| 三要素响应 | `data` 内 PII 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/sms/valid` | `/common/channel/api/sms/valid` |

---

## sms/noValid — 免短信验证码注册登录

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/sms/noValid` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/sms/noValid` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#smsnovalid--免短信验证码注册登录](./HUAAN_API_SAMPLES.md#smsnovalid--免短信验证码注册登录) |

### 请求体（渠道 → 本服务）

> 待测试。须带公共字段 + `mobile`（加密）。

### 响应体（本服务 → 渠道）

> 待测试。预期 `data` 含 `userId`、`idCard`、`name`、`phoneNo`、`channelCode`（三要素字段为密文）。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 手机号字段 | 请求 `mobile` 须加密 | 明文 `mobile` |
| 三要素响应 | `data` 内 PII 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/sms/noValid` | `/common/channel/api/sms/noValid` |

---

## priceByUser — 查询产品价格

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/priceByUser` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/priceByUser` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#pricebyuser--查询产品价格](./HUAAN_API_SAMPLES.md#pricebyuser--查询产品价格) |

### 请求体（渠道 → 本服务）

> 待测试。须带 `productCode`、`idCard`（加密）、`hasSocialSecurity`（整数）；`productPriceList` 可选。

### 响应体（本服务 → 渠道）

> 待测试。预期 `data` 含产品价格与生效天数等字段（无三要素）。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 身份证 | 请求 `idCard` 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/priceByUser` | `/common/channel/api/priceByUser` |

---

## policy/phone — 查询用户投保情况

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/policy/phone` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/policy/phone` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#policyphone--查询用户投保情况](./HUAAN_API_SAMPLES.md#policyphone--查询用户投保情况) |

### 请求体（渠道 → 本服务）

> 待测试。`phoneNo`（加密）与 `userId` 至少传一项。

### 响应体（本服务 → 渠道）

> 待测试。`data` 为基础版产品数组；`insuredList` 等 PII 字段为密文。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 手机号 | 请求 `phoneNo` 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/policy/phone` | `/common/channel/api/policy/phone` |

---

## getUserInfoByPhoneNo — 查询用户信息

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getUserInfoByPhoneNo` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getUserInfoByPhoneNo` |
| 采集时间 | 2026-06-25 |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getuserinfobyphoneno--查询用户信息](./HUAAN_API_SAMPLES.md#getuserinfobyphoneno--查询用户信息) |

### 请求体（渠道 → 本服务）

`phoneNo`（加密）与 `userId` 至少传一项。

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "***",
  "sign": "***",
  "phoneNo": "AES-GCM-Base64-密文"
}
```

### 响应体（本服务 → 渠道）

结构与华安一致；`phoneNo`/`name`/`idCard` 为 **AES 密文**（示例为脱敏明文）：

```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "userId": "xxx",
    "idCard": "xxx",
    "name": "xxx",
    "phoneNo": "xxx",
    "channelCode": "xxx"
  }
}
```

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 手机号 | 请求 `phoneNo` 须加密 | 明文 |
| 三要素 | 响应须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/getUserInfoByPhoneNo` | `/common/channel/api/getUserInfoByPhoneNo` |

---

## getLiabilitiesByProductId — 查询可选责任列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getLiabilitiesByProductId` |
| 华安路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getLiabilitiesByProductId` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getliabilitiesbyproductid--查询可选责任列表](./HUAAN_API_SAMPLES.md#getliabilitiesbyproductid--查询可选责任列表) |

### 请求体（渠道 → 本服务）

> 待测试。须带 `productCode`、`idCard`（加密）、`hasSocialSecurity`、`productType`（1 体验版 2 正式版）。

### 响应体（本服务 → 渠道）

> 待测试。`data` 为责任项数组，可组装为 priceByUser 的 `productPriceList`。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 签名 | 渠道密钥 | 华安密钥或空 |
| 身份证 | 请求 `idCard` 须加密 | 明文 |
| 上游路径 | 本服务 `/upChannelApi/getLiabilitiesByProductId` | `/common/channel/api/getLiabilitiesByProductId` |

---

## proInsurance — 预投保

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/proInsurance`（本服务转发华安 `POST /proxy/upChannelApi/proInsurance`） |
| 采集时间 | 2026-06-02（字段 2026-06-25 对齐规范） |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#proinsurance--预投保](./HUAAN_API_SAMPLES.md#proinsurance--预投保) |

姓名、证件号录入后发起预投保；`code=200` 则继续投保。`productCode` 通常来自 `getProductInfoByChannel` 或产品查询类接口。

### 请求体（渠道 → 本服务）

与华安字段相同，三要素为密文：

```json
{
  "timestamp": "1780456222458",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "YOUR_CHANNEL_KEY",
  "sign": "YOUR_SIGN",
  "productCode": "PROD2025001",
  "phoneNo": "Base64密文(手机号)",
  "name": "Base64密文(姓名)",
  "idCard": "Base64密文(身份证号)",
  "hasSocialSecurity": 1,
  "isUpgrade": 0,
  "autoRenew": 1,
  "productPriceList": [
    {
      "liabilityName": "罕见保险金",
      "price": 0.0,
      "kindCode": "102",
      "uwCount": 10
    }
  ]
}
```

### 响应体（本服务 → 渠道）

规范成功响应（华安可能额外返回 `userId`、`policyStatus` 等，本服务原样转发）：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "policyId": "xxx",
    "policyNo": "xxx"
  }
}
```

历史实测（旧环境）仍可见 `policyStatus`、`userId`，见 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md)。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 三要素 | 请求须密文 | 明文（直连华安） |
| 华安 URL | 本服务拼 `/proxy/upChannelApi/proInsurance` | 同左（非 `/upChannelApi/proInsurance`） |
| `hasSocialSecurity` / `isUpgrade` / `autoRenew` | 整数 | 整数 |
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
| 采集时间 | 2026-06-02 |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getproductpricesbyproductcode--查询产品价格](./HUAAN_API_SAMPLES.md#getproductpricesbyproductcode--查询产品价格) |

`productCode` 来自 `getProductInfoByChannel`；须带 `hasSocialSecurity` 与三要素。

### 请求体（渠道 → 本服务）

```json
{
  "timestamp": "1780456222050",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "YOUR_CHANNEL_KEY",
  "sign": "YOUR_SIGN",
  "productCode": "ZFHLW1041001",
  "hasSocialSecurity": "1",
  "phoneNo": "Base64密文(手机号)",
  "name": "Base64密文(姓名)",
  "idCard": "Base64密文(身份证号)"
}
```

### 响应体（本服务 → 渠道）

`data` 数组与华安一致（无 PII 字段）：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "productCode": "ZFHLW1041001",
      "productId": "cc55f8cff47a4f0d88da3a9ca057772a",
      "price": 0.65,
      "productName": "百万医疗险-体验版",
      "productType": "1"
    },
    {
      "productCode": "ZFBWYL1038002",
      "productId": "fade25f092f44d2088fe8d7eeba4d493",
      "price": 150.8,
      "productName": "百万医疗险-正式版",
      "productType": "2"
    }
  ]
}
```

`productCode=ZFHLW1040003` 时响应见 [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md#getproductpricesbyproductcode--查询产品价格)。

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 三要素 | 请求须密文 | 明文 |
| `data[]` 结构 | 与华安相同 | 体验版 + 正式版两条报价 |

---

## getProductPricesByPolicyId — 按保单 ID 查价格（收银台）

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getProductPricesByPolicyId`（转发华安 `/common/channel/api/getProductPricesByPolicyId`） |
| 采集时间 | 2026-06-25 |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getproductpricesbypolicyid--按保单-id-查价格收银台](./HUAAN_API_SAMPLES.md#getproductpricesbypolicyid--按保单-id-查价格收银台) |

### 请求体（渠道 → 本服务）

仅需 `policyId` 与公共签名字段，无三要素。

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "***",
  "sign": "***",
  "policyId": "POL2026060500001"
}
```

### 响应体（本服务 → 渠道）

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "productCode": "PROD2025001",
    "productName": "基础版医疗险",
    "price": 128.0,
    "historyBinds": [
      {
        "bankName": "中国工商银行",
        "bankCode": "ICBC",
        "isPaySuccess": 1
      }
    ]
  }
}
```

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 华安 URL | 本服务拼 `/common/channel/api/getProductPricesByPolicyId` | 同左 |
| 三要素 | 不需要 | 不需要 |
| 响应 | 原样转发 | 同左 |

---

## getPolicyInfoByPolicyId — 根据保单 ID 查询保单信息

| 项 | 值 |
|----|-----|
| 路径 | `POST {平台域名}/upChannelApi/getPolicyInfoByPolicyId`（转发华安 `/common/channel/api/getPolicyInfoByPolicyId`） |
| 采集时间 | 2026-06-25 |
| 对照华安样例 | [HUAAN_API_SAMPLES.md#getpolicyinfobypolicyid--根据保单-id-查询保单信息](./HUAAN_API_SAMPLES.md#getpolicyinfobypolicyid--根据保单-id-查询保单信息) |

### 请求体（渠道 → 本服务）

仅需 `policyId` 与公共签名字段。

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "key": "***",
  "sign": "***",
  "policyId": "POL2026060500001"
}
```

### 响应体（本服务 → 渠道）

结构与华安一致；`phoneNo`/`name`/`idCard` 及 `insuredList` 内三要素为 **AES 密文**（示例为脱敏明文）：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "insuredList": [
      {
        "insuredCardNo": "341xxxxxxx18",
        "hasSocialSecurity": "1",
        "insuredName": "xxx",
        "phoneNo": "18xxxxxx1"
      }
    ],
    "optionalRiskList": [],
    "productId": "ccxxxx8da3a9ca057772a",
    "idCard": "3xxxxxxxxxxxx8",
    "hasSocialSecurity": "1",
    "policyEndDate": "2026-06-27",
    "policyStartDate": "2026-05-29",
    "policyStatus": "6",
    "productName": "百万医疗险-体验版",
    "phoneNo": "18xxxxx1",
    "productCode": "ZFxxx41001",
    "policyId": "f8xxxxx339d7795c1",
    "payPremium": 0.65,
    "riskList": [
      {
        "riskName": "一般医疗保险金、重大疾病医疗保险金",
        "riskCode": null
      },
      {
        "riskName": "罕见特定恶性肿瘤—重度疾病保险金",
        "riskCode": "1038102"
      }
    ],
    "createTime": "2026-05-28",
    "name": "李xxx",
    "autoRenew": null
  }
}
```

### 与华安原文对照

| 差异点 | 渠道侧 | 华安原文 |
|--------|--------|----------|
| 华安 URL | 本服务拼 `/common/channel/api/getPolicyInfoByPolicyId` | 同左 |
| 三要素 | 响应须加密 | 明文 |
| 请求 | 仅需 policyId + 验签 | 仅需 policyId + channelCode 等 |

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
