# Lesson: TUI Task 的 Verify Criteria 必须包含视觉验证条件

## 问题

deep-drill-analytics feature 的 Forge task 完成后，agent 认为任务已完成（编译通过、测试通过），但运行时出现 11 个 bug 和 5 个样式问题。根因是 task 的 verify criteria 只检查了功能性，没有覆盖视觉渲染正确性。

### 具体表现

以下 bug 全部通过了"编译 + 测试"的 verify gate，但在运行时暴露：

| Bug | 根因 | 为什么测试没发现 |
|-----|------|----------------|
| 屏幕乱码 | 内容溢出触发终端换行 | 功能测试不检查渲染输出尺寸 |
| 列不对齐 | `len()` 替代 `runewidth` | 值正确但视觉对齐错误 |
| 百分比换行 | 宽度未考虑 scrollbar | scrollbar 是运行时动态的 |
| 路径截断丢失上下文 | 按字符而非按段截断 | 截断后路径仍唯一可读 |
| Rx/Ex 后缀错位 | 未右 pad 路径到统一宽度 | 各行数据正确但列不对齐 |

**Why:** TUI 渲染是"视觉正确性"，不是"逻辑正确性"。Go 的 `_test.go` 可以验证函数返回值，但无法验证 lipgloss 渲染后的终端输出是否对齐、是否溢出、是否在窄终端下正常。这需要不同类型的验证条件。

## 规则

对于涉及 TUI 渲染（`internal/model/*.go` 的 View()/Render() 函数）的 task，verify criteria **必须**包含以下三类检查：

### 1. Golden File 对比（自动）

所有新增或修改的 View()/Render() 函数必须有对应的 golden test：

```go
func TestDashboard_View(t *testing.T) {
    m := newDashboardModelWithTestData()
    got := m.View()

    // 维度检查：行数 == 终端高度，每行宽度 <= 终端宽度
    lines := strings.Split(got, "\n")
    if len(lines) != m.height {
        t.Errorf("lines = %d, want %d", len(lines), m.height)
    }
    for i, line := range lines {
        if lipgloss.Width(line) > m.width {
            t.Errorf("line %d width %d > terminal %d", i, lipgloss.Width(line), m.width)
        }
    }

    // Golden 对比
    goldenFile := filepath.Join("testdata", "dashboard_populated.golden")
    // ...
}
```

**检查项**：
- View() 输出行数 == 终端高度
- 每行 lipgloss.Width <= 终端宽度
- 输出与 golden file 一致（忽略 ANSI 色码差异）

### 2. 边界条件场景（自动）

在 verify criteria 中列出必须测试的场景：

```markdown
## Verify
- [ ] `go test ./internal/model/... -run TestDashboard` passes
- [ ] Golden: 行数=m.height, 每行宽度<=m.width
- [ ] Narrow terminal (80x24): 布局不溢出
- [ ] Wide terminal (140x40): 布局不变形
- [ ] Mixed digit widths (1 vs 100): 列对齐正确
- [ ] Long path (>50 chars): 截断后保留尾部段
- [ ] Scrollbar present: 内容宽度自动减 1
- [ ] Empty data: 显示 "No data" 不报错
```

### 3. 测试数据真实性

Golden test 的测试数据**必须**包含真实数据中会出现的边界值。不允许只用简单 ASCII 短字符串：

```go
// Bad: 测试数据全是短 ASCII
files: []FileEntry{{Path: "a.go", ReadCount: 1}}

// Good: 包含边界值
files: []FileEntry{
    {Path: "a.go", ReadCount: 1},                                          // 短路径 + 单数字
    {Path: "internal/pkg/handler/middleware/auth_handler.go", ReadCount: 100}, // 长路径 + 三位数
    {Path: "组件/配置文件.go", ReadCount: 5},                                // CJK 路径
}
```

**强制最低要求**（任一缺失则 verify 不通过）：

| 边界值 | 示例 | 防止的 bug |
|--------|------|-----------|
| CJK 字符串 | `"配置文件"` | runewidth vs len 对齐错误 |
| 长路径 (>50 chars) | `"internal/pkg/handler/middleware/auth.go"` | 截断逻辑错误 |
| 多位数字 (>9) | `ReadCount: 100` | 列对齐偏移 |
| 空字段 | `EditCount: 0` | 缺失字段布局错误 |

### 3. 视觉证据（手动 / agent 自检）

对于没有 golden test 覆盖的新面板，task executor 应在 verify 阶段输出渲染结果：

```bash
# 运行 golden test 并打印实际输出（--update 不适用时）
go test ./internal/model/... -run TestDashboard -v 2>&1 | head -50
```

如果 task executor 能运行程序，截取终端输出并对比 tech design 中的 ASCII mockup。

### Verify Criteria 模板

在 `/breakdown-tasks` 阶段，当 task 涉及 `internal/model/` 的 View/Render 函数时，自动附加以下 verify 模板：

```markdown
## Verify

### Functional
- [ ] `just compile` passes
- [ ] `just test` passes

### TUI Rendering
- [ ] Golden test exists for new/modified View() function
- [ ] Dimension check: output lines == terminal height, each line width <= terminal width
- [ ] Boundary scenarios covered (see list in task description)
- [ ] No hardcoded widths — all derived from m.width
- [ ] Colors from palette only (docs/conventions/tui-layout-ui.md)
```

## How to apply

### 当前（手动，skill 修改前）

1. **撰写 task 时**：在 task description 的 verify criteria 中，手动追加 TUI Rendering 部分和测试数据真实性要求
2. **执行 task 时**：agent 在 verify 阶段运行 golden test 后，额外检查输出维度（行数、宽度），而非仅依赖 `just test` exit code
3. **Review task 时**：检查 golden test 测试数据是否包含 4 个强制边界值

### 后续（Forge skill 改进后）

1. `/breakdown-tasks` skill 识别 task scope 含 `internal/model` 时，自动追加 TUI verify 模板
2. `/execute-task` 的 quality gate 在 `just test` 之后增加 golden test 维度检查
3. `/eval-design` 对 TUI feature 检查 tech design 是否包含 boundary scenarios

## 预期效果

| 指标 | 改进前 | 改进后（预期） |
|------|--------|---------------|
| 后功能 fix 提交 | 11 | 2-3 |
| Bug 在 task verify 阶段被捕获 | 0/11 | 8-9/11 |
| Golden test 覆盖率 | 事后补 | task 级别即有 |

## Tags

`verify-criteria`, `tui`, `golden-test`, `process`, `forge-improvement`
