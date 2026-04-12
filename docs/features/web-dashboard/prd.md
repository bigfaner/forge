# Web Dashboard PRD

## 概述
为 ZCode 项目构建一套基于 React + Vite 的 Web 看板，提供任务管理、Feature 追踪、执行记录查询等功能。

## 目标用户
- 项目开发者：查看 Feature 进度与任务状态
- Claude Agent：认领任务、写入执行记录

## 核心功能

### Dashboard
- 全局任务统计（总数 / 完成 / 进行中 / 阻塞）
- Active Features 卡片网格，附进度条
- 点击跳转 Feature 详情

### Feature Detail
- Kanban / List / DAG 三种任务视图
- 内联 PRD / Design 文档查看
- Claim Task 一键认领

### Records & Lessons
- 时间线展示历史执行记录
- Lessons 知识库，支持分类过滤和搜索

## 技术栈
- React 18 + TypeScript
- Vite 6 + TanStack Query
- Tailwind CSS v3
- React Router v6
- @xyflow/react (DAG 视图)
