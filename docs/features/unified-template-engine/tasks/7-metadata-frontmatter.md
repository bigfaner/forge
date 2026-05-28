---
id: "7"
title: "统一模板 metadata frontmatter"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [3, 4, 5]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 7: 统一模板 metadata frontmatter

## Description
为全部 41 个模板文件（21 prompt + 2 task creation + 12 autogen body + 6 record）添加统一的 metadata frontmatter（type, category, variables 字段）。Record 模板重构：输出 frontmatter（status, started 等）从 metadata 区移入 body 区，由 Go template 渲染。模板加载器在 `template.Parse()` 前剥离 metadata frontmatter。实现 `ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 反射字段交叉校验。autogen body 模板的 `<!-- body-only -->` 注释被 metadata frontmatter 替代。

## Reference Files
- `forge-cli/pkg/prompt/data/*.md`: 21 个 prompt 模板添加 metadata frontmatter (source: proposal.md#Frontmatter-规范统一)
- `forge-cli/pkg/task/data/record-*.md`: 6 个 record 模板重构为双 frontmatter 结构 (source: proposal.md#Frontmatter-规范统一)
- `forge-cli/pkg/task/data/*.md`: 12 个 autogen body 模板添加 metadata，移除 <!-- body-only --> (source: proposal.md#Frontmatter-规范统一)
- `forge-cli/pkg/prompt/prompt.go`: 模板加载器增加 metadata frontmatter 剥离逻辑 + ValidateTemplates() (source: proposal.md#Frontmatter-规范统一)
- `forge-cli/pkg/template/data/*.md`: 2 个任务创建模板添加 metadata frontmatter (source: proposal.md#Frontmatter-规范统一)

## Acceptance Criteria
- [ ] 全部 41 个模板文件包含 metadata frontmatter（type, category, variables 字段），格式符合提案规范
- [ ] Record 模板结构正确：metadata frontmatter + body（含 `---` 输出 frontmatter + 内容），渲染输出不含 metadata
- [ ] 模板加载器在 `template.Parse()` 前正确剥离 metadata frontmatter，渲染输出不含 metadata 内容
- [ ] `ValidateTemplates()` 使用 `variables` 字段与对应 struct 反射字段做正向交叉校验（每个声明 variable 必须在 struct 中有匹配字段）
- [ ] autogen body 模板中无 `<!-- body-only -->` 注释残留

## Implementation Notes
- metadata frontmatter 由模板加载器在 `template.Parse()` 前剥离，不参与渲染输出
- Record 模板文件头部需添加注释解释双 frontmatter 结构
- autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释——frontmatter 的 `variables` 字段本身表达"此模板期望哪些变量"
- `missingkey=error` 确保反向覆盖（struct 字段存在但模板未使用时不报错，但模板使用不存在的字段时报错）
