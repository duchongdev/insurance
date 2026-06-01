---
name: deploy-test
description: >-
  打包并部署 insurance-bridge 到测试环境服务器（10.41.61.41:/home/ins）。
  先停服再上传解压启动，保留服务器 .env 与 config/config.yaml。
  当用户说部署、部署测试、测试环境、更新测试环境服务、发布测试、deploy test、
  更新测试服务器、上线测试 等时使用本 skill，并立即执行部署脚本。
---

# 测试环境自动部署

## 触发词

用户消息包含以下任一意图时**立即应用本 skill**（无需再确认，直接执行）：

- 部署 / 部署测试 / 部署到测试环境
- 测试环境 / 更新测试环境 / 更新测试环境服务
- 发布测试 / 上线测试 / 更新测试服务器

## 默认环境（见 `.cursor/rules/deploy.mdc`）

| 项 | 值 |
|----|-----|
| SSH | `root@10.41.61.41`（本机免密） |
| 目录 | `/home/ins/` |
| 原则 | **先停服，再覆盖代码，再启动**；不覆盖服务器 `.env`、`config/config.yaml` |

## 执行方式

**默认执行项目脚本**（不要临时拼命令）：

```bash
bash scripts/deploy-test.sh
```

或：

```bash
make deploy-test
```

脚本流程：本机 `make package` → SSH 停服 → `scp` 上传 → 远程解压 → `sync-dsn` + `start.sh` → `curl /health/ready`。

## Agent 行为

1. 读到触发词后**先读本 skill**，再**直接运行** `scripts/deploy-test.sh`。
2. 全程用 Shell 工具执行，不要只输出命令让用户自己跑。
3. 部署结束汇报：是否成功、健康检查结果、失败时的 `docker compose logs` 指引。
4. **不要**修改、上传或覆盖服务器上的 `.env`、`config/config.yaml`（压缩包本身也不含这两项）。
5. 未经用户明确要求，不要 `git commit`、不要部署到其他机器或目录。

## 首次部署（服务器尚无配置）

若脚本报错缺少 `.env` / `config/config.yaml`，引导用户在服务器执行：

```bash
ssh root@10.41.61.41
cd /home/ins
./scripts/install.sh
# 编辑 .env 与 config/config.yaml 后
./scripts/start.sh
```

之后再在本机用 `bash scripts/deploy-test.sh` 做增量更新。

## 失败排查

```bash
ssh root@10.41.61.41 'cd /home/ins && docker compose ps'
ssh root@10.41.61.41 'cd /home/ins && docker compose logs --tail=50 bridge'
ssh root@10.41.61.41 'curl -v http://127.0.0.1:5051/health/ready'
```

## 环境变量（可选覆盖）

| 变量 | 默认 |
|------|------|
| `DEPLOY_HOST` | `root@10.41.61.41` |
| `DEPLOY_DIR` | `/home/ins` |

```bash
DEPLOY_HOST=root@10.41.61.41 DEPLOY_DIR=/home/ins bash scripts/deploy-test.sh
```
