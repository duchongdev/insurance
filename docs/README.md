# 项目文档

> **开发前必读**：[../README.md](../README.md) 为项目**全局要求与标准**（架构分层、编码规范、Git/部署约定等）。本文档目录为专题文档索引。

| 文档 | 说明 |
|------|------|
| [../README.md](../README.md) | **全局开发标准** + 项目说明、快速启动 |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | 系统架构：渠道层 / 华安层分离、数据流与代码结构 |
| [TESTING.md](./TESTING.md) | 测试说明：华安直连集成测试、全流程测试、环境变量 |
| [HUAAN_API_SAMPLES.md](./HUAAN_API_SAMPLES.md) | 华安上游 10 个接口原始 JSON 样例与联调结论 |
| [CHANNEL_API.md](./CHANNEL_API.md) | **下游渠道商接口文档**（对外）：10 个 API、签名、加解密、示例 |
| [CHANNEL_API_SAMPLES.md](./CHANNEL_API_SAMPLES.md) | 渠道请求/响应 JSON 样例归档（联调对照） |
| [DEPLOY.md](./DEPLOY.md) | 交付部署：Docker 安装、配置、启动与运维 |
| [ALIYUN.md](./ALIYUN.md) | **阿里云上云**：资源采购、安全组、SSH、域名 HTTPS、白名单 |
| [CHANNEL_INTEGRATION.md](./CHANNEL_INTEGRATION.md) | 已合并至 CHANNEL_API.md（保留跳转） |
| [PII_CHANNEL_ENCRYPTION.md](./PII_CHANNEL_ENCRYPTION.md) | **方案**：按渠道控制三要素是否加解密、Redis 渠道配置缓存（设计稿，未编码） |
