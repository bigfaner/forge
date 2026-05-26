---
created: "2026-05-26"
tags: [architecture, testing]
---

# Surface Resolution Shortcut for Task Generation

## Problem

quick-tasks skill 的 Surface-Key/Type Inference 步骤要求对每个 task 的受影响文件逐一调用 `forge surfaces --json <file-path>` 解析 surface。在单 surface 项目中，所有文件返回相同结果，导致 N 次（N = coding task 数）冗余外部命令调用。

## Root Cause

1. SKILL.md 模板按字面描述了 "per-task per-file query" 流程，未区分单 surface vs 多 surface 场景
2. Agent 机械执行模板步骤，未先判断项目是否需要逐文件解析
3. `forge surfaces` 查询的是项目级配置，不是文件级属性——同一项目所有文件必然返回相同 surface

## Solution

在执行 Surface-Key/Type Inference 步骤前，先判断项目 surface 拓扑：

- **单 surface 项目**（`surfaces: api` scalar 形式，或 map 只有一个 key）：直接使用该 surface，跳过所有文件级查询
- **多 surface 项目**：根据 file_path 前缀快速判断归属 surface（如 `backend/` → api, `frontend/` → web），仅对路径无法推断的文件调用 `forge surfaces --json`

## Reusable Pattern

Surface resolution 应分两层：
1. **项目级**：读 `.forge/config.yaml` 的 `surfaces` 字段，一次判断单/多 surface
2. **文件级**（仅多 surface 需要）：路径前缀匹配优先，`forge surfaces` 兜底

判断顺序：项目级 → 路径前缀 → 外部命令。每层能解决就跳过后续层。

## Example

```yaml
# 单 surface，无需逐文件查询
surfaces: cli

# 多 surface，按路径前缀推断
surfaces:
  backend: api    # forge-cli/ → backend
  frontend: web   # web/ → frontend
```

## References

- `skills/quick-tasks/SKILL.md` — Step 2 Surface-Key/Type Inference
