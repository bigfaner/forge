新的工作流程：
1. 用户使用 brainstorm skill 把模糊的想法变成明确的提案，输出proposal.md
2. 用户使用 write-prd skill 根据proposal.md 产出prd
3. 用户使用 tech-design skill 根据prd设计技术方案
4. 用户使用 ui-design skill 根据prd设计界面，产出UI设计稿
5. 用户使用 breakdown-tasks skill 根据prd、技术方案、UI设计拆分任务。
6. 用户使用 run-tasks / execute-task (slash command)执行任务
7. AI端到端测试+人工验收



新的目录结构：
```
project-root/
└── docs/
    ├── proposals/<slug>/           # /brainstorm 产出
    │   └── proposal.md             # 提案
    ├── features/<slug>/            # Feature 工作区
    │   ├── manifest.md             # Feature 索引 & 可追溯性映射
    │   ├── prd/
    │   │   ├── prd-spec.md         # PRD Spec（需求文档）
    │   │   ├── prd-user-stories.md # 用户故事
    │   │   └── prd-ui-functions.md # UI 功能要点（可选）
    │   ├── design/
    │   │   ├── tech-design.md      # 技术设计
    │   │   └── api-handbook.md     # API 文档
    │   ├── ui/
    │   │   └── ui-design.md        # UI 设计规格（可选）
    │   └── tasks/
    │       ├── index.json          # 任务定义（核心）
    │       ├── process/            # 运行时状态（不提交）
    │       │   ├── state.json      # 当前任务状态
    │       │   └── record.json     # 进行中的记录
    │       ├── 1.1-<title>.md     # 任务详情
    │       └── records/            # 执行记录
    │           └── 1.1-<title>.md
    └── lessons/                    # 经验教训
```
