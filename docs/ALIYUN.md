# 阿里云上云说明

本文档汇总 insurance-bridge **单机单节点**部署到阿里云的方案，涵盖资源采购、网络与安全组、SSH 运维、应用部署及投产检查。Docker 安装与 `.env` 配置细节见 [DEPLOY.md](./DEPLOY.md)。

---

## 1. 部署架构

### 1.1 阶段一：单机（当前推荐）

```
                    互联网（渠道 / 管理后台 / 华安回调）
                              │
                    ┌─────────▼─────────┐
                    │  域名 + HTTPS      │  投产：api.example.com
                    │  （DNS → 公网 IP） │
                    └─────────┬─────────┘
                              │ :443 / :80
                    ┌─────────▼─────────────────────────┐
                    │  ECS 单机（1 台）                   │
                    │  Docker Compose                     │
                    │   nginx  → bridge → 华安上游        │
                    │   mysql / redis（仅容器内网）       │
                    └─────────────────────────────────────┘
```

### 1.2 阶段二：扩容（后期按需）

```
互联网 → 域名 → SLB（HTTPS）→ ECS-1 / ECS-2（bridge + nginx）
                              ↘ RDS MySQL + 云 Redis
```

**原则**：第一阶段 **不买 SLB、不买 RDS/云 Redis**，全部跑在一台 ECS 的 Compose 内；流量或可用性要求提高后再拆。

---

## 2. 需要采购什么

### 2.1 必买 / 必配

| 资源 | 建议规格 | 说明 |
|------|----------|------|
| **ECS** | 4 核 8G；系统盘 40G + 数据盘 100G SSD | 同时跑 MySQL、Redis、应用、Nginx；2C4G 仅适合极低流量测试 |
| **公网带宽** | 固定带宽 **3～5 Mbps** 起步 | 与公网 IP 一起在买 ECS 时勾选「分配公网 IPv4」 |
| **域名** | 1 个 | 投产 HTTPS 与渠道对接强烈建议；联调阶段可先用 IP |
| **SSL 证书** | 阿里云免费 DV 证书 | 绑域名，Nginx 格式部署 |

### 2.2 第一阶段不需要

| 资源 | 原因 |
|------|------|
| **SLB 负载均衡** | 单机无意义；多节点时再加 |
| **独立 EIP** | 买 ECS 时已带公网 IP 即可；仅当需 IP 与实例解绑、迁移时再转 EIP |
| **RDS / 云 Redis** | Compose 内自建；多节点或高可用时再迁 |
| **CDN** | 渠道 API 不适合 CDN |
| **非标准端口 5051** | 已改为标准 **80 / 443** |

### 2.3 地域

- 选离渠道用户、华安上游较近的地域（如华东1 杭州、华东2 上海）。
- ECS、域名解析、证书、后续 SLB 须 **同地域**（证书与 SLB 绑定在同 region）。

---

## 3. 公网 IP 与公网带宽

二者 **不是同一个东西**：

| 概念 | 含义 | 示例 |
|------|------|------|
| **公网 IP** | 服务器在互联网上的地址 | `47.116.192.222` |
| **公网带宽** | 该地址对外收发数据的速率上限 | `3 Mbps`（按固定带宽计费） |

- **公网 IP** = 别人能找到你  
- **公网带宽** = 找到你之后数据传输有多快  

买 ECS 时勾选「分配公网 IPv4」即同时获得 IP 与带宽，**一般无需单独再买 EIP**。

控制台「转换为弹性公网 IP」：将普通公网 IP 转为可独立绑定的 EIP，便于以后换机器保留 IP；单机初期可不做。

---

## 4. 域名

### 4.1 是否必须

| 场景 | 是否买域名 |
|------|------------|
| 内网/短期联调 | 可不买，用 IP |
| 正式对渠道开放、HTTPS | **建议买** |
| 长期生产 | **建议买** |

### 4.2 后期绑定当前公网 IP

**可以。** 常见顺序：

1. 现在：用 `http://<公网IP>/` 部署、联调  
2. 以后：买域名 → DNS **A 记录** 指向同一公网 IP → 申请 SSL → 启用 HTTPS  
3. **无需换服务器、无需换公网 IP**（实例不释放的前提下）

示例 DNS：

| 类型 | 主机记录 | 记录值 |
|------|----------|--------|
| A | `api` | `47.116.192.222` |

投产访问地址：

- 管理后台：`https://api.yourdomain.com/admin/`
- 渠道 API：`https://api.yourdomain.com/upChannelApi/...`
- 健康检查：`https://api.yourdomain.com/health/ready`

---

## 5. 对外端口

项目 **不再使用 5051**，统一为标准端口：

| 端口 | 用途 | 安全组 |
|------|------|--------|
| **443** | HTTPS（投产主入口） | 按白名单或 `0.0.0.0/0` |
| **80** | HTTP → HTTPS 跳转 | 同上 |
| **22** | SSH 运维 | **仅运维 IP** |
| 3306 / 6379 | MySQL / Redis | **禁止对公网开放** |

无 SSL 证书时：`./scripts/start.sh` 仅监听 **HTTP 80**。  
证书放入 `deploy/ssl/` 后重启，自动启用 **443** 并将 80 跳转 HTTPS。详见 [DEPLOY.md § HTTPS](./DEPLOY.md)。

---

## 6. 安全组

### 6.1 推荐入方向规则（投产）

| 协议 | 端口 | 授权对象 | 说明 |
|------|------|----------|------|
| TCP | 443 | 见 §7 白名单 | HTTPS |
| TCP | 80 | 见 §7 白名单 | 跳转 HTTPS（可与 443 同源） |
| TCP | 22 | `你的公网IP/32` | SSH，如 `220.181.41.13/32` |

### 6.2 应删除的规则

| 规则 | 原因 |
|------|------|
| **全部 RDP 3389** | Linux 不用远程桌面 |
| **ICMP 0.0.0.0/0** | 非必须，减少探测 |
| **5051** | 已废弃 |
| **SSH 22 → 大段阿里云内网 IP** | 改为仅自己的运维 IP；Console Workbench 可兜底 |

### 6.3 操作顺序（避免锁死）

1. **先新增** SSH：`22` → `你的IP/32`  
2. 本机验证 SSH 能登录  
3. **再删除** 3389、旧 22 规则、ICMP 等  
4. 配置 80/443 白名单后，删除 `0.0.0.0/0` 的 80/443（若曾开放）

---

## 7. 上下游 IP 白名单

安全组 **入方向** 控制「谁连到你 ECS 的 80/443」；**出站** 由你访问对方时，在 **对方** 白名单登记你的公网 IP。

### 7.1 流量方向对照

| 流量 | 方向 | 配置位置 |
|------|------|----------|
| 渠道调 `/upChannelApi/...` | 渠道 → 你（**入站**） | **你的安全组**：放行渠道 IP |
| 华安回调 `/huaan/callback/...` | 华安 → 你（**入站**） | **你的安全组**：放行华安回调 IP |
| 你调华安上游 API | 你 → 华安（**出站**） | **华安侧**白名单：填你的公网 IP |
| 你转发回调到渠道 `callbackUrl` | 你 → 渠道（**出站**） | **渠道侧**白名单：填你的公网 IP |
| 管理后台浏览器访问 | 你/同事 → 你（**入站**） | **你的安全组**：加办公网 IP |

### 7.2 提供给上下游的出口 IP

将 ECS **公网 IP** 提供给华安与各渠道（示例）：

```text
47.116.192.222
```

### 7.3 80/443 能否改为具体 IP

**可以。** 在渠道 IP、华安回调 IP、办公网 IP 均能列全时，**删除** `0.0.0.0/0`，改为多条规则，每条一个来源 `/32` 或网段。  
注意：漏 IP、换宽带、新增渠道都需同步改安全组；先加新规则测通，再删 `0.0.0.0/0`。

---

## 8. SSH 连接服务器

### 8.1 基本命令

```bash
# 密钥登录（创建 ECS 时若选了密钥对）
chmod 400 ~/path/to/your-key.pem
ssh -i ~/path/to/your-key.pem root@47.116.192.222

# 密码登录（若实例启用了密码且 sshd 允许）
ssh root@47.116.192.222
```

Alibaba Cloud Linux 默认用户名为 **`root`**。

### 8.2 可选：SSH 配置别名

`~/.ssh/config`：

```text
Host aliyun-ins
    HostName 47.116.192.222
    User root
    IdentityFile ~/.ssh/aliyun-ecs.pem
```

之后：`ssh aliyun-ins`

### 8.3 常见报错

| 现象 | 含义 | 处理 |
|------|------|------|
| `Connection timed out` | 安全组未放行 22 或网络不通 | 检查安全组、实例是否运行中 |
| `Permission denied (publickey,...)` 且无 `password` | 已连上 SSH，**仅允许密钥** | 使用 `-i xxx.pem`；或 Workbench 登录后配置密码/公钥 |
| `Permission denied (password)` | 密码错误 | 控制台重置实例密码并重启 |

**Workbench 兜底**：控制台 → ECS → **远程连接** → Workbench，可重置密码或写入 `~/.ssh/authorized_keys`。

启用密码登录（Workbench 内执行，按需）：

```bash
sudo sed -i 's/^#\?PasswordAuthentication.*/PasswordAuthentication yes/' /etc/ssh/sshd_config
sudo sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config
sudo systemctl restart sshd
```

---

## 9. 部署步骤

### 9.1 服务器初始化

```bash
# 安装 Docker、Compose、镜像加速 — 见 DEPLOY.md
mkdir -p /home/ins
cd /home/ins
```

### 9.2 上传并解压

开发机：

```bash
make package
scp dist/insurance-bridge-*.tar.gz root@47.116.192.222:/home/ins/
```

ECS：

```bash
cd /home/ins
tar -xzf insurance-bridge-*.tar.gz
chmod +x scripts/*.sh
./scripts/install.sh
# 编辑 .env、config/config.yaml（密码、华安地址、32 字节密钥等）
./scripts/start.sh
```

若服务器无法访问 Docker Hub，开发机使用 `make deploy-test`（修改 `DEPLOY_HOST` 指向 ECS），见 `scripts/deploy-test.sh`。

### 9.3 验证

```bash
curl http://127.0.0.1/health/ready
curl http://127.0.0.1/health/live
docker compose ps
```

浏览器（安全组放行 80 后）：

- `http://47.116.192.222/admin/`
- `http://47.116.192.222/openapi.yaml`

### 9.4 启用 HTTPS（投产）

```bash
mkdir -p deploy/ssl
# 上传阿里云 Nginx 格式证书
# deploy/ssl/fullchain.pem
# deploy/ssl/privkey.pem
chmod 600 deploy/ssl/privkey.pem

./scripts/stop.sh && ./scripts/start.sh
```

验证：

```bash
curl -k https://127.0.0.1/health/ready   # 或改用域名
```

---

## 10. 配置要点（`.env`）

| 变量 | 说明 |
|------|------|
| `MYSQL_*` / `REDIS_PASSWORD` | 强随机密码 |
| `BRIDGE_HUAAN_BASE_URL` | 华安生产环境地址 |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | **32 字节**，生成后勿改 |
| `BRIDGE_SECURITY_JWT_SECRET` | 足够长的随机串 |

`BRIDGE_DATABASE_DSN`、`BRIDGE_REDIS_PASSWORD` 由脚本自动生成，勿手填。

---

## 11. 投产检查清单

### 网络与安全

- [ ] 安全组：443、80 按白名单或业务需要配置；22 仅运维 IP  
- [ ] 未对公网开放 3306、6379、5051、3389  
- [ ] `.env`、`config/config.yaml` 权限适当，不入 Git  

### 服务

- [ ] `https://<域名>/health/ready` 正常  
- [ ] 管理后台可登录，默认密码已修改  
- [ ] Nginx HTTPS 已启用（公网渠道须 TLS）  

### 业务与对接

- [ ] 管理后台创建渠道，记录 `channelCode` / `channelKey`  
- [ ] 渠道 `piiEncrypted=true`（公网默认须加密）  
- [ ] 已将 **47.116.192.222** 提供给华安、渠道做出口白名单  
- [ ] 安全组已放行华安回调 IP、各渠道 IP  
- [ ] 渠道文档平台域名改为 `https://api.yourdomain.com`  

---

## 12. 运维与更新

```bash
./scripts/stop.sh
# 解压新版本（保留 .env、config/config.yaml、logs/、deploy/ssl/）
./scripts/start.sh
```

| 操作 | 命令 |
|------|------|
| 日志 | `docker compose logs -f bridge nginx` |
| 备份 | 数据盘快照；或 `mysqldump` |
| 监控 | 阿里云云监控 CPU/内存/磁盘 |

---

## 13. 后期扩容

| 信号 | 动作 |
|------|------|
| CPU/内存长期偏高 | 垂直扩容 ECS |
| 需消除单点 | SLB + 2 台 ECS（bridge + nginx） |
| 多应用节点 | MySQL → RDS，Redis → 云 Redis |
| 固定入口、换机 | 公网 IP 绑 SLB |

---

## 14. 费用粗算（参考）

| 项目 | 大致范围 |
|------|----------|
| ECS 4C8G + 5M 带宽 + 100G 盘 | ¥300～600/月 |
| 域名 | ¥50～100/年 |
| SSL | 免费 DV |
| SLB / RDS / 云 Redis | 第一阶段 ¥0 |

---

## 15. 相关文档

| 文档 | 内容 |
|------|------|
| [DEPLOY.md](./DEPLOY.md) | Docker 安装、Compose 启动、HTTPS 证书路径 |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | 系统架构与数据流 |
| [CHANNEL_API.md](./CHANNEL_API.md) | 渠道对接、签名、回调 |
| [PII_CHANNEL_ENCRYPTION.md](./PII_CHANNEL_ENCRYPTION.md) | 公网须 TLS、PII 加密策略 |

---

## 16. 实施时间线（建议）

```
第 1 步  购买 ECS（带公网 IP），配置安全组（22 仅自己 IP）
第 2 步  SSH 登录，安装 Docker / Compose
第 3 步  上传交付包，install.sh，填写 .env / config.yaml
第 4 步  start.sh，http://公网IP/health/ready 通过
第 5 步  与华安/渠道交换 IP 白名单，收紧 80/443 入站规则
第 6 步  购买域名，DNS 解析，申请 SSL，deploy/ssl/ 后重启
第 7 步  渠道联调、投产检查清单逐项确认
```
