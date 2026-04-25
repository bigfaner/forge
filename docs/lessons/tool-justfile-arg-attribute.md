# just `[arg]` Attribute: 官方文档未覆盖的 Recipe 选项参数

## Problem

justfile 中使用 `[arg("feature", long)]` 为 recipe 参数生成长选项，意图实现 `just test-e2e --feature` 的 flag 式调用，但实际执行报错：

```
error: Recipe `test-e2e` option `--feature` missing value
```

查阅官方 just book（stable release）的 Attributes 章节，完全找不到 `[arg(...)]` 属性，容易误判为无效语法。

## Root Cause

两层误解叠加：

1. **文档滞后于实现**：just 1.50.0 已支持 `[arg]` 属性（源码 `src/arg_attribute.rs`），但官方 book 尚未收录。只能通过阅读 casey/just 源码仓库确认。

2. **`[arg(long)]` 产生的是命名选项，不是布尔 flag**：
   - `[arg(long)]` 将参数变为 `--param <value>` 形式，必须提供值
   - 正确调用：`just test-e2e --feature yes`（任意非空值触发 feature 分支）
   - 错误调用：`just test-e2e --feature`（缺少值）

## Solution

**推荐方案：`[arg(long)]` + feature slug 作为参数值，不在脚本内嵌 `task feature`**

```just
[arg(long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        scripts_dir="docs/features/{{feature}}/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        # regression 模式
    fi
```

调用：
- `just test-e2e` — 回归测试
- `just test-e2e --feature <slug>` — 指定 feature 的测试

**好处**：
- 调用方显式传入 slug，不依赖 `task feature` 的运行时状态
- 脚本可独立运行，无需 task-cli 环境
- 外层编排器（CI、task all-completed）负责解析 feature slug 并传参，职责清晰

## Key Takeaway

1. **just 1.50.0 起支持 `[arg]` 属性**，但官方 book 文档未更新。遇到文档缺失时，查源码仓库 `src/arg_attribute.rs` 和 `src/parameter.rs` 确认。

2. **`[arg(long)]` ≠ 布尔 flag**。它生成 `--name <value>` 命名选项，调用时必须带值。如需无值 flag 效果，只能改用位置参数 + 空字符串默认值 + body 内条件判断。

3. **验证 justfile 的可靠方法**：`just --dry-run <recipe>` 可在不执行命令的情况下检查语法和参数替换结果。比仅查文档更可靠。

4. **just 版本差异大**：winget 安装的 1.50.0 远新于部分 Linux 发行版仓库中的版本。生产环境应锁定版本，避免 `[arg]` 等新特性在旧版上报错。

## Reference

- just 源码属性解析：`casey/just` → `src/arg_attribute.rs`
- just 官方 book（不含 `[arg]`）：https://just.systems/man/en/chapter_34.html
- just 源码 zread 文档（含 `[arg]` 说明）：`casey/just` → "配方参数与参数" 章节
