# Example: Debugging Race Condition

## Session Input

```
> Error: `go test` failed with race condition in `internal/app/state.go`
> Cause: Two goroutines accessing map without synchronization
> Fix: Added `sync.RWMutex` to protect map access
```

## Generated Lesson

```markdown
# Protect Map Access with sync.RWMutex

## Problem
Race condition detected when multiple goroutines access a shared map.

## Root Cause
Go maps are not goroutine-safe. Concurrent read/write causes data race.

## Solution
Wrap map access with `sync.RWMutex`:
- `RLock()/RUnlock()` for reads
- `Lock()/Unlock()` for writes

## Reusable Pattern
When sharing state across goroutines, always protect with mutex.
Prefer `sync.RWMutex` when reads outnumber writes.

## Example
```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]string
}

func (s *SafeMap) Get(key string) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.m[key]
}

func (s *SafeMap) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[key] = value
}
```

## Related Files
- `internal/app/state.go`

## References
- https://go.dev/doc/articles/race_detector
- https://pkg.go.dev/sync#RWMutex
```
