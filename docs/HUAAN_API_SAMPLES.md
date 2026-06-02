# 华安接口响应样例

本文档保存华安上游接口的**原始 JSON 响应**，便于对照字段格式与联调。数据来自直连华安环境的集成测试，**非**本服务加工后的渠道响应（渠道侧三要素字段会加密，此处为华安明文）。文档中手机号、姓名、身份证号等**用户三要素**均以 `***` 脱敏展示。

## getBankList — 获取银行列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getBankList` |
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
| `payChannelId` | data[] | string | 支付渠道 ID（华安返回，本服务落库时不保存） |

本服务将 `getBankList` 成功响应同步至 `bank_info_t` 时，使用 `bankCode`、`bankName`、`debitCard`、`creditCard`；未返回 `status` 时默认 `1`（启用）。解析逻辑见 `internal/service/bank_sync.go`。

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

## proInsurance — 投保

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/proInsurance` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

`productCode` 取自上文 [getProductInfoByChannel](#getproductinfobychannel--渠道查询产品) 响应，分别实测一次。**实测须传 `hasSocialSecurity`**（如 `"1"`），否则返回 500。

### 请求体示例（成功用例）

```json
{
  "timestamp": "1717300000000",
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

| 请求字段 | 实测值 | 说明 |
|----------|--------|------|
| `productCode` | 见下表各用例 | 来自 getProductInfoByChannel |
| `phoneNo` | `***` | 明文（直连华安，非渠道加密） |
| `name` | `***` | 同上 |
| `idCard` | `***` | 同上 |
| `hasSocialSecurity` | `"1"` | 是否有社保；未传时华安解析为 `null` 并报价格规则错误 |

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
    "policyId": "2c61ea7a04bb4b91a688ff1bf1b04dc4",
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
    "policyId": "91a51cc70e724d4885ac2e27530099ec",
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
| `policyStatus` | data | string | 保单状态（实测为 `"0"`） |
| `userId` | data | string | 华安用户 ID，与 verifyNoCode 一致 |

### 复现采集

```bash
set -a && source .env.huaan && set +a
export HUAAN_TEST_PHONE="***"
export HUAAN_TEST_NAME="***"
export HUAAN_TEST_ID_CARD="***"
# 联调时将上述 *** 替换为真实三要素

HUAAN_TEST_PRODUCT_CODE=ZFHLW1041001 go test -tags=integration ./internal/huaan/ -v -count=1 \
  -run 'TestHuaAnDirect_AllPaths/proInsurance'
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
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，`key`/`sign` 均为空字符串） |

`productCode` 取自上文 [getProductInfoByChannel](#getproductinfobychannel--渠道查询产品) 响应。实测除 `productCode` 外还需三要素及 `hasSocialSecurity`；仅传 `productCode` 时返回 `601`。

### 请求体示例（成功用例）

```json
{
  "timestamp": "1717300000000",
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

| 请求字段 | 实测值 | 说明 |
|----------|--------|------|
| `productCode` | 见下表各用例 | 来自 getProductInfoByChannel |
| `phoneNo` | `***` | 明文（直连华安） |
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
      "price": 150.80,
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
# 联调时将上述 *** 替换为真实三要素；集成测试已默认 hasSocialSecurity=1

HUAAN_TEST_PRODUCT_CODE=ZFHLW1041001 go test -tags=integration ./internal/huaan/ -v -count=1 \
  -run 'TestHuaAnDirect_AllPaths/getProductPricesByProductCode'
```

## upGradeIns — 升级险种

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/upGradeIns` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭 |

`policyId` 取自上文 [proInsurance](#proinsurance--投保) 成功响应。

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

`policyId` 来自 proInsurance；`bankCode`/`cardType` 取自 [getBankList](#getbanklist--获取银行列表) 第一组（BOC，`cardType=1` 表示储蓄卡）；`userId` 来自 [verifyNoCode](#verifynocode--无验证码实名)。

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "policyId": "91a51cc70e724d4885ac2e27530099ec",
  "bankCode": "BOC",
  "cardType": "1",
  "userId": "5c819c252ab343e2a2a6df98b3a014ab",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***"
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

#### policyId = `91a51cc70e724d4885ac2e27530099ec`

```json
{
  "code": 500,
  "message": "支付通道不存在",
  "data": null
}
```

---

## getProductPricesByPolicyId — 按保单报价

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getProductPricesByPolicyId` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭 |

`policyId` 来自 proInsurance；另需三要素及 `hasSocialSecurity: "1"`（与 getProductPricesByProductCode 一致）。

### 请求体示例（成功用例）

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM",
  "key": "",
  "sign": "",
  "policyId": "91a51cc70e724d4885ac2e27530099ec",
  "phoneNo": "***",
  "name": "***",
  "idCard": "***",
  "hasSocialSecurity": "1"
}
```

### 响应体（华安原始 JSON）

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

| 字段 | 层级 | 类型（实测） | 说明 |
|------|------|--------------|------|
| `code` | 根 | integer | `200` 表示成功 |
| `message` | 根 | string | 提示信息 |
| `data` | 根 | object | 含 `price` 等 |
| `price` | data | number | 保费价格 |

---

## getPolicyInfoByPolicyId — 按保单 ID 查询

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getPolicyInfoByPolicyId` |
| 采集时间 | 2026-06-02 |
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
| `policyStatus` | data | string | 保单状态（实测 `"0"`） |
| `productCode` / `productName` | data | string | 产品编码 / 名称 |
| `payPremium` | data | number | 保费 |
| `riskList` | data | array | 险种列表 |
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
- [CHANNEL_INTEGRATION.md](./CHANNEL_INTEGRATION.md) — 渠道侧调用说明
