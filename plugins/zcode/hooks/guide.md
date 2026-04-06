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

**典型场景**: 开始工作前 `task feature` → `task claim` 认领任务 → 完成后 `task record` 记录。

### Key Commands

| Command | When to Use |
|---------|-------------|
| `task feature <slug>` | 切换或设置当前 feature 上下文 |
| `task claim` | 认领下一个可用任务（自动选择依赖已满足的最高优先级任务） |
| `task record <id> --data docs/features/{slug}/tasks/process/record.json` | 任务完成后，生成执行记录并更新状态 |

### record.json 生成机制

**工作流程：**
```
1. task claim           → 写入 process/state.json（当前任务状态）
2. 任务执行期间         → 写入 process/record.json（执行记录）
3. task record --data   → 读取 JSON，生成 records/*.md，清空 process/
```

**生成方式：** 在任务执行过程中，agent 将执行信息写入 `docs/features/{slug}/tasks/process/record.json`：

```json
{
  "taskId": "3.3.1",
  "status": "completed",
  "summary": "实现了什么功能",
  "filesCreated": ["src/components/Button.tsx"],
  "filesModified": ["src/utils/helpers.ts"],
  "keyDecisions": ["使用 useCallback 优化性能"],
  "testsPassed": 12,
  "testsFailed": 0,
  "coverage": 85.6,
  "acceptanceCriteria": [
    { "criterion": "按钮点击响应正确", "met": true }
  ]
}
```

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

**⚠️ 强制规则：**
- 唯一允许路径：`docs/features/{slug}/tasks/process/record.json`
- 必须使用 `task record` CLI 命令，禁止直接写入 `index.json` 或手动创建记录文件

> 完整命令列表和参数说明请运行 `task -h` 或 `task [command] -h`