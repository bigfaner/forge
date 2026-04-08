新的工作流程：
1. 用户使用 brainstorm skill 把模糊的想法变成明确的提案，输出proposal.md
2. 用户使用 write-prd skill 根据proposal.md 产出prd
3. 用户使用 tech-design skill 根据prd设计技术方案
4. 用户使用 ui-design skill 根据prd设计界面，产出UI设计稿
5. 用户使用 breakdown-tasks skill 根据prd、技术方案 overview.md、ui overivew.md 拆分任务。
6. 用户使用 run-tasks / execute-task (slash command)执行任务
7. AI端到端测试+人工验收



新的目录结构，不完整，可调整：
```
project-root/
└── docs/
    ├── proposal/<slug>/           # Proposal
    │   └── proposal.md            # 提案
    ├── features/<slug>/           # Feature 工作区
    │   ├── prd/                   # 需求
    │   │   ├── overview.md        # 需求总览   
    │   │   ├── user-stories.md    # 用户故事
    │   │   └── ui-functions.md    # UI功能要点，非设计稿，只是关键要素的说明     
    │   ├── design/                # 设计
    │   │   ├── overview.md        # 设计总览
    │   │   ├── ui/
    │   │   └── api.md             # API文档      
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