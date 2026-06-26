---
name: deploy-test
description: >-
  打包并部署 insurance-bridge 到测试环境服务器（10.41.61.41:/home/ins）。
  先停服再上传解压启动，保留服务器 .env 与 config/config.yaml。
  当用户说部署、部署测试、测试环境、更新测试环境服务、发布测试、deploy test、
  更新测试服务器、上线测试、改造完自动更新远程服务器 等时使用本 skill，并立即执行部署脚本。
  在本仓库完成功能/修复/改造后，也应默认执行本 skill（见 .cursor/rules/deploy.mdc）。
---

# 测试环境自动部署

## 触发条件（满足任一即执行，无需再确认）

**用户明确要求：**

- 部署 / 部署测试 / 部署到测试环境
- 测试环境 / 更新测试环境 / 更新测试环境服务 / 更新远程服务器
- 发布测试 / 上线测试 / 更新测试服务器
- 改造完自动更新服务器 / 每次改完部署

**Agent 默认行为（本仓库）：**

- 完成可交付的功能、修复或管理后台改造后，**任务收尾时主动运行**本脚本（除非用户说先不部署，或改动仅文档注释无运行时影响）。

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

1. 触发后**先读本 skill**，再**直接运行** `scripts/deploy-test.sh`（或 `make deploy-test`）。
2. 全程用 Shell 工具执行，不要只输出命令让用户自己跑。
3. 部署结束汇报：是否成功、健康检查结果、失败时的 `docker compose logs` 指引。
4. 功能/改造类任务：代码改完 →（按需文档同步）→ **部署** → 再向用户总结（含部署结果）。
5. **不要**修改、上传或覆盖服务器上的 `.env`、`config/config.yaml`（压缩包本身也不含这两项）。
6. 未经用户明确要求，不要 `git commit`、不要部署到其他机器或目录。

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
ssh root@10.41.61.41 'curl -v http://127.0.0.1/health/ready'
```

## 环境变量（可选覆盖）

| 变量 | 默认 |
|------|------|
| `DEPLOY_HOST` | `root@10.41.61.41` |
| `DEPLOY_DIR` | `/home/ins` |

```bash
DEPLOY_HOST=root@10.41.61.41 DEPLOY_DIR=/home/ins bash scripts/deploy-test.sh
```
