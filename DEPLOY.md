# 交付部署说明

本文档面向**收到压缩包后在服务器上部署**的使用方。无需安装 Go，仅需 Docker 与 Docker Compose。

## 环境要求

| 项目 | 要求 |
|------|------|
| 操作系统 | Linux x86_64（推荐 Ubuntu 20.04+ / CentOS 7+） |
| Docker | 20.10+ |
| Docker Compose | v2（`docker compose`）或 v1（`docker-compose`） |
| Node.js | 18+（构建管理后台前端；`scripts/start.sh` 会自动调用 `scripts/build-admin.sh`） |
| 网络 | 首次构建需访问外网拉取基础镜像与 Go 依赖 |
| 端口 | 默认占用 **5051**（Nginx 反向代理，可在 `.env` 修改 `APP_PORT`） |

验证 Docker：

```bash
docker version
docker compose version   # 或 docker-compose version
```

## 部署步骤（四步）

### 1. 上传并解压

```bash
mkdir -p /opt/insurance-bridge
cd /opt/insurance-bridge
tar -xzf insurance-bridge-*.tar.gz
```

解压后目录应包含 `docker-compose.yml`、`Dockerfile`、`scripts/`、`config/config.yaml.example`、`.env.example` 等。

### 2. 初始化配置

```bash
chmod +x scripts/*.sh
./scripts/install.sh
```

脚本会：

- 从 `.env.example` 生成 `.env`
- 从 `config/config.yaml.example` 生成 `config/config.yaml`
- 创建 `logs/` 目录
- 根据 MySQL 账号自动同步 `BRIDGE_DATABASE_DSN`

### 3. 修改配置（必做）

编辑 **`.env`**（主配置，含密钥与数据库密码）：

| 变量 | 说明 |
|------|------|
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 |
| `MYSQL_USER` | 应用数据库用户名（默认 `bridge`） |
| `MYSQL_PASSWORD` | 应用数据库用户密码 |
| `MYSQL_DATABASE` | 数据库名（默认 `insurance_bridge`） |
| `REDIS_PASSWORD` | Redis 认证密码 |
| `BRIDGE_HUAAN_BASE_URL` | 华安上游域名 |
| `BRIDGE_SECURITY_DATA_ENCRYPTION_KEY` | **32 字节**加密密钥 |
| `BRIDGE_SECURITY_JWT_SECRET` | 管理后台 JWT 密钥 |
| `BRIDGE_ADMIN_DEFAULT_PASSWORD` | 首次创建的管理员密码 |
| `APP_PORT` | Nginx 对外端口，默认 5051 |

编辑 **`config/config.yaml`**（可与 `.env` 保持一致，尤其 `huaan.base_url`、`security` 段）。

`BRIDGE_DATABASE_DSN` **无需手动填写**；`install.sh` / `start.sh` 会根据 `MYSQL_USER`、`MYSQL_PASSWORD`、`MYSQL_DATABASE` 自动生成。

`BRIDGE_REDIS_PASSWORD` **无需手动填写**；`install.sh` / `start.sh` 会根据 `REDIS_PASSWORD` 自动生成。

修改 MySQL 或 Redis 密码后，执行 `./scripts/install.sh` 或 `./scripts/start.sh` 即可重新同步。

### 4. 启动服务

```bash
./scripts/start.sh
```

首次启动会构建 Docker 镜像（约数分钟）；若 `web/admin/dist/` 不存在，会先构建 Vue 管理后台。

也可提前构建前端：

```bash
bash scripts/build-admin.sh
./scripts/start.sh
```

## 验证

```bash
# 容器状态
docker compose ps

# 健康检查
curl http://127.0.0.1:5051/health/live
curl http://127.0.0.1:5051/health/ready

# 应用日志
docker compose logs -f bridge

# Nginx 访问/错误日志
docker compose logs -f nginx
```

浏览器访问：

- 管理后台：`http://<服务器IP>:5051/admin/`
- API 文档：`http://<服务器IP>:5051/openapi.yaml`

默认管理员用户名见 `config/config.yaml` 中 `admin.default_username`（默认 `admin`），密码为你在 `.env` 中配置的 `BRIDGE_ADMIN_DEFAULT_PASSWORD`。**首次登录后请修改密码。**

## 常用运维命令

```bash
./scripts/stop.sh              # 停止服务（保留数据库数据）
./scripts/start.sh             # 启动/更新后重启
docker compose logs -f bridge  # 查看应用日志
docker compose logs -f nginx   # 查看 Nginx 日志
docker compose logs -f mysql   # 查看数据库日志
docker compose logs -f redis   # 查看 Redis 日志
docker compose down -v         # 停止并删除数据卷（清空数据库与 Redis，慎用）
```

## 更新版本

```bash
./scripts/stop.sh
# 解压新版本覆盖代码（保留 .env、config/config.yaml、logs/）
./scripts/start.sh
```

## 安全建议

- 不要将 `.env` 或含真实密钥的 `config/config.yaml` 外传或提交到 Git
- 生产环境在防火墙仅开放 `APP_PORT`（Nginx）；应用与数据库默认不映射到宿主机端口
- 可在 `deploy/nginx/conf.d/bridge.conf` 中调整反向代理参数；生产环境建议在 Nginx 前或之上配置 HTTPS
- 定期备份 Docker 卷 `mysql_data` 或导出 MySQL 数据

## 故障排查

| 现象 | 处理 |
|------|------|
| `请在 .env 中设置 MYSQL_ROOT_PASSWORD` | 未创建或未编辑 `.env`，执行 `./scripts/install.sh` |
| `ready` 探针失败 | 数据库或 Redis 未就绪，执行 `docker compose logs mysql` / `docker compose logs redis` |
| 构建失败 / 无法拉镜像 | 检查服务器网络与 Docker 镜像源 |
| 5051 无法访问 | 检查防火墙、`APP_PORT`、Nginx 与 bridge 容器是否运行 |

## 目录说明

```
.
├── DEPLOY.md                 # 本文件
├── docker-compose.yml        # 编排 MySQL + Redis + 应用 + Nginx
├── deploy/nginx/             # Nginx 反向代理配置
├── Dockerfile                # 应用镜像构建
├── .env.example              # 环境变量模板
├── config/
│   ├── config.yaml.example   # 应用配置模板
│   └── config.yaml           # 实际配置（install 后生成，勿外传）
├── scripts/
│   ├── install.sh            # 初始化配置
│   ├── start.sh              # 构建并启动
│   ├── stop.sh               # 停止
│   ├── sync-redis.sh         # 同步 Redis 密码到 BRIDGE_REDIS_PASSWORD
│   └── package.sh            # 交付方打包容器（开发用）
└── logs/                     # 应用日志目录
```
---
## 一键安装 Docker 24.0.9（24.0.x 最新稳定小版本）
```
# 卸载旧版本（如果有）
yum remove -y docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine

# 安装依赖
yum install -y yum-utils device-mapper-persistent-data lvm2

# 添加阿里云 Docker 官方源（国内速度最快）
yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

# 安装 Docker
yum install -y docker-ce docker-ce-cli containerd.io

# 启动 Docker 并设置开机自启
systemctl start docker
systemctl enable docker

# 验证安装
docker --version 
```
##  一键安装 Docker Compose V2 2.24.5（和 Docker 24.0.x 最兼容、企业用得最多）。docker compose up -d（无横线）。
```
# 下载 Docker Compose 最新稳定版
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# 添加执行权限
chmod +x /usr/local/bin/docker-compose

# 创建软链接（让命令全局可用）
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

# 验证安装
docker-compose --version
```

## 配置阿里云镜像加速（必须配置，否则拉取镜像极慢）
```
# 创建 docker 配置目录
mkdir -p /etc/docker

# 写入阿里云镜像加速地址
tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://mirror.baidubce.com"
  ]
}
EOF

# 重启 Docker 生效
systemctl daemon-reload
systemctl restart docker
```

## 测试是否安装成功（跑一个 hello-world）
```
docker run --rm hello-world
```

