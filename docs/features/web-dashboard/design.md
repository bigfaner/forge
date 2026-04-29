# Web Dashboard Design

## 设计方向
使用 Tailwind CSS CSS Variable 方案，支持 Light/Dark 双主题，以 OLED 深色为主打。

## 色彩体系
- Background: `#0F172A` (dark) / `#F8FAFC` (light)
- Accent: `#22C55E` (运行绿)
- Destructive: `#EF4444`

## 布局
- 左侧固定侧边栏 (w-56)
- 右侧 main 滚动区
- 内容最大宽度 1200px

## 组件规范
- 卡片：rounded-lg border border-border bg-card
- 徽章：rounded border px-1.5 py-0.5 text-xs
- 按钮：bg-accent text-white rounded-md
