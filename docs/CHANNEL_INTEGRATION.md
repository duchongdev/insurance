# 渠道接入指南

## 1. 接入准备

1. 由平台管理员在管理后台创建渠道，获取：
   - `channelCode`：渠道编码
   - `channelKey`：渠道签名密钥（请求体 `key` 字段）
2. 向平台索取 **三要素加密密钥**（与服务器 `security.data_encryption_key` 一致，32 字节）。
3. 平台配置华安分配的 `huaAnKey`，用于本服务向上游签名（渠道无需感知）。

## 2. 接口地址

```
POST {平台域名}/upChannelApi/{接口路径}
Content-Type: application/json; charset=utf-8
```

接口路径与 `ZF保险.md` 完全一致。

## 3. 签名算法

与华安文档一致：

1. 除 `sign` 外所有参数按 key 的 ASCII 升序排序。
2. 拼接为 `key1=value1&key2=value2`（`key` 字段使用渠道密钥）。
3. 对拼接字符串做 MD5，32 位**小写**十六进制，写入 `sign`。

## 4. 用户三要素加密

以下字段在请求本服务时必须 **AES-256-GCM** 加密后 **Base64** 编码：

| 字段 | 说明 |
|------|------|
| phoneNo | 手机号 |
| name | 姓名 |
| idCard | 身份证号 |

本服务解密后以明文转发华安。华安响应中的同名敏感字段，本服务会加密后再返回渠道。

### 加密格式

- 算法：AES-256-GCM
- 密钥：32 字节 UTF-8 字符串
- Nonce：随机 12 字节，置于密文前
- 输出：`Base64(nonce || ciphertext+tag)`

Go 示例见 `internal/pkg/cipher/cipher.go`。

## 5. 响应约定

- `code == 200` 表示成功，其余为失败。
- 结构与原华安文档一致；仅三要素相关字段为密文。

## 6. 错误码（本服务网关层）

| code | 说明 |
|------|------|
| 400 | 参数/解密错误 |
| 401 | 签名或密钥错误 |
| 403 | 渠道无效或已禁用 |
| 502 | 华安不可达 |
| 504 | 调用华安超时 |

业务错误仍以华安返回的 `code`、`message` 为准。

## 7. 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) — 系统架构与分层
- [TESTING.md](./TESTING.md) — 测试说明
- [DEPLOY.md](./DEPLOY.md) — 部署与配置
