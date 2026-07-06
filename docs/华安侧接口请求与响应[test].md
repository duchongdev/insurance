# 华安侧接口请求与响应[test]

本文档记录在本机通过 `huaan.Client` 逻辑**直连华安**（不经渠道验签、三要素**明文**）的实测结果。文档中手机号、姓名、身份证号等**用户三要素**均以 `***` 脱敏展示。

## 测试环境

| 项 | 值 |
|----|-----|
| 华安地址 | `http://47.97.156.18:9040` |
| 渠道编码 | `BLtJjF` |
| 测试时间 | 2026-07-06T15:53:02+08:00 |
| 签名 | **开启**（`HUAAN_SIGN_ENABLED=true`） |
| 签名规则 | 非空参数按 key 字典序 `k=v&` 拼接，末尾 `&key=渠道密钥`，MD5 **大写** 32 位；写入 HTTP 头 `sign`；**密钥不写 body** |
| 参数来源 | `productCode` ← **product/info**；`bankCode`/`payChannelId` ← **getBankList**；`policyId`/`userId` ← **proInsurance** |
| 探测脚本 | `go run ./scripts/huaan_direct_probe/main.go`（`HUAAN_PROBE_SKIP_SMS=false`） |

三要素测试账号（文档脱敏为 `***`）：手机号、姓名、身份证号来自 `HUAAN_TEST_*` 环境变量。

## 接口明细

### 1. getBankList — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getBankList` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/getBankList` |
| 业务码 | `200` |
| 耗时 | 133 ms |
| 说明 | 无业务参数；bankCode/payChannelId 取自本接口 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "timestamp": "1783324382058"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "",
  "data": [
    {
      "bankName": "建设银行",
      "bankCode": "CCB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 1,
      "payChannelId": "ef70a5a2afea440d86c9a804d0458736",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "农业银行",
      "bankCode": "ABC",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 4,
      "payChannelId": "ef70a5a2afea440d86c9a804d0458736",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "邮储银行",
      "bankCode": "PSBC",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 5,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "中国银行",
      "bankCode": "BOC",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 6,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "华夏银行",
      "bankCode": "HX",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 6,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "平安银行",
      "bankCode": "SZPA",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 7,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "中信银行",
      "bankCode": "ECITIC",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 8,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "广发银行",
      "bankCode": "GDB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 9,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "北京银行",
      "bankCode": "BCCB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 10,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "青岛银行",
      "bankCode": "QDYH",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 12,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "浦发银行",
      "bankCode": "SPDB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 13,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "上海银行",
      "bankCode": "SHB",
      "debitCard": "1",
      "creditCard": "1",
      "sortOrder": 14,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "光大银行",
      "bankCode": "CEB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 15,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "南京银行",
      "bankCode": "NJCB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 16,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "恒丰银行",
      "bankCode": "HFB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 17,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "兴业银行",
      "bankCode": "CIB",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 18,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    },
    {
      "bankName": "民生银行",
      "bankCode": "CMBC",
      "debitCard": "1",
      "creditCard": "0",
      "sortOrder": 19,
      "payChannelId": "4bf5dd6ece6847e68ee6c5d4345afee4",
      "isActBank": "0",
      "activityScore": null
    }
  ]
}
```

### 2. product/info — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/product/info` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/product/info` |
| 业务码 | `200` |
| 耗时 | 93 ms |
| 说明 | 三要素明文；productCode 取自本接口 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "idCard": "***",
  "name": "***",
  "phoneNo": "***",
  "timestamp": "1783324382192"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "productCode": "HuaAnBao",
      "productName": "百万医疗-体验版"
    },
    {
      "productCode": "HuaAnBaoB",
      "productName": "重疾险-体验版"
    }
  ]
}
```

### 3. getUserInfoByPhoneNo — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getUserInfoByPhoneNo` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/getUserInfoByPhoneNo` |
| 业务码 | `200` |
| 耗时 | 90 ms |
| 说明 | phoneNo 明文 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "phoneNo": "***",
  "timestamp": "1783324382287"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "idCard": "***",
    "name": "***",
    "userId": "5c819c252ab343e2a2a6df98b3a014ab",
    "phoneNo": "***"
  }
}
```

### 4. policy/phone — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/policy/phone` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/policy/phone` |
| 业务码 | `200` |
| 耗时 | 131 ms |
| 说明 | phoneNo 明文 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "phoneNo": "***",
  "timestamp": "1783324382380"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "insuredList": [],
      "productCode": "HuaAnBao",
      "productId": "859db15daa024f43bde30eab9b5f5da2",
      "payPremium": "",
      "riskList": [],
      "isPolicy": "0",
      "policyNo": "",
      "policyStatus": "",
      "productName": "百万医疗-体验版"
    },
    {
      "insuredList": [],
      "productCode": "HuaAnBaoB",
      "productId": "b465c1f3cc4c42ac8033b146fed44f0a",
      "payPremium": "",
      "riskList": [],
      "isPolicy": "0",
      "policyNo": "",
      "policyStatus": "",
      "productName": "重疾险-体验版"
    }
  ]
}
```

### 5. getProductPricesByProductCode — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getProductPricesByProductCode` |
| 华安路径 | `POST http://47.97.156.18:9040/upChannelApi/getProductPricesByProductCode` |
| 依赖 | product/info → productCode=HuaAnBao |
| 业务码 | `200` |
| 耗时 | 153 ms |
| 说明 | hasSocialSecurity 为字符串；含三要素 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "hasSocialSecurity": "1",
  "idCard": "***",
  "name": "***",
  "phoneNo": "***",
  "productCode": "HuaAnBao",
  "timestamp": "1783324382515"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "productCode": "HuaAnBao",
      "productId": "859db15daa024f43bde30eab9b5f5da2",
      "price": 0.6,
      "productName": "百万医疗-体验版",
      "productType": "1"
    },
    {
      "productCode": "HuaAnBao",
      "productId": "d3177469d882467482c58d69aed4972c",
      "price": 155.4,
      "productName": "百万医疗-正式版",
      "productType": "2"
    }
  ]
}
```

### 6. priceByUser — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/priceByUser` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/priceByUser` |
| 依赖 | product/info → productCode=HuaAnBao |
| 业务码 | `200` |
| 耗时 | 118 ms |
| 说明 | hasSocialSecurity 为整数 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "hasSocialSecurity": 1,
  "idCard": "***",
  "productCode": "HuaAnBao",
  "timestamp": "1783324382670"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "formalEffectAfterDays": 30,
      "productCode": "HuaAnBao",
      "productId": "859db15daa024f43bde30eab9b5f5da2",
      "price": 0.6,
      "yearPremium": 7.2,
      "expEffectAfterDays": 1,
      "productName": "百万医疗-体验版",
      "productType": "1"
    },
    {
      "formalEffectAfterDays": 30,
      "productCode": "HuaAnBao",
      "productId": "d3177469d882467482c58d69aed4972c",
      "price": 155.4,
      "yearPremium": 1864.8,
      "expEffectAfterDays": 1,
      "productName": "百万医疗-正式版",
      "productType": "2"
    }
  ]
}
```

### 7. getLiabilitiesByProductId — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getLiabilitiesByProductId` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/getLiabilitiesByProductId` |
| 依赖 | product/info → productCode=HuaAnBao |
| 业务码 | `200` |
| 耗时 | 104 ms |
| 说明 | productType 1=体验版 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "hasSocialSecurity": 1,
  "idCard": "***",
  "productCode": "HuaAnBao",
  "productType": 1,
  "timestamp": "1783324382794"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": []
}
```

### 8. proInsurance — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/proInsurance` |
| 华安路径 | `POST http://47.97.156.18:9040/proxy/upChannelApi/proInsurance` |
| 依赖 | product/info → productCode=HuaAnBao |
| 业务码 | `200` |
| 耗时 | 536 ms |
| 说明 | 成功时返回 policyId/userId |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "autoRenew": 1,
  "channelCode": "BLtJjF",
  "hasSocialSecurity": 1,
  "idCard": "***",
  "isUpgrade": 0,
  "name": "***",
  "phoneNo": "***",
  "productCode": "HuaAnBao",
  "timestamp": "1783324382899"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "policyId": "790ebf3612a343e18e2c6b0793560b3e",
    "policyStatus": "0",
    "userId": "5c819c252ab343e2a2a6df98b3a014ab"
  }
}
```

### 9. getPolicyInfoByPolicyId — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getPolicyInfoByPolicyId` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/getPolicyInfoByPolicyId` |
| 依赖 | proInsurance → policyId=790ebf3612a343e18e2c6b0793560b3e |
| 业务码 | `200` |
| 耗时 | 123 ms |
| 说明 | 仅需 policyId |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "policyId": "790ebf3612a343e18e2c6b0793560b3e",
  "timestamp": "1783324383444"
}
```

**响应体**：

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
    "productId": "859db15daa024f43bde30eab9b5f5da2",
    "idCard": "***",
    "hasSocialSecurity": "1",
    "policyEndDate": "2026-08-06",
    "policyStartDate": "2026-07-07",
    "policyStatus": "0",
    "productName": "百万医疗-体验版",
    "phoneNo": "***",
    "productCode": "HuaAnBao",
    "policyId": "790ebf3612a343e18e2c6b0793560b3e",
    "payPremium": 0.6,
    "riskList": [
      {
        "riskName": "融盛-百万医疗-体验版-必选责任",
        "riskCode": null
      }
    ],
    "createTime": "2026-07-06",
    "name": "***",
    "autoRenew": "1"
  }
}
```

### 10. getProductPricesByPolicyId — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getProductPricesByPolicyId` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/getProductPricesByPolicyId` |
| 依赖 | proInsurance → policyId=790ebf3612a343e18e2c6b0793560b3e |
| 业务码 | `200` |
| 耗时 | 115 ms |
| 说明 | 仅需 policyId |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "policyId": "790ebf3612a343e18e2c6b0793560b3e",
  "timestamp": "1783324383570"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "productCode": "HuaAnBao",
    "price": 0.6,
    "historyBinds": [],
    "productName": "百万医疗-体验版"
  }
}
```

### 11. getSignUrl — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getSignUrl` |
| 华安路径 | `POST http://47.97.156.18:9040/upChannelApi/getSignUrl` |
| 依赖 | proInsurance → policyId=790ebf3612a343e18e2c6b0793560b3e; getBankList → bankCode=CCB |
| 业务码 | `200` |
| 耗时 | 907 ms |
| 说明 | product=百万医疗-体验版 bank=CCB |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "bankCode": "CCB",
  "cardType": "1",
  "channelCode": "BLtJjF",
  "idCard": "***",
  "name": "***",
  "payChannelId": "ef70a5a2afea440d86c9a804d0458736",
  "phoneNo": "***",
  "policyId": "790ebf3612a343e18e2c6b0793560b3e",
  "timestamp": "1783324383695",
  "userId": "5c819c252ab343e2a2a6df98b3a014ab"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "signId": "1638199514013630464",
    "url": "http://yl.yming.xin/bx/sign-test.html?signId=1638199514013630464"
  }
}
```

### 12. sms/send — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/sms/send` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/sms/send` |
| 业务码 | `200` |
| 耗时 | 431 ms |
| 说明 | 三要素明文 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "idCard": "***",
  "name": "***",
  "phoneNo": "***",
  "timestamp": "1783324384602"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": "验证码发送成功"
}
```

### 13. sms/noValid — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/sms/noValid` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/sms/noValid` |
| 业务码 | `200` |
| 耗时 | 114 ms |
| 说明 | mobile 明文 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "mobile": "***",
  "timestamp": "1783324385041"
}
```

**响应体**：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "idCard": "",
    "name": "",
    "userId": "aeebb589ed224ca5a97d4f8f1df252e8",
    "phoneNo": null,
    "channelCode": "BLtJjF"
  }
}
```

### 14. sms/valid — 成功

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/sms/valid` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/sms/valid` |
| 依赖 | sms/send |
| 业务码 | `200` |
| 耗时 | 65 ms |
| 说明 | code 明文 |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "code": "3875",
  "phoneNo": "***",
  "timestamp": "1783324385159"
}
```

**响应体**：

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

### 15. getPhoneByToken — 失败

| 项 | 值 |
|----|-----|
| 渠道路径 | `POST /upChannelApi/getPhoneByToken` |
| 华安路径 | `POST http://47.97.156.18:9040/common/channel/api/getPhoneByToken` |
| 业务码 | `500` |
| 耗时 | 0 ms |
| 说明 | 占位 token；需 SDK 真实 userInformation+token |

**请求体**（签名开启时 `sign` 在 HTTP 请求头，body 不含 `key`/`sign`）：

```json
{
  "channelCode": "BLtJjF",
  "userInformation": "test",
  "token": "test",
  "timestamp": "1783324385407"
}
```

**响应体**：

```json
{
  "code": 500,
  "message": "token 校验失败：token校验失败[MSISDN_OWNERS Processor]",
  "data": null
}
```

## 复现命令

```bash
HUAAN_BASE_URL='http://47.97.156.18:9040' \
HUAAN_CHANNEL_CODE='BLtJjF' \
HUAAN_KEY='fb9ec7236b6c45b7bfd562672e0373ea' \
HUAAN_SIGN_ENABLED=true \
HUAAN_TEST_PHONE='***' \
HUAAN_TEST_NAME='***' \
HUAAN_TEST_ID_CARD='***' \
HUAAN_PROBE_SKIP_SMS=false \
HUAAN_TEST_SMS_CODE='3875' \
go run ./scripts/huaan_direct_probe/
```
