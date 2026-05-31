---
id: "1"
title: "编写环境配置文档 environment-setup.md"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: 编写环境配置文档 environment-setup.md

## Description
创建 `docs/user-guide/environment-setup.md`，面向新用户提供完整的环境配置指南。当前 README.md 安装部分仅 25 行，缺少 OS 要求、Claude Code 版本兼容性、Go 环境验证和安装后验证步骤。本任务需要从代码库中提取实际的前置条件和安装流程，形成独立、完整的环境配置文档。

## Reference Files
- `README.md`: 第 46-71 行安装部分，覆盖 marketplace 和 local 两种安装方式 (source: proposal.md#Evidence)
- `.forge/config.yaml`: 前置条件和环境需求参考 (source: proposal.md#Proposed-Solution)
- `plugins/forge/commands/init-forge.md`: forge init 命令的安装验证流程 (source: proposal.md#Proposed-Solution)
- `docs/conventions/forge-distribution.md`: Forge 分发模型，理解安装路径机制 (source: proposal.md#Constraints-&-Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/user-guide/environment-setup.md` | 环境配置指南：前置条件、安装方式、环境验证、常见问题 |

### Modify
| File | Changes |
|------|---------|
| (无) | |

### Delete
| File | Reason |
|------|--------|
| (无) | |

## Acceptance Criteria
- [ ] 文档覆盖 3 种安装方式：Marketplace 安装、本地构建安装、开发模式安装
- [ ] 包含完整前置条件清单（操作系统、Go 版本、Claude Code CLI 版本及验证命令）
- [ ] 包含安装后验证步骤（`forge --version`、环境检查命令）
- [ ] 包含至少 3 条常见安装问题及解决方案（如 Go 版本不兼容、Claude Code 未安装、权限问题）
- [ ] 所有代码示例可直接复制执行，无需额外修改

## Implementation Notes
- 文档使用中文
- 每个安装方式提供从零开始的完整步骤，不假设用户已具备任何前置知识
- 从 README.md 和代码库中提取实际的命令和版本要求，不要凭记忆编写
- 文档顶部标注"最后更新"日期和对应版本（v3.0.0）
