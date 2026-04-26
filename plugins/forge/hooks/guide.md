# Forge Guide

## Document Index

```
project-root/
├── docs/
│   ├── proposals/<slug>/           # /brainstorm 产出
│   │   └── proposal.md
│   ├── features/<slug>/            # Feature 工作区
│   │   ├── manifest.md             # Feature 索引 & 可追溯性映射
│   │   ├── prd/
│   │   │   ├── prd-spec.md         # PRD Spec (需求文档)
│   │   │   ├── prd-user-stories.md # 用户故事
│   │   │   └── prd-ui-functions.md # UI 功能要点（可选）
│   │   ├── design/
│   │   │   ├── tech-design.md      # 技术设计
│   │   │   └── api-handbook.md     # API 文档
│   │   ├── ui/
│   │   │   └── ui-design.md        # UI 设计规格（可选）
│   │   ├── testing/                # 测试产物（由标准任务生成）
│   │   │   ├── test-cases.md       # 测试用例
│   │   │   ├── scripts/            # 可执行测试脚本
│   │   │   └── results/            # 测试结果和截图
│   │   └── tasks/
│   │       ├── index.json          # 任务定义（核心）
│   │       ├── process/            # 运行时状态（不提交）
│   │       │   ├── state.json      # 当前任务状态
│   │       │   └── record.json     # 进行中的记录
│   │       ├── 1.1-<title>.md     # 任务详情
│   │       └── records/            # 执行记录
│   │           └── 1.1-<title>.md
│   ├── sitemap/
│   │   └── sitemap.json            # 页面元素地图（项目级，由 /gen-sitemap 生成）
│   ├── README.md                   # 知识库索引 (本文件)
│   ├── ARCHITECTURE.md             # 分层架构
│   ├── decisions/                  # 技术决策（按类别分目录）
│   └── lessons/                    # 经验教训
├── tests/e2e/                      # 已毕业的 e2e 测试（回归套件）
│   ├── .graduated/                 # 毕业标记（每 feature 一个）
│   │   └── <slug>                  # 包含毕业时间戳
│   ├── config.yaml                 # 测试环境配置
│   ├── helpers.ts                  # 共享测试工具
│   ├── package.json                # e2e 测试依赖
│   ├── tsconfig.json               # TypeScript 配置
│   └── <target>/                   # 按 target 组织的 spec 文件
```

## Skill Workflow

```
/brainstorm → /write-prd → /eval-prd → /tech-design ─→ /eval-design → /breakdown-tasks
     ↓             ↓            ↓            ↓              ↓               ↓
 proposal.md   prd/*.{3}  eval report  design/*.{2}   eval report     tasks/*.md
               manifest.md             manifest.md                   manifest.md
                            ↘ /ui-design ─→ /eval-design ↗
                                 ↓
                            ui/ui-design.md

/breakdown-tasks 追加标准测试任务（T-test-1, T-test-2），执行时依次调用：
/gen-sitemap → /gen-test-cases → /gen-test-scripts → /run-e2e-tests
     ↓              ↓                  ↓                   ↓
 sitemap.json  test-cases.md    testing/scripts/*    testing/results/
```

每个 skill 执行前会 `ls` 检查上一阶段产物是否存在；缺失则中止并提示用户先完成上一步。

### Manifest

`manifest.md` 是 Feature 的单一入口，AI agent 读取此文件即可了解完整上下文：
- **Documents** 表：列出所有文档路径和自动生成的摘要
- **Traceability** 表：PRD → Design → Tasks 的追溯映射
- **Status**：prd → design → tasks → in-progress → done

## Task-CLI

Task CLI 管理 feature 生命周期中的任务流转。

**典型场景**: 开始工作前 `task feature` → `task claim` 认领任务 →  `task record` 记录任务结果 + 更新任务状态。

### Key Commands

| 操作 | 命令 | 说明 |
|------|------|------|
| 切换 feature | `task feature <slug>` | 设置当前工作上下文 |
| 认领任务 | `task claim` | 获取下一个可用任务 |
| 完成任务 | `task record <id> --data docs/features/{slug}/tasks/process/record.json` | **一步完成记录+状态更新** |

### `task record` 工作流

```
1. task claim           → 写入 process/state.json
2. 任务执行期间          → 写入 process/record.json
3. task record --data docs/features/{slug}/tasks/process/record.json  → 生成 records/*.md + 更新 index.json
```

**一条命令完成2件事：** 生成并输出记录文件 → 更新任务状态。

### record.json
**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `taskId` | string | ✓ | 任务 ID |
| `status` | string | | 状态，默认 `completed` |
| `summary` | string | ✓ | 实现摘要 |
| `filesCreated` | array | | 新建文件列表 |
| `filesModified` | array | | 修改文件列表 |
| `keyDecisions` | array | | 关键设计决策 |
| `testsPassed` | int | | 通过测试数 |
| `testsFailed` | int | | 失败测试数 |
| `coverage` | float | | 覆盖率 |
| `acceptanceCriteria` | array | | `{criterion, met}` 对象数组 |

**示例**
```json
{
  "taskId": "3.3.1",
  "status": "completed",
  "summary": "实现了什么功能",
  "filesCreated": ["src/components/Button.tsx"],
  "filesModified": ["src/utils/helpers.ts"],
  "keyDecisions": ["关键决策"],
  "testsPassed": 12,
  "testsFailed": 0,
  "coverage": 85.6,
  "acceptanceCriteria": [{ "criterion": "验收标准", "met": true }]
}
```

### 强制规则

**禁止操作：**
- ❌ 直接写入 `records/*.md` 或 `index.json`
- ❌ 用 Python/JavaScript/Node 修改 JSON
- ❌ 写入格式错误的 `process/record.json`

### Validation Rules（enforced by CLI）

`task record` 会拒绝以下组合：

| Condition | Error | Fix |
|-----------|-------|-----|
| `status=completed` + `testsPassed=0` + `testsFailed=0` + `coverage >= 0` | No test evidence | Run tests, or set `coverage: -1.0` for no-test tasks |
| `status=completed` + any `acceptanceCriteria.met=false` | Unmet AC | Fix the issue, or set `status: "blocked"` |
| `summary` is empty or whitespace | Missing summary | Provide a summary |

Override any validation error with `--force`:
```bash
task record <TASK_ID> --data record.json --force
```

> 完整命令说明请运行 `task -h` 或 `task [command] -h`
