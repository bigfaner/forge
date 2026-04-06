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
| `task record <id> --data file.json` | 任务完成后，生成执行记录并更新状态 |

> 完整命令列表和参数说明请运行 `task -h` 或 `task [command] -h`