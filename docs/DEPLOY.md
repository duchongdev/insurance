# 部署文档

本文档说明 insurance-bridge 的完整部署流程，涵盖**测试环境**与**生产环境**。部署采用 Docker Compose，无需在服务器安装 Go。

**阿里云上云**（ECS 采购、安全组、域名、HTTPS）见 **[ALIYUN.md](./ALIYUN.md)**。

---

## 1. 部署架构

Compose 编排四个服务，对外仅暴露 Nginx 的 80/443 端口：

```
                    ┌─────────────┐
  渠道 / 浏览器 ──► │   Nginx     │ :80 / :443
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   bridge    │ Go 应用（渠道 API + 管理后台 API）
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──────┐ ┌───▼───┐ ┌──────▼──────┐
       │   MySQL     │ │ Redis │ │  华安上游   │
       │ (MariaDB)   │ │       │ │  (外网)     │
       └─────────────┘ └───────┘ └─────────────┘
```

| 容器 | 说明 | 对外端口 |
|------|------|----------|
| `ins-nginx` | 反向代理、静态管理后台 | 80、443 |
| `ins-app` | Go 后端 | 无（内网） |
| `ins-mysql` | MariaDB 数据持久化 | 无（内网） |
| `ins-redis` | 缓存 | 无（内网） |

---

## 2. 测试环境与生产环境的区分

项目**没有**全局 `ENV=test|prod` 开关，通过以下三层区分：

| 层级 | 测试环境 | 生产环境 |
|------|---------|---------|
| **部署实例** | 内网测试机 `10.41.61.41:/home/ins` | 独立 ECS（如阿里云） |
| **配置文件** | 各服务器独立的 `.env`、`config/config.yaml` | 同上，密码/密钥/域名不同 |
| **华安上游** | 管理后台「华安配置」`envType=test`，绑定测试域名 | `envType=prod`，绑定生产域名 |
| **渠道** | 测试渠道关联测试华安配置 | 生产渠道关联生产华安配置 |
| **日志** | `BRIDGE_LOG_LEVEL=debug` | `info` |
| **HTTPS** | 可不配，HTTP :80 联调 | 必配 `deploy/ssl/` 证书 |

华安域名示例（以运营提供为准）：

| 环境 | 地址示例 |
|------|---------|
| 华安测试 | `http://47.97.156.18:9040` |
| 华安生产 | `https://ins.api.hahealth.ink/` |

**同一 bridge 实例可同时存在测试、生产两套华安配置**，各渠道通过 `huaan_setting_id` 绑定对应上游，互不影响。

---

## 3. 环境要求

| 项目 | 要求 |
|------|------|
| 操作系统 | Linux x86_64（Ubuntu 20.04+ / CentOS 7+） |
| Docker | 20.10+ |
| Docker Compose | v2（`docker compose`）或 v1（`docker-compose`） |
| Node.js | 18+（仅服务器**本地构建**管理后台时需要；`make deploy-test` 已在开发机构建好前端） |
| 网络 | 首次构建需拉取基础镜像；测试机若无法访问 Docker Hub，须用 `make deploy-test` 预构建镜像 |
| 端口 | 默认 **80**（HTTP）、**443**（HTTPS，需证书） |

验证：

```bash
docker version
docker compose version   # 或 docker-compose version
```

---

## 4. 部署方式概览

| 方式 | 适用场景 | 命令 |
|------|---------|------|
| **开发机一键部署测试环境** | 日常开发完成后更新测试机 | `make deploy-test` |
| **服务器手动部署** | 首次安装、生产环境、无 SSH 免密时 | 见 §5 |
| **仅打交付包** | 交付给第三方自行部署 | `make package` |

---

## 5. 首次部署（服务器手动）

### 5.1 上传并解压

```bash
mkdir -p /home/ins          # 测试环境默认目录；生产可改为 /opt/insurance-bridge
cd /home/ins
tar -xzf insurance-bridge-*.tar.gz
```

解压后应包含 `docker-compose.yml`、`Dockerfile`、`scripts/`、`config/config.yaml.example`、`.env.example`、`deploy/admin-dist/` 等。

> 交付包**不包含** `.env` 与 `config/config.yaml`（含密钥，勿入 Git），须在本机生成。

### 5.2 初始化配置

```bash
chmod +x scripts/*.sh
./scripts/install.sh
```

脚本会：

- 从 `.env.example` 生成 `.env`
- 从 `config/config.yaml.example` 生成 `config/config.yaml`
- 创建 `logs/` 目录
- 根据 `MYSQL_*` 自动生成 `BRIDGE_DATABASE_DSN`
- 根据 `REDIS_PASSWORD` 自动生成 `BRIDGE_REDIS_PASSWORD`

### 5.3 修改配置（必做）

#### `.env` 必改项

| 变量 | 说明 |
|------|------|
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 |
| `MYSQL_USER` | 应用库用户名（默认 `bridge`） |
| `MYSQL_PASSWORD` | 应用库密码 |
| `MYSQL_DATABASE` | 数据库名（默认 `insurance_bridge`） |
| `REDIS_PASSWORD` | Redis 认证密码 |
| `BRIDGE_HUAAN_BASE_URL` | 华安上游域名（测试填测试地址） |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | **32 字节**加密密钥（渠道三要素 AES） |
| `BRIDGE_SECURITY_JWT_SECRET` | 管理后台 JWT 密钥 |

**无需手填**（脚本自动生成）：

- `BRIDGE_DATABASE_DSN` ← `MYSQL_*`
- `BRIDGE_REDIS_PASSWORD` ← `REDIS_PASSWORD`

#### `.env` 可选项

| 变量 | 说明 | 测试建议 | 生产建议 |
|------|------|---------|---------|
| `BRIDGE_LOG_LEVEL` | 日志级别 | `debug` | `info` |
| `BRIDGE_HUAAN_SIGN_ENABLED` | 是否生成华安 sign | 按联调要求 | 通常 `true` |
| `BRIDGE_ADMIN_CORS_ORIGINS` | 管理后台跨域 | 本地 dev 时需要 | Nginx 同源留空 |
| `BRIDGE_IMAGE` | 应用镜像名 | `insurance-bridge:deploy` | 同左或自定义 |
| `NGINX_IMAGE` | Nginx 镜像 | `nginx:1.26.1-alpine` | 同左 |
| `REDIS_IMAGE` | Redis 镜像 | `redis:7-alpine` | 同左 |
| `MYSQL_IMAGE` | MariaDB 镜像 | `mariadb:10.6` | 同左 |

#### `config/config.yaml` 建议同步

| 配置段 | 关键项 | 说明 |
|--------|--------|------|
| `huaan` | `base_url`、`channel_code`、`key`、`sign_enabled` | 与 `.env` 中华安相关项保持一致 |
| `security` | `data_encryption_key`、`jwt_secret` | 必须与 `.env` 一致 |
| `log` | `level` | 测试 `debug`，生产 `info` |
| `server` | `mode`、`upstream_timeout` | 默认 `release`、25s |

> `BRIDGE_*` 环境变量会覆盖 yaml 同名项；Docker 部署以 `.env` 为主，`config/config.yaml` 作为补充。

修改 MySQL 或 Redis 密码后，执行 `./scripts/install.sh` 或 `./scripts/start.sh` 重新同步 DSN。

### 5.4 启动服务

```bash
./scripts/start.sh
```

首次启动若本地无预构建镜像，会构建 Docker 镜像（约数分钟）。`start.sh` 会：

- 同步 DSN / Redis 密码
- 检查 `deploy/admin-dist/`（交付包已含；缺失时尝试构建）
- 检测 `deploy/ssl/` 证书，自动切换 HTTP / HTTPS 模式
- 启动全部容器

---

## 6. 开发机一键部署测试环境

适用于已配置 SSH 免密（`ssh root@10.41.61.41`）的开发机。脚本在本机构建 linux/amd64 镜像并上传，**测试服务器无需访问 Docker Hub**。

```bash
make deploy-test
# 等价于
bash scripts/deploy-test.sh
```

流程（8 步）：

1. 打包源码（含预构建管理后台）
2. 本机构建 `insurance-bridge:deploy` 镜像
3. 导出应用镜像与 Redis 镜像
4. SSH 停服（`scripts/stop.sh`）
5. 上传压缩包与镜像
6. 解压、加载镜像、写入 `.env` 镜像变量
7. 启动（`--no-build`，不重建镜像）
8. 健康检查 `curl http://127.0.0.1/health/ready`

**前提**：测试服务器上已有 `.env` 与 `config/config.yaml`（首次须按 §5 手动初始化）。增量部署**保留**服务器现有配置，不覆盖密钥。

可选环境变量：

```bash
DEPLOY_HOST=root@10.41.61.41 \
DEPLOY_DIR=/home/ins \
NGINX_IMAGE=nginx:1.26.1-alpine \
make deploy-test
```

部署完成后：

- 管理后台：http://10.41.61.41/admin/
- 健康检查：http://10.41.61.41/health/ready

---

## 7. 部署后业务配置

配置文件解决「本服务怎么跑」；以下在**管理后台**完成「业务怎么对接」。

### 7.1 登录管理后台

- 地址：`http://<服务器IP>/admin/`（HTTPS 部署后改用 `https://<域名>/admin/`）
- 默认账号：**admin** / **Admin123!@#**（首次启动自动创建，**生产环境务必修改密码**）

### 7.2 华安配置

路径：管理后台 → **华安配置**

| 字段 | 说明 |
|------|------|
| 名称 | 便于识别，如「华安测试环境」 |
| 环境类型 | `test` 或 `prod` |
| 上游地址 | 华安 `base_url` |
| 渠道编码 | 华安侧 `channelCode` |
| 渠道密钥 | 华安侧 `key`（请求体字段） |

首次启动若库内无记录，系统会从 `config/config.yaml` 初始化一条默认测试配置。

### 7.3 渠道管理

路径：管理后台 → **渠道管理**

创建渠道时需配置：

- **关联华安配置**：选择测试或生产环境
- **callbackUrl**：投保结果回调地址（必填）
- **huaAnKey**：华安提供的渠道密钥
- **piiEncrypted**：三要素是否加密（新建默认开启）

创建后获得本服务的 `channelCode` 与 `channelKey`，提供给下游渠道商对接。详见 [CHANNEL_API.md](./CHANNEL_API.md)。

---

## 8. HTTPS（生产必做）

1. 在阿里云等平台申请 SSL 证书，下载 **Nginx** 格式。
2. 上传到服务器：

```text
deploy/ssl/fullchain.pem   # 证书链
deploy/ssl/privkey.pem     # 私钥
chmod 600 deploy/ssl/privkey.pem
```

3. 重启：`./scripts/start.sh`

检测到证书后自动：

- 启用 Nginx **443**（HTTPS）
- **80** 301 跳转至 HTTPS

4. 安全组放行 **443**、**80**；**不要**对公网开放 3306、6379。

无证书时默认仅 **HTTP :80**，便于测试联调。

---

## 9. 验证

```bash
# 容器状态
docker compose ps

# 健康检查
curl http://127.0.0.1/health/live
curl http://127.0.0.1/health/ready

# 应用日志
docker compose logs -f bridge

# Nginx 日志
docker compose logs -f nginx
```

浏览器：

| 资源 | 地址 |
|------|------|
| 管理后台 | `http://<IP>/admin/` |
| OpenAPI | `http://<IP>/openapi.yaml` |
| 渠道 API | `http://<IP>/upChannelApi/...` |

---

## 10. 版本更新

### 开发机推送（测试环境）

```bash
make deploy-test
```

### 服务器手动更新

```bash
./scripts/stop.sh
# 解压新版本覆盖代码（保留 .env、config/config.yaml、logs/、deploy/ssl/）
./scripts/start.sh
```

> **切勿覆盖** `.env` 与 `config/config.yaml`，否则密钥与数据库连接会丢失。

---

## 11. 常用运维命令

```bash
./scripts/stop.sh              # 停止服务（保留数据库数据）
./scripts/start.sh             # 启动 / 更新后重启
./scripts/install.sh           # 重新初始化或校验配置
docker compose logs -f bridge  # 应用日志
docker compose logs -f nginx   # Nginx 日志
docker compose logs -f mysql   # 数据库日志
docker compose logs -f redis   # Redis 日志
docker compose down -v         # 停止并删除数据卷（清空库，慎用）
```

应用日志文件：`logs/app.log`（宿主机映射）。

---

## 12. 配置项完整参考

### 12.1 `.env` 全部变量

| 变量 | 必填 | 说明 |
|------|------|------|
| `MYSQL_ROOT_PASSWORD` | 是 | MySQL root 密码 |
| `MYSQL_DATABASE` | 否 | 默认 `insurance_bridge` |
| `MYSQL_USER` | 否 | 默认 `bridge` |
| `MYSQL_PASSWORD` | 是 | 应用库密码 |
| `REDIS_PASSWORD` | 是 | Redis 密码 |
| `BRIDGE_HUAAN_BASE_URL` | 是 | 华安域名 |
| `BRIDGE_HUAAN_API_PATH` | 否 | 默认 `/upChannelApi` |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | 是 | 32 字节 |
| `BRIDGE_SECURITY_JWT_SECRET` | 是 | JWT 密钥 |
| `BRIDGE_DATABASE_DSN` | 自动 | 勿手填 |
| `BRIDGE_REDIS_PASSWORD` | 自动 | 勿手填 |
| `BRIDGE_LOG_LEVEL` | 否 | `debug` / `info` |
| `BRIDGE_HUAAN_CHANNEL_CODE` | 否 | 可被管理后台覆盖 |
| `BRIDGE_HUAAN_KEY` | 否 | 华安 body `key` 字段 |
| `BRIDGE_HUAAN_SIGN_ENABLED` | 否 | 默认 `false` |
| `BRIDGE_SERVER_UPSTREAM_TIMEOUT` | 否 | 默认 25s |
| `BRIDGE_ADMIN_CORS_ORIGINS` | 否 | 逗号分隔 |
| `BRIDGE_IMAGE` | 否 | 预构建镜像部署时使用 |
| `NGINX_IMAGE` / `REDIS_IMAGE` / `MYSQL_IMAGE` | 否 | 镜像覆盖 |

### 12.2 `config/config.yaml` 结构

```yaml
server:
  port: 8080
  mode: release
  read_timeout: 15s
  write_timeout: 30s
  upstream_timeout: 25s

huaan:
  base_url: https://华安域名/
  api_path: /upChannelApi
  channel_code: ""
  key: ""
  sign_enabled: false

database:
  dsn: ...   # Docker 部署由 BRIDGE_DATABASE_DSN 覆盖

redis:
  addr: redis:6379
  password: ...
  db: 0
  bank_list_ttl: 0

security:
  data_encryption_key: "32字节密钥"
  jwt_secret: "..."

log:
  level: info
  file_path: ./logs/app.log
  retention_days: 90
  archive_enabled: true
  max_size_mb: 100

admin:
  cors_origins: ""
```

---

## 13. 安全建议

- 勿将 `.env`、`config/config.yaml` 提交 Git 或外传
- 生产防火墙仅开放 **80、443**；SSH 仅放行运维 IP
- MySQL、Redis 不对公网映射端口
- `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` 生成后**勿改**（否则已加密数据无法解密）
- 管理后台默认密码须在生产环境修改
- 定期备份 Docker 卷 `mysql_data` 或导出 MySQL

---

## 14. 故障排查

| 现象 | 处理 |
|------|------|
| `请在 .env 中设置 MYSQL_ROOT_PASSWORD` | 未编辑 `.env`，执行 `./scripts/install.sh` |
| `配置仍为占位符` | `.env` 或 `config/config.yaml` 含「请修改」等占位文字 |
| `ready` 探针失败 | `docker compose logs mysql` / `redis` 查数据库是否就绪 |
| 构建失败 / 无法拉镜像 | 检查网络；测试机改用 `make deploy-test` |
| 80/443 无法访问 | 查安全组、Nginx 与 bridge 容器；HTTPS 查 `deploy/ssl/` |
| 缺少 `deploy/admin-dist` | 开发机执行 `bash scripts/build-admin.sh` 后重新打包 |
| 服务器缺少 Nginx 镜像 | 在 `.env` 设置已有镜像名，或手动 `docker load` |
| 华安接口 401/500 | 检查管理后台华安配置、IP 白名单、`sign_enabled` |

---

## 15. 目录说明

```
.
├── docker-compose.yml        # MySQL + Redis + bridge + Nginx
├── Dockerfile
├── .env.example              # 环境变量模板
├── config/
│   ├── config.yaml.example
│   └── config.yaml           # 实际配置（勿提交 Git）
├── deploy/
│   ├── admin-dist/           # 管理后台静态资源（打包时生成）
│   ├── nginx/                # Nginx 配置
│   └── ssl/                  # HTTPS 证书（不提交 Git）
├── scripts/
│   ├── install.sh            # 初始化配置
│   ├── start.sh              # 构建并启动
│   ├── stop.sh               # 停止
│   ├── package.sh            # 打交付包
│   ├── deploy-test.sh        # 一键部署测试环境
│   ├── sync-dsn.sh           # 同步数据库 DSN
│   └── sync-redis.sh         # 同步 Redis 密码
└── logs/                     # 应用日志
```

---

## 16. Docker 安装（新服务器）

### 16.1 安装 Docker（CentOS 示例）

```bash
yum remove -y docker docker-client docker-client-latest docker-common \
  docker-latest docker-latest-logrotate docker-logrotate docker-engine
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum install -y docker-ce docker-ce-cli containerd.io
systemctl start docker && systemctl enable docker
docker --version
```

### 16.2 安装 Docker Compose

```bash
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
  -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
docker-compose --version
```

### 16.3 配置镜像加速（国内推荐）

```bash
mkdir -p /etc/docker
tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://mirror.baidubce.com"
  ]
}
EOF
systemctl daemon-reload && systemctl restart docker
docker run --rm hello-world
```

---

## 17. 相关文档

| 文档 | 说明 |
|------|------|
| [ALIYUN.md](./ALIYUN.md) | 阿里云 ECS、安全组、域名 HTTPS |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | 系统架构 |
| [TESTING.md](./TESTING.md) | 华安直连集成测试（`.env.huaan`，与部署配置独立） |
| [CHANNEL_API.md](./CHANNEL_API.md) | 渠道商接口文档 |
| [../README.md](../README.md) | 项目说明与本地开发 |
