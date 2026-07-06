# 华安接口响应样例

本文档保存华安上游接口的**原始 JSON 响应**，便于对照字段格式与联调。数据来自直连华安环境的集成测试，**非**本服务加工后的渠道响应（渠道侧三要素字段会加密，此处为华安明文）。文档中手机号、姓名、身份证号等**用户三要素**均以 `***` 脱敏展示。

渠道商调用本服务时的请求/响应样例见 **[CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md)**，可与本文档成对对照。

## getBankList — 获取支持的银行列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getBankList` |
| 采集时间 | 2026-06-02（路径 2026-06-25 对齐规范） |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "bankName": "中国银行",
      "bankCode": "BOC",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 1,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "建设银行",
      "bankCode": "CCB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 2,
      "payChannelId": "ef70a5a2afea440d86c9a804d0458736"
    },
    {
      "bankName": "农业银行",
      "bankCode": "ABC",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 4,
      "payChannelId": "ef70a5a2afea440d86c9a804d0458736"
    },
    {
      "bankName": "邮储银行",
      "bankCode": "PSBC",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 5,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "华夏银行",
      "bankCode": "HX",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 6,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "平安银行",
      "bankCode": "SZPA",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 7,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "中信银行",
      "bankCode": "ECITIC",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 8,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "广发银行",
      "bankCode": "GDB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 9,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "北京银行",
      "bankCode": "BCCB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 10,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "青岛银行",
      "bankCode": "QDYH",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 12,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "浦发银行",
      "bankCode": "SPDB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 13,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "上海银行",
      "bankCode": "SHB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 14,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "光大银行",
      "bankCode": "CEB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 15,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "南京银行",
      "bankCode": "NJCB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 16,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "恒丰银行",
      "bankCode": "HFB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 17,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "兴业银行",
      "bankCode": "CIB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 18,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    },
    {
      "bankName": "民生银行",
      "bankCode": "CMBCHINA",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 19,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4"
    }
  ]
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息，如「操作成功」 |
| `data` | 根 | array | 银行列表 |
| `bankName` | data[] | string | 银行名称 |
| `bankCode` | data[] | string | 银行编码 |
| `debitCard` | data[] | string | `"1"` 支持储蓄卡，`"0"` 不支持 |
| `creditCard` | data[] | string | `"1"` 支持信用卡，`"0"` 不支持 |
| `sortOrder` | data[] | integer | 排序序号（华安返回，本服务落库时不保存） |
| `payChannelId` | data[] | string | 支付渠道 ID（与 bankCode 配套） |
| `isActBank` | data[] | integer | `1` 活跃银行，`0` 非活跃（规范字段；历史实测可能未返回，默认视为活跃） |
| `sortOrder` | data[] | integer | 排序序号（历史实测返回，本服务落库时不保存） |

本服务将 `getBankList` 成功响应同步至 `bank_info_t` 时，使用 `bankCode`、`bankName`、`debitCard`、`creditCard`、`payChannelId`、`isActBank`；未返回 `isActBank` 时默认 `1`（活跃）。解析逻辑见 `internal/service/bank_sync.go`。

### 复现采集

```bash
set -a && source .env.huaan && set +a
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/getBankList'
```

## getProductInfoByChannel — 渠道查询产品

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getProductInfoByChannel` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串；实测亦可返回 200） |

### 请求体示例

除 `timestamp`、`channelCode` 外无额外业务字段；`key`/`sign` 由本服务 `huaan.Client.Call` 按配置注入。

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "productCode": "ZFHLW1041001",
      "productName": "百万医疗险-体验版"
    },
    {
      "productCode": "ZFHLW1040003",
      "productName": "抗癌险-体验版"
    }
  ]
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息，如「操作成功」 |
| `data` | 根 | array | 该渠道可售产品列表 |
| `productCode` | data[] | string | 产品编码（如 `proInsurance` 等业务接口入参） |
| `productName` | data[] | string | 产品名称 |

### 复现采集

```bash
set -a && source .env.huaan && set +a
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/getProductInfoByChannel'
```

## product/info — 获取渠道产品信息

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/product/info` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 环境 | `https://ins.api.hahealth.ink/`（计划） |
| 渠道编码 | 待确认 |
| 签名 | 待确认 |

按手机号查询渠道产品详情；**上游路径不在 `/upChannelApi` 下**，本服务 `huaan.Client` 经 `upstreamPathOverrides` 转发。

### 请求体示例

除公共字段外须传 `phoneNo`（集成测试为明文三要素）：

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "phoneNo": "***",
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

> 待测试。文档约定示例：

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "productCode": "PRO2026001",
    "productName": "月度体验版意外险",
    "companyName": "XX人寿保险",
    "productType": 1,
    "channelCode": "hy333"
  }
}
```

`productType`：`1` 基础版，`2` 升级版。

### 字段说明

| 字段 | 层级 | 类型 | 说明 |
|------|------|------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `msg` / `message` | 根 | string | 提示信息 |
| `data` | 根 | object | 产品详情 |
| `productCode` | data | string | 产品编码 |
| `productName` | data | string | 产品名称 |
| `companyName` | data | string | 保险公司名称 |
| `productType` | data | integer | `1` 基础版，`2` 升级版 |
| `channelCode` | data | string | 渠道编码 |

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="***" \
HUAAN_TEST_NAME="张三" \
HUAAN_TEST_ID_CARD="110101199001011234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/productInfo'
```

## sms/send — 获取用户短信验证码

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/sms/send` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 环境 | `https://ins.api.hahealth.ink/`（计划） |
| 渠道编码 | 待确认 |
| 签名 | 待确认 |

向指定手机号发送短信验证码；**上游路径不在 `/upChannelApi` 下**。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "phoneNo": "***",
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

> 待测试。文档约定示例：

```json
{
  "code": 200,
  "data": "验证码发送成功",
  "message": "操作成功"
}
```

### 字段说明

| 字段 | 层级 | 类型 | 说明 |
|------|------|------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | string | 发送结果描述 |

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="***" \
HUAAN_TEST_NAME="张三" \
HUAAN_TEST_ID_CARD="110101199001011234" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/smsSend'
```

## sms/valid — 短信验证码校验

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/sms/valid` |
| 采集时间 | 2026-06-25 |
| 联调状态 | **测试环境已通过**（`http://47.97.156.18:9040`，渠道 `BLtJjF`）；hahealth 待测 |
| 环境 | `http://47.97.156.18:9040` |
| 渠道编码 | `BLtJjF` |
| 签名 | 请求头 `sign`；参与字段 `channelCode`、`phoneNo`、`code`、`timestamp`，末尾标准 `&key=华安密钥`，MD5 大写 |

校验短信验证码；请求字段为 **`phoneNo`**（手机号）与 **`code`**（验证码，非 `smsCode`）。

### 请求体示例

```json
{
  "timestamp": "1782703244953",
  "channelCode": "BLtJjF",
  "phoneNo": "***",
  "code": "3875"
}
```

### 上游签名示例

```text
明文: channelCode=BLtJjF&code=3875&phoneNo=***&timestamp=1782703244953&key=fb9ec7236b6c45b7bfd562672e0373ea
sign: F13AC45360ED163467F3D48FD11A33BE
```

### 响应体（华安原始 JSON）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "idCard": "***",
    "name": "***",
    "userId": "5c819c252ab343e2a2a6df98b3a014ab",
    "phoneNo": "***",
    "channelCode": "I7fZcM"
  }
}
```

### 字段说明

| 字段 | 层级 | 类型 | 说明 |
|------|------|------|------|
| `phoneNo` | 请求 | string | 手机号 |
| `code` | 请求 | string | 短信验证码 |
| `code` | 根 | integer | 响应业务码，`200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | object | 用户信息 |
| `userId` | data | string | 用户 ID |
| `idCard` | data | string | 身份证号 |
| `name` | data | string | 姓名 |
| `phoneNo` | data | string | 手机号 |
| `channelCode` | data | string | 渠道编码 |

### 复现采集

```bash
HUAAN_BASE_URL="http://47.97.156.18:9040" \
HUAAN_CHANNEL_CODE="BLtJjF" \
HUAAN_KEY="fb9ec7236b6c45b7bfd562672e0373ea" \
HUAAN_SIGN_ENABLED=true \
HUAAN_TEST_PHONE="***" \
HUAAN_TEST_SMS_CODE="3875" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/smsValid'
```

## sms/noValid — 免短信验证码注册登录

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/sms/noValid` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 环境 | `https://ins.api.hahealth.ink/`（计划） |
| 渠道编码 | 待确认 |
| 签名 | 待确认 |

仅凭手机号注册/登录，无需 `smsCode`。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "mobile": "***",
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

> 待测试。文档约定示例：

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

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="***" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/smsNoValid'
```

## priceByUser — 查询产品价格

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/priceByUser` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 环境 | `https://ins.api.hahealth.ink/`（计划） |
| 渠道编码 | 待确认 |
| 签名 | 待确认 |

按产品编码与用户身份证查询价格；`hasSocialSecurity` 为**整数**（`1`/`0`）；`productPriceList` 可省略。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "productCode": "PROD2025001",
  "idCard": "***",
  "hasSocialSecurity": 1,
  "productPriceList": [
    {
      "liabilityName": "罕见保险金",
      "price": 0.0,
      "kindCode": "102",
      "uwCount": 20
    }
  ],
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

> 待测试。文档约定示例：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "productId": "PID2025001",
    "productName": "综合意外险",
    "productCode": "PROD2025001",
    "productType": "意外险",
    "price": 95.0,
    "yearPremium": 2000,
    "expEffectAfterDays": 1,
    "formalEffectAfterDays": 1
  }
}
```

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PRODUCT_CODE="PROD2025001" \
HUAAN_TEST_ID_CARD="***" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/priceByUser'
```

## policy/phone — 查询用户投保情况

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/policy/phone` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 环境 | `https://ins.api.hahealth.ink/`（计划） |
| 说明 | 仅返回**基础版**产品，用于二次进入页面判断是否已购并跳转 |

### 请求体示例

`phoneNo` 与 `userId` **至少传一项**：

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "phoneNo": "***",
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

> 待测试。文档约定示例：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "insuredList": [],
      "productCode": "Z041001",
      "productId": "cc55f8cff47a4f0d88da3a9ca057772a",
      "payPremium": "",
      "riskList": [],
      "isPolicy": "0",
      "policyNo": "",
      "policyStatus": "",
      "productName": "百万医疗险-体验版"
    }
  ]
}
```

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="***" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/policyByPhone'
```

## getUserInfoByPhoneNo — 查询用户信息

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getUserInfoByPhoneNo` |
| 采集时间 | 2026-06-25（新增） |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 环境 | `https://ins.api.hahealth.ink/`（计划） |

### 请求体示例

`phoneNo` 与 `userId` **至少传一项**：

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "phoneNo": "***",
  "key": "",
  "sign": ""
}
```

或：

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "userId": "xxx",
  "key": "",
  "sign": ""
}
```

### 响应体（规范）

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

### 字段说明

| 字段 | 层级 | 类型 | 说明 |
|------|------|------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 如「请求成功」 |
| `data` | 根 | object | 用户信息 |
| `userId` | data | string | 用户 ID |
| `idCard` | data | string | 证件号（直连明文；经本服务返回渠道时加密） |
| `name` | data | string | 姓名 |
| `phoneNo` | data | string | 手机号 |
| `channelCode` | data | string | 渠道编码 |

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PHONE="***" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/getUserInfoByPhoneNo'
```

## getLiabilitiesByProductId — 查询可选责任列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getLiabilitiesByProductId` |
| 采集时间 | — |
| 联调状态 | **待测试**（见 [TESTING.md §7](./TESTING.md#7-待测试接口列表)） |
| 说明 | 配合 [priceByUser](#pricebyuser--查询产品价格) 使用，返回可选责任及默认份数 |

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "YOUR_CHANNEL_CODE",
  "productCode": "PROD2025001",
  "idCard": "***",
  "hasSocialSecurity": 1,
  "productType": 1,
  "key": "",
  "sign": ""
}
```

### 响应体（华安原始 JSON）

> 待测试。文档约定示例：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "liabilityName": "高保险金",
      "price": 3.16,
      "kindCode": "101",
      "uwCount": 20
    },
    {
      "liabilityName": "罕见保险金",
      "price": 0.05,
      "kindCode": "102",
      "uwCount": 20
    }
  ]
}
```

### 复现采集

```bash
HUAAN_BASE_URL="https://ins.api.hahealth.ink/" \
HUAAN_CHANNEL_CODE="你的渠道编码" \
HUAAN_TEST_PRODUCT_CODE="PROD2025001" \
HUAAN_TEST_ID_CARD="***" \
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/liabilitiesByProductId'
```

## proInsurance — 预投保

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/proxy/upChannelApi/proInsurance` |
| 采集时间 | 2026-06-02（路径/字段 2026-06-25 对齐规范） |
| 环境 | `https://zf.ins.api.york.xin/`（历史实测）；新环境见 [TESTING.md §7](./TESTING.md#7-待测试接口列表) |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

`productCode` 取自上文 [getProductInfoByChannel](#getproductinfobychannel--渠道查询产品) 响应（集成测试自动解析）。**须传** `hasSocialSecurity`、`isUpgrade`、`autoRenew`（整数）；未传 `hasSocialSecurity` 时华安返回 500。成功时 `policyId` 每次预投保新生成；规范响应含 `policyNo`。

### 请求体示例（成功用例，华安直连明文三要素）

```json
{
  "timestamp": "1780456222458",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "productCode": "ZFHLW1041001",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***",
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

| 请求字段 | 实测值（示例） | 说明 |
|----------|----------------|------|
| `productCode` | `ZFHLW1041001` / `ZFHLW1040003` | 来自 getProductInfoByChannel |
| `phoneNo` | `***` | 明文（直连华安；经本服务转发时须 AES 密文） |
| `name` | `***` | 同上 |
| `idCard` | `***` | 同上 |
| `hasSocialSecurity` | `1` | 是否有社保（**整数**） |
| `isUpgrade` | `0` | 是否升级：1 升级，0 不升级 |
| `autoRenew` | `1` | 是否次年自动续保 |
| `productPriceList` | 可选 | 不传则使用后台默认份数 |

### 未传 hasSocialSecurity 时的响应（失败）

#### productCode = `ZFHLW1041001`

```json
{
  "code": 500,
  "message": "预投保失败。未找到价格规则，productId=cc55f8cff47a4f0d88da3a9ca057772a, age=32, hasSocialSecurity=null",
  "data": null
}
```

#### productCode = `ZFHLW1040003`

```json
{
  "code": 500,
  "message": "预投保失败。未找到价格规则，productId=3ecd3842febd42e2854403d9cfe67c36, age=32, hasSocialSecurity=null",
  "data": null
}
```

### 响应体（华安原始 JSON，hasSocialSecurity=1）

#### productCode = `ZFHLW1041001`（百万医疗险-体验版）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "policyId": "9697d878b6474c619fa43ffa9ca3b4b0",
    "policyStatus": "0",
    "userId": "5c819c252ab343e2a2a6df98b3a014ab"
  }
}
```

#### productCode = `ZFHLW1040003`（抗癌险-体验版）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "policyId": "ab185a401820431aad4be32c329ab003",
    "policyStatus": "0",
    "userId": "5c819c252ab343e2a2a6df98b3a014ab"
  }
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息或失败原因 |
| `data` | 根 | object / null | 成功时为投保结果 |
| `policyId` | data | string | 保单 ID，可供 getProductPricesByPolicyId 等接口使用 |
| `policyNo` | data | string | 保单号（规范字段；历史实测可能未返回） |
| `policyStatus` | data | string | 保单状态（历史实测为 `"0"`，可能额外返回） |
| `userId` | data | string | 华安用户 ID（可能额外返回，与 verifyNoCode 一致） |

### 复现采集

```bash
set -a && source .env.huaan && set +a
export HUAAN_TEST_PHONE="***"
export HUAAN_TEST_NAME="***"
export HUAAN_TEST_ID_CARD="***"

go test -tags=integration ./internal/huaan/ -v -count=1 \
  -run 'TestHuaAnDirect_ProductFromChannel/ZFHLW1041001/proInsurance'
# 或跑渠道下全部产品的报价+投保：
# go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_ProductFromChannel
```

## verifyNoCode — 无验证码实名

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/verifyNoCode` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***"
}
```

| 请求字段 | 实测值 | 说明 |
|----------|--------|------|
| `phoneNo` | `***` | 明文（直连华安，非渠道加密） |
| `name` | `***` | 同上 |
| `idCard` | `***` | 同上 |

### 响应体（华安原始 JSON）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "idCard": "***",
    "name": "***",
    "userId": "5c819c252ab343e2a2a6df98b3a014ab",
    "phoneNo": "***",
    "channelCode": "I7fZcM"
  }
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息，如「操作成功」 |
| `data` | 根 | object | 实名结果 |
| `userId` | data | string | 华安侧用户 ID |
| `phoneNo` | data | string | 手机号（明文） |
| `name` | data | string | 姓名（明文） |
| `idCard` | data | string | 身份证号（明文） |
| `channelCode` | data | string | 渠道编码 |

经本服务转发给渠道时，`phoneNo`、`name`、`idCard` 等 PII 字段会加密返回。

### 复现采集

```bash
set -a && source .env.huaan && set +a
export HUAAN_TEST_PHONE="***"
export HUAAN_TEST_NAME="***"
export HUAAN_TEST_ID_CARD="***"
# 联调时将上述 *** 替换为真实三要素
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/verifyNoCode'
```

## getPolicyInfoByPhoneNo — 按手机号查保单

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getPolicyInfoByPhoneNo` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

`userId` 取自上文 [verifyNoCode](#verifynocode--无验证码实名) 响应中的 `data.userId`。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "userId": "5c819c252ab343e2a2a6df98b3a014ab",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***"
}
```

| 请求字段 | 实测值 | 说明 |
|----------|--------|------|
| `userId` | `5c819c252ab343e2a2a6df98b3a014ab` | 来自 verifyNoCode 响应 |
| `phoneNo` | `***` | 与 verifyNoCode 一致 |
| `name` | `***` | 同上 |
| `idCard` | `***` | 同上 |

### 响应体（华安原始 JSON）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "insuredList": [],
      "productCode": "ZFHLW1041001",
      "productId": "cc55f8cff47a4f0d88da3a9ca057772a",
      "payPremium": "",
      "riskList": [],
      "isPolicy": "0",
      "policyNo": "",
      "policyStatus": "",
      "productName": "百万医疗险-体验版"
    },
    {
      "insuredList": [],
      "productCode": "ZFHLW1040003",
      "productId": "3ecd3842febd42e2854403d9cfe67c36",
      "payPremium": "",
      "riskList": [],
      "isPolicy": "0",
      "policyNo": "",
      "policyStatus": "",
      "productName": "抗癌险-体验版"
    }
  ]
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | array | 保单/产品列表 |
| `productCode` | data[] | string | 产品编码 |
| `productName` | data[] | string | 产品名称 |
| `productId` | data[] | string | 产品 ID |
| `isPolicy` | data[] | string | `"0"` 表示暂无有效保单（实测） |
| `policyNo` | data[] | string | 保单号，未投保时为空 |
| `policyStatus` | data[] | string | 保单状态，未投保时为空 |
| `payPremium` | data[] | string | 保费 |
| `insuredList` | data[] | array | 被保人列表 |
| `riskList` | data[] | array | 险种列表 |

### 复现采集

先执行 verifyNoCode 取得 `userId`，再调用本接口（集成测试默认不含 `userId`，需自行拼接或脚本串联）：

```bash
set -a && source .env.huaan && set +a
export HUAAN_TEST_PHONE="***"
export HUAAN_TEST_NAME="***"
export HUAAN_TEST_ID_CARD="***"
# 联调时将上述 *** 替换为真实三要素

# 1) 获取 userId
go test -tags=integration ./internal/huaan/ -v -count=1 -run 'TestHuaAnDirect_AllPaths/verifyNoCode'

# 2) 将 verifyNoCode 响应中的 userId 写入请求后调用 getPolicyInfoByPhoneNo
#    （或参考 verifyNoCode 与 getPolicyInfoByPhoneNo 串联脚本）
```

## getProductPricesByProductCode — 查询产品价格

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getProductPricesByProductCode` |
| 采集时间 | 2026-06-02（`TestHuaAnDirect_ProductFromChannel` 直连实测） |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

`productCode` 取自上文 [getProductInfoByChannel](#getproductinfobychannel--渠道查询产品) 响应。实测除 `productCode` 外还需三要素及 `hasSocialSecurity`；仅传 `productCode` 时返回 `601`。成功时 `data` 为**数组**，通常含体验版（`productType: "1"`）与正式版（`productType: "2"`）两条报价。

### 请求体示例（成功用例，华安直连明文三要素）

```json
{
  "timestamp": "1780456222050",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "productCode": "ZFHLW1041001",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***",
  "hasSocialSecurity": "1"
}
```

| 请求字段 | 实测值（示例） | 说明 |
|----------|----------------|------|
| `productCode` | `ZFHLW1041001` / `ZFHLW1040003` | 来自 getProductInfoByChannel |
| `phoneNo` | `***` | 明文（直连华安；经本服务转发时须 AES 密文） |
| `name` | `***` | 同上 |
| `idCard` | `***` | 同上 |
| `hasSocialSecurity` | `"1"` | 是否有社保；未传时华安可能报「未找到价格规则… hasSocialSecurity=null」 |

### 仅 productCode 时的响应（失败）

```json
{
  "code": 601,
  "message": null,
  "data": null
}
```

### 响应体（华安原始 JSON）

#### productCode = `ZFHLW1041001`（百万医疗险-体验版）

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

#### productCode = `ZFHLW1040003`（抗癌险-体验版）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "productCode": "ZFHLW1040003",
      "productId": "3ecd3842febd42e2854403d9cfe67c36",
      "price": 0.70,
      "productName": "抗癌险-体验版",
      "productType": "1"
    },
    {
      "productCode": "ZFHLW1040002",
      "productId": "c19338f72a354ab69307381b9647ec47",
      "price": 42.22,
      "productName": "抗癌险-正式版",
      "productType": "2"
    }
  ]
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | array | 价格列表（含体验版与正式版等） |
| `productCode` | data[] | string | 产品编码 |
| `productName` | data[] | string | 产品名称 |
| `productId` | data[] | string | 产品 ID |
| `price` | data[] | number | 价格 |
| `productType` | data[] | string | `"1"` 体验版，`"2"` 正式版（实测） |

### 复现采集

```bash
set -a && source .env.huaan && set +a
export HUAAN_TEST_PHONE="***"
export HUAAN_TEST_NAME="***"
export HUAAN_TEST_ID_CARD="***"

go test -tags=integration ./internal/huaan/ -v -count=1 \
  -run 'TestHuaAnDirect_ProductFromChannel/ZFHLW1041001/getProductPricesByProductCode'
# 或跑渠道下全部产品的报价+投保：
# go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_ProductFromChannel
```

## upGradeIns — 升级险种

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/upGradeIns` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭 |

`policyId` 取自上文 [proInsurance](#proinsurance--预投保) 成功响应。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "policyId": "91a51cc70e724d4885ac2e27530099ec"
}
```

### 响应体（华安原始 JSON）

#### policyId = `2c61ea7a04bb4b91a688ff1bf1b04dc4`（百万医疗险-体验版）

```json
{
  "code": 601,
  "message": null,
  "data": null
}
```

#### policyId = `91a51cc70e724d4885ac2e27530099ec`（抗癌险-体验版）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": "升级成功"
}
```

---

## getSignUrl — 获取签约链接

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getSignUrl` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭 |

`policyId` 来自 proInsurance；`userId` 来自 [verifyNoCode](#verifynocode--无验证码实名)。

**`bankCode` 与 `payChannelId` 必须配套**：二者取自 [getBankList](#getbanklist--获取支持的银行列表) **同一条** `data[]` 记录，不得将 A 银行的 `bankCode` 与 B 银行的 `payChannelId` 组合。`cardType` 与所选银行的 `debitCard`/`creditCard` 能力一致（储蓄卡 `"1"`，信用卡 `"2"`）。

集成测试 `TestHuaAnDirect_GetSignUrl` 会先调 `getBankList` 再调 `getSignUrl`，并在 `-v` 日志中输出两次请求与响应。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "policyId": "91a51cc70e724d4885ac2e27530099ec",
  "bankCode": "BOC",
  "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
  "cardType": "1",
  "userId": "5c819c252ab343e2a2a6df98b3a014ab",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***"
}
```

### 仅测 getSignUrl（华安直连）

```bash
set -a && source .env.huaan && set +a
go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_GetSignUrl
```

环境变量见 [TESTING.md](./TESTING.md) 第 3.6 节。

### 响应体（华安原始 JSON）

#### policyId = `2c61ea7a04bb4b91a688ff1b04dc4`

```json
{
  "code": 500,
  "message": "保单不存在",
  "data": null
}
```

#### policyId = `91a51cc70e724d4885ac2e27530099ec`

```json
{
  "code": 500,
  "message": "支付通道不存在",
  "data": null
}
```

---

## getProductPricesByPolicyId — 按保单 ID 查价格（收银台）

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getProductPricesByPolicyId` |
| 采集时间 | 2026-06-02（路径/字段 2026-06-25 对齐规范） |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭 |

`policyId` 来自 proInsurance；**仅需 `policyId`**，无需三要素。

### 请求体示例（成功用例）

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "policyId": "91a51cc70e724d4885ac2e27530099ec"
}
```

### 响应体（规范）

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

### 历史实测（旧环境，字段较少）

#### policyId = `2c61ea7a04bb4b91a688ff1b04dc4`

```json
{
  "code": 601,
  "message": null,
  "data": null
}
```

#### policyId = `91a51cc70e724d4885ac2e27530099ec`

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "price": 0.70
  }
}
```

### 字段说明

| 字段 | 层级 | 类型 | 说明 |
|------|------|------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | object | 收银台展示数据 |
| `productCode` | data | string | 产品编码 |
| `productName` | data | string | 产品名称 |
| `price` | data | number | 支付金额 |
| `historyBinds` | data | array | 历史绑卡/支付记录（逐步启用） |
| `bankName` | historyBinds[] | string | 银行名称 |
| `bankCode` | historyBinds[] | string | 银行编码 |
| `isPaySuccess` | historyBinds[] | integer | `1` 成功，`0` 失败；缺省表示回调未到，勿重复绑卡 |

---

## getPolicyInfoByPolicyId — 根据保单 ID 查询保单信息

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/common/channel/api/getPolicyInfoByPolicyId` |
| 采集时间 | 2026-06-02（路径 2026-06-25 对齐规范） |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭 |

`policyId` 来自 proInsurance。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "policyId": "91a51cc70e724d4885ac2e27530099ec"
}
```

### 响应体（华安原始 JSON）

#### policyId = `2c61ea7a04bb4b91a688ff1b04dc4`

```json
{
  "code": 500,
  "message": "保单不存在",
  "data": null
}
```

#### policyId = `91a51cc70e724d4885ac2e27530099ec`（抗癌险-体验版）

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "insuredList": [
      {
        "insuredCardNo": "***",
        "hasSocialSecurity": "1",
        "insuredName": "***",
        "phoneNo": "***"
      }
    ],
    "optionalRiskList": [],
    "productId": "3ecd3842febd42e2854403d9cfe67c36",
    "idCard": "***",
    "hasSocialSecurity": "1",
    "policyEndDate": null,
    "policyStartDate": null,
    "policyStatus": "0",
    "productName": "抗癌险-体验版",
    "phoneNo": "***",
    "productCode": "ZFHLW1040003",
    "policyId": "91a51cc70e724d4885ac2e27530099ec",
    "payPremium": 0.70,
    "riskList": [
      {
        "riskName": "恶性肿瘤医疗保险金",
        "riskCode": null
      },
      {
        "riskName": "罕见特定恶性肿瘤——重度疾病保险金",
        "riskCode": "1040102"
      }
    ],
    "createTime": "2026-06-02",
    "name": "***",
    "autoRenew": null
  }
}
```

### 字段说明

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | object | 保单详情 |
| `policyId` | data | string | 保单 ID |
| `policyStatus` | data | string | 保单状态（如 `"0"`、`"6"`，视保单阶段而定） |
| `productCode` / `productName` / `productId` | data | string | 产品编码 / 名称 / ID |
| `payPremium` | data | number | 保费 |
| `policyStartDate` / `policyEndDate` | data | string / null | 保障起止日期 |
| `createTime` | data | string | 创建日期 |
| `autoRenew` | data | integer / null | 是否自动续保 |
| `riskList` | data | array | 已选险种（`riskCode` 可为 null） |
| `optionalRiskList` | data | array | 可选险种 |
| `insuredList` | data | array | 被保人列表（含三要素，明文） |

### 保单 ID 联调说明

| policyId | 来源 productCode | upGradeIns | getSignUrl | getProductPricesByPolicyId | getPolicyInfoByPolicyId |
|----------|------------------|------------|------------|----------------------------|-------------------------|
| `2c61ea7a04bb4b91a688ff1bf1b04dc4` | ZFHLW1041001 | 601 | 500 保单不存在 | 601 | 500 保单不存在 |
| `91a51cc70e724d4885ac2e27530099ec` | ZFHLW1040003 | 200 | 500 支付通道不存在 | 200 | 200 |

百万医疗险 `policyId` 仅 `proInsurance` 返回成功，按保单查询类接口多不可用，可能与 `policyStatus: "0"` 未生效有关。

## 相关文档

- [TESTING.md](./TESTING.md) — 华安直连集成测试
- [ARCHITECTURE.md](./ARCHITECTURE.md) — getBankList 缓存与转发
- [CHANNEL_API.md](./CHANNEL_API.md) — 渠道商接口文档
- [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md) — 渠道 → 本服务请求/响应样例
