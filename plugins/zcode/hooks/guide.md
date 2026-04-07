# ZCode Guide

## Document Index

```
project-root/
└── docs/
    ├── features/<slug>/           # Feature 工作区
    │   ├── prd.md                 # 需求文档
    │   ├── design.md              # 设计文档
    │   └── tasks/
    │       ├── index.json         # 任务定义（核心）
    │       ├── process/           # 运行时状态（不提交）
    │       │   ├── state.json     # 当前任务状态
    │       │   └── record.json    # 进行中的记录
    │       ├── 1.1-<title>.md     # 任务详情
    │       └── records/               # 执行记录
    │           └── 1.1-<title>.md
    ├── README.md                   # 知识库索引 (本文件)
    ├── ARCHITECTURE.md             # 分层架构
    ├── DECISIONS.md                # 技术决策
    └── lessons/                    # 经验教训
```

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

> 完整命令说明请运行 `task -h` 或 `task [command] -h`
