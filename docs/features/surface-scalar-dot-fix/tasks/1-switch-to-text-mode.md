---
id: "1"
title: "统一所有 skill 使用 forge surfaces 文本模式"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "high"
mainSession: false
breaking: false
---

# 1: 统一所有 skill 使用 forge surfaces 文本模式

## Description

当用户在 `.forge/config.yaml` 中使用 scalar 简写形式配置 surface（如 `surfaces: tui`），Go 解析器将其转换为哨兵 key `"."`。`forge surfaces --json` 将此哨兵直接输出为 `"key": "."`，导致依赖 JSON 输出解析 key 的 skill 拒绝该值。

解决方案：统一所有 skill 从 `forge surfaces --json` 切换到 `forge surfaces`（文本模式）。文本模式天然区分 scalar 和 named 两种形式，CLI 内部已吸收 `"."` 哨兵。

解析规则：每行按 `=` 分割，无 `=` 则为 scalar（只有 type，无 key）。Named key 形式（如 `app=tui`）下，`=` 左侧为 key，右侧为 type。

## Reference Files
- `docs/proposals/surface-scalar-dot-fix/proposal.md` — Proposed Solution, Scope > In Scope, Success Criteria
- `plugins/forge/skills/init-justfile/SKILL.md`: forge surfaces --json 调用和 `[a-zA-Z0-9_-]+` key 校验需替换为文本模式解析 (ref: Proposed Solution)
- `plugins/forge/skills/run-tests/SKILL.md`: recipe-prefix JSON 解析需替换为文本模式，scalar 无 prefix (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/SKILL.md`: 数据源从直接读取 config.yaml 切换到 `forge surfaces` 文本模式 (ref: Proposed Solution)
- `plugins/forge/skills/breakdown-tasks/SKILL.md`: Surface-Key/Type Inference 两层 JSON 策略需替换为文本模式 (ref: Proposed Solution)
- `plugins/forge/skills/quick-tasks/SKILL.md`: Surface-Key/Type Inference 两层 JSON 策略需替换为文本模式 (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | 数据源从 `--json` 切换到文本模式；移除 `[a-zA-Z0-9_-]+` key 校验；scalar 形式无 prefix recipe |
| `plugins/forge/skills/run-tests/SKILL.md` | 数据源从 `--json` 切换到文本模式；recipe-prefix 推导改为文本解析；scalar 无 prefix |
| `plugins/forge/skills/test-guide/SKILL.md` | 数据源从 config.yaml 切换到 `forge surfaces` 文本模式；统一解析规则 |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Surface-Key/Type Inference 从 JSON 两层策略切换到文本模式解析 |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Surface-Key/Type Inference 从 JSON 两层策略切换到文本模式解析 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] init-justfile 从 `forge surfaces --json` 切换到 `forge surfaces`（文本模式），scalar 形式（文本输出无 `=`）生成无 prefix recipe：`test`、`build`、`dev`、`teardown`
- [ ] run-tests 从 `forge surfaces --json` 切换到 `forge surfaces`（文本模式），scalar 形式调用 `just test` 而非 `just <key>-test`
- [ ] test-guide 从直接读取 config.yaml 切换到 `forge surfaces`（文本模式），统一数据源与其他 skill 一致
- [ ] breakdown-tasks 和 quick-tasks 的 Surface-Key/Type Inference 从 `--json` 切换到 `forge surfaces`（文本模式），scalar 形式下 surface-key 留空、surface-type 为 type 值
- [ ] Named key 形式（如 `app=tui`）下，recipe 名为 `<key>-<verb>`（如 `app-test`），行为与当前一致
- [ ] 所有 skill 使用统一解析规则：每行按 `=` 分割，无 `=` 则为 scalar（只有 type，无 key）；有 `=` 则左侧为 key，右侧为 type

## Hard Rules
- 仅修改 `plugins/forge/skills/` 下的 5 个 SKILL.md 文件，不修改 Go CLI 代码
- 不修改 `forge surfaces` CLI 的输出格式或行为

## Implementation Notes

### 解析规则定义（所有 skill 统一）

```
forge surfaces 文本输出解析：
  对每一行：
    if line contains '=':
      key = part before '='
      type = part after '='
      → named surface, recipe prefix = "<key>-"
    else:
      key = (empty)
      type = line
      → scalar surface, recipe prefix = ""
```

### Key Risks 缓解

- **文本解析逻辑不一致**：在所有 skill 中写明相同的解析规则段落（见上方定义），确保一致性
- **forge surfaces 输出格式变更**：`key=type` 格式已是 CLI 稳定接口，文本模式在 CLI 层已吸收 `"."` 哨兵

### 各 skill 变更要点

1. **init-justfile**：移除 `forge surfaces --json` 调用和 `jq` 解析，改为 `forge surfaces` 文本输出逐行解析。移除 `[a-zA-Z0-9_-]+` key 校验中对 `"."` 的防御。Scalar 形式下 recipe 名直接使用动词（无 prefix）
2. **run-tests**：移除 JSON 解析逻辑（array length 判断、key 字段提取），改为文本模式解析。Scalar 时 `recipe_prefix=""`，直接调用 `just <verb>`
3. **test-guide**：移除 Step 1a 中直接读取 config.yaml 的逗号分隔解析，改为调用 `forge surfaces` 并按统一规则解析
4. **breakdown-tasks** 和 **quick-tasks**：移除两层 JSON 策略（project-level + per-file），改为单次 `forge surfaces` 文本输出解析。Project-level 判断改为文本行数（单行=单 surface，多行=multi-surface）
