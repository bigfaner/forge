# Sitemap: Architecture Information Must Guide Exploration Strategy

## Problem

Sitemap 生成了 351 个元素，其中 ~120 个是侧边栏重复（每页 ~10 个侧边栏元素 × 12 个认证页面）。用户指出后才发现。

## Root Cause

**不是"去重遗漏"，是拿到架构信息后没有用它。**

探索流程中，我读了 `App.tsx` 的路由定义，明确看到 `<Route element={<AppLayout />}>` 包裹所有认证路由。这已经告诉我：侧边栏是共享布局组件，不是页面内容。但我对这个信息视而不见，对每个页面机械执行 snapshot → 抄录所有元素，没有任何"这个元素属于布局还是页面"的判断。

本质错误：**把生成 sitemap 当成机械的逐页抄录任务，而不是需要架构理解的结构化提取。**

## Solution

**正确的生成顺序应该是：**

1. 先读路由定义，识别布局嵌套（哪些路由共享 `AppLayout`、`PermissionRoute`）
2. 只探索首页，提取共享布局元素（sidebar、header），归入 sitemap 的 `layout` 字段
3. 后续每个页面只记录路由特有的内容区元素，跳过布局部分

```
App.tsx 路由结构 → 识别 layout 层 → 提取共享元素一次 → 逐页只记录页面元素
```

## Key Takeaway

**在开始机械执行前，先问：我已经掌握的信息中，有没有能改变执行策略的？**

拿到 App.tsx 的路由结构后，应该先建立"共享 vs 页面特有"的心智模型，再开始探索。而不是先无脑执行再回头补去重。架构理解应该前置，不是后置。
