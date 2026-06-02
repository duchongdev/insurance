---
name: pre-commit-review
description: >-
  在 git commit 前对暂存/待提交改动做结构化 code review 并打分（满分 100，及格 60）。
  不及格时禁止提交并给出改造建议；及格后才允许 git add / git commit。
  当用户要求提交代码、Agent 准备 git commit、或用户提到提交前审查/code review/代码评审时使用。
---

# 提交前 Code Review

## 触发条件（满足任一即执行，**先于** `git commit`）

- 用户要求「提交代码 / git commit / 提交」
- Agent 按项目规范准备自动提交
- 用户明确要求 code review 或提交前审查

**硬约束**：未通过 review（总分 < 60）时 **不得** 执行 `git commit`；须先修复或向用户说明阻塞原因与改造建议。

## 执行流程

```
1. 收集 diff → 2. 构建检查 → 3. 多维度评分 → 4. 判定 → 5. 输出报告 → 6. 通过才 commit
```

### Step 1：收集改动

并行执行：

```bash
git status
git diff          # 未暂存
git diff --cached # 已暂存
```

阅读与改动相关的现有代码上下文；涉及 HTTP 接口时对照 `api/openapi.yaml` 与 `README.md`。

### Step 2：构建与测试（Go 项目）

```bash
go build ./...
go test ./...     # 或 make test
```

- 编译失败 → 直接判定 **不通过**，无需继续打分
- 测试失败：区分代码问题 vs 环境/架构问题；代码问题不通过，环境问题在报告中注明但仍需用户确认

### Step 3：评分维度（满分 100）

| 维度 | 分值 | 检查要点 |
|------|------|----------|
| **正确性** | 25 | 逻辑正确、无遗漏引用/死代码、依赖注入与路由一致、删除功能时无残留调用 |
| **架构与分层** | 20 | handler → service → repository 方向、职责边界、无循环依赖、改动位置符合 `.cursor/rules/project-standards.mdc` |
| **安全性** | 20 | 无密钥/PII 泄露、日志不记录敏感 body、鉴权未削弱、SQL 参数化 |
| **文档与契约同步** | 15 | 改 handler/OpenAPI 时 openapi.yaml + README 已同步；破坏性变更已 bump version |
| **代码质量** | 10 | 命名清晰、注释适度（见 `go-coding.mdc`）、最小 diff、无无关改动 |
| **提交原子性** | 10 | 单次 commit 只含同一功能/修复；skill/文档/业务代码不混无关变更 |

**扣分规则**：

- 🔴 阻塞问题（编译失败、逻辑 bug、安全漏洞、接口未同步文档）：该维度 ≤ 40%，总分通常不及格
- 🟡 建议改进（命名、注释、可维护性）：每处扣 1～3 分
- 🟢 可选优化：不扣分，仅在报告中列出

### Step 4：判定

| 总分 | 结论 |
|------|------|
| ≥ 60 | ✅ **通过**，可执行 `git commit` |
| < 60 | ❌ **不通过**，禁止提交 |

### Step 5：输出报告（固定格式）

```markdown
## 提交前 Code Review

**结论**：✅ 通过（{总分}/100） / ❌ 不通过（{总分}/100，及格线 60）

### 评分摘要
| 维度 | 得分 | 说明 |
|------|------|------|
| 正确性 | x/25 | … |
| 架构与分层 | x/20 | … |
| 安全性 | x/20 | … |
| 文档与契约同步 | x/15 | … |
| 代码质量 | x/10 | … |
| 提交原子性 | x/10 | … |

### 阻塞问题（必须修复）
- （无则写「无」）

### 改造建议（不通过时必填；通过时列可选优化）
1. …

### 验证情况
- go build：通过/失败
- go test：通过/失败/跳过（原因）
```

### Step 6：后续动作

**通过**：

1. 按 `.cursor/rules/project-standards.mdc` 撰写中文 commit message
2. `git add` 相关文件 → `git commit`
3. 若本次为可交付功能/修复，按 `deploy-test` skill 部署测试环境（用户说先不部署除外）

**不通过**：

1. **不要** `git commit`
2. 优先修复阻塞问题；修复后 **重新执行本 skill 全流程**
3. 向用户展示报告与改造建议，必要时询问是否继续修复

## 项目特定参考

- 架构与 Git 约定：`.cursor/rules/project-standards.mdc`
- Go 风格：`.cursor/rules/go-coding.mdc`
- API 文档同步：`.cursor/rules/api-documentation.mdc`
- 管理后台改动需检查 `web/admin/` 路由、API 客户端、页面是否一致

## 示例

**不通过**：删除了 handler 路由但未同步 openapi.yaml → 文档与契约同步 ≤ 5/15，总分不及格 → 禁止提交，建议补 openapi 与 README。

**通过**：移除业务落库逻辑，handler/service/model/repository/前端/OpenAPI/README 同步更新，`go build ./...` 通过 → 可提交 `refactor: 移除业务数据落库，渠道 API 改为透明转发`。
