# Raycast UI PRD

## 概述
以 `web/` 为基础，fork 出 `web-raycast/` 副本，按照 `DESIGN.md` 中的 Raycast 设计规范完整重写前端视觉层。

## 设计规范来源
DESIGN.md（由 `npx getdesign@latest add raycast` 生成）

## 核心要求

### 视觉
- 背景必须为 `#07080a`（近黑蓝色调），禁用纯黑
- 主品牌色 Raycast Red `#FF6363`，仅用于点睛，不泛用
- 所有 body text 使用 weight 500 + letter-spacing +0.2px
- 卡片使用 double-ring shadow 系统

### 交互
- 按钮 hover 用 opacity 0.6，不改变背景色
- 输入框 focus 用蓝色 glow `hsla(202,100%,67%,0.15)`
- 导航链接 active 用 `rgba(255,255,255,0.06)` 背景

### 隔离
- 独立目录 `web-raycast/`，不修改原 `web/`
- Vite dev server 端口 7301（原版 7300）

## 验收标准
- [ ] TypeScript 0 error
- [ ] 视觉通过 DESIGN.md Section 7 Do's and Don'ts 全部 checklist
- [ ] 背景色用色值确认为 `#07080a`
- [ ] 响应式断点 600px / 768px / 1024px 测试通过
