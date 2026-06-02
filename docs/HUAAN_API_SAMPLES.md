# 华安接口响应样例

本文档保存华安上游接口的**原始 JSON 响应**，便于对照字段格式与联调。数据来自直连华安环境的集成测试，**非**本服务加工后的渠道响应（渠道侧三要素字段会加密，此处为华安明文）。

## getBankList — 获取银行列表

| 项 | 值 |
|----|-----|
| 路径 | `POST {HUAAN_BASE_URL}/upChannelApi/getBankList` |
| 采集时间 | 2026-06-02 |
| 环境 | `https://zf.ins.api.york.xin/` |
| 渠道编码 | `I7fZcM` |
| 签名 | 关闭（`sign_enabled=false`，请求体不含 `key`/`sign`） |

### 请求体示例

```json
{
  "timestamp": "1717300000000",
  "channelCode": "I7fZcM"
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

## 相关文档

- [TESTING.md](./TESTING.md) — 华安直连集成测试
- [ARCHITECTURE.md](./ARCHITECTURE.md) — getBankList 缓存与转发
- [CHANNEL_INTEGRATION.md](./CHANNEL_INTEGRATION.md) — 渠道侧调用说明
