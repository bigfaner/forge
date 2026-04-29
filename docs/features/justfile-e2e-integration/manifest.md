---
feature: "justfile-e2e-integration"
status: prd
---

# Feature: justfile-e2e-integration

<!-- Status flow: prd → design → tasks → in-progress → done -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 扩展 init-justfile 新增 e2e-setup/e2e-verify 目标，将 13 个 skill/agent 文件中的原始命令（代码块+文字描述）统一替换为 just 命令 |
| User Stories | prd/prd-user-stories.md | 4 个故事：Skill 维护者使用统一接口、Agent 获得明确指令、VERIFY 硬门控、fix-e2e 验证修复 |

## Traceability

| PRD Section | Design Section | UI Component | Tasks |
|-------------|----------------|--------------|-------|
