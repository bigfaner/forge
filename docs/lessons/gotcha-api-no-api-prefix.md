# API 路由无 /api 前缀

## Problem

调用后端接口时使用了 `/api/v1/...` 路径，返回 404。

## Root Cause

本项目的 Gin 路由定义为：

```go
v1 := r.Group(basePath + "/v1")
```

`basePath` 默认为空字符串，因此所有 API 路由的实际路径是 `/v1/...`，**没有 `/api` 前缀**。

Vite 代理配置也印证了这一点：

```ts
proxy: {
  '/v1': 'http://localhost:8080',
}
```

## Solution

使用正确的路径前缀：

| 错误 | 正确 |
|------|------|
| `POST /api/v1/auth/login` | `POST /v1/auth/login` |
| `GET /api/v1/teams/1/main-items` | `GET /v1/teams/1/main-items` |

## Key Takeaway

不要假设后端 API 带有 `/api` 前缀。在调用接口前，先查看 `backend/internal/handler/router.go` 中的 `r.Group(...)` 定义，或查看 `frontend/vite.config.ts` 中的 proxy 配置来确认实际路径前缀。
