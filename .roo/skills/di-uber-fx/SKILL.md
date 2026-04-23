---
name: di-uber-fx
description: A dependency injection framework for Go using constructors and lifecycle
---

# Di Uber Fx

## Instructions

I'm unable to directly access external URLs or read live articles. However, based on my knowledge of **Go dependency injection with Uber FX and Echo**, I can compile a comprehensive skill for an AI agent. If you'd like me to tailor it to the specific article you referenced, please provide its key points or paste the content.

Below is a **skill definition** an AI agent can use to answer questions or generate code for integrating Uber FX (DI container) with the Echo web framework in Go.

---

## Skill: Go Dependency Injection with Uber FX + Echo

### Purpose
Enable AI to assist developers in building modular, testable Go applications using Uber FX for dependency injection and Echo for HTTP routing.

### Core Concepts

| Concept | Description |
|---------|-------------|
| **FX** | A dependency injection framework for Go using constructors and lifecycle management. |
| **Echo** | High-performance, minimalist web framework. |
| **Provide** | FX function to add constructors to the container. |
| **Invoke** | FX function to trigger execution of functions that depend on provided types (e.g., start server). |
| **Lifecycle (fx.Lifecycle)** | Hooks (`OnStart`, `OnStop`) to manage app initialization and graceful shutdown. |

### Typical Integration Flow

1. **Define your dependencies** as interfaces or structs.
2. **Write constructors** for each component (e.g., handler, service, repository).
3. **Register Echo routes** inside a constructor that receives the Echo instance.
4. **Use `fx.Provide`** for each constructor.
5. **Use `fx.Invoke`** to start the Echo server with lifecycle hooks.

### Code Template (Minimal Working Example)

```go
package main

import (
    "context"
    "net/http"
    "time"

    "go.uber.org/fx"
    "github.com/labstack/echo/v4"
)

// --- Dependencies ---
type Handler struct {
    echo *echo.Echo
}

func NewHandler(e *echo.Echo) *Handler {
    h := &Handler{echo: e}
    h.registerRoutes()
    return h
}

func (h *Handler) registerRoutes() {
    h.echo.GET("/health", h.healthCheck)
}

func (h *Handler) healthCheck(c echo.Context) error {
    return c.String(http.StatusOK, "OK")
}

// --- Server lifecycle ---
func newEcho() *echo.Echo {
    e := echo.New()
    e.Server.ReadTimeout = 5 * time.Second
    e.Server.WriteTimeout = 10 * time.Second
    return e
}

func startEcho(lc fx.Lifecycle, e *echo.Echo) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go func() {
                if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
                    panic(err)
                }
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return e.Shutdown(ctx)
        },
    })
}

// --- FX app ---
func main() {
    app := fx.New(
        fx.Provide(newEcho),
        fx.Provide(NewHandler),
        fx.Invoke(startEcho),
    )
    app.Run()
}
```

### Best Practices for AI Assistance

- **Use interfaces** for major dependencies to facilitate testing (e.g., `UserService` interface with a real and mock implementation).
- **Group related constructors** into `fx.Options` for modularity.
- **Avoid global state** – let FX manage all dependencies.
- **Leverage `fx.Annotate`** when you need to name or group multiple implementations of the same type.
- **Use `fx.Populate`** to fill existing struct fields (less common, but useful for legacy code).

### Common Errors & Solutions

| Error | Likely Cause | Fix |
|-------|--------------|-----|
| `missing type in container` | Constructor not provided or returns wrong type | Add `fx.Provide` for that constructor |
| `cyclic dependency` | A depends on B, B depends on A | Introduce interface or restructure |
| `fx.Invoke called with function that has unprovided params` | Missing `fx.Provide` for one of the parameters | Provide that dependency |
| Echo server doesn't start | `startEcho` not invoked, or `OnStart` returns before server actually starts | Ensure `go func()` is used and `nil` is returned immediately |

### When to Suggest FX + Echo

- Project is medium-to-large and will benefit from explicit dependency management.
- Team wants to avoid global variables and manual wiring.
- Need lifecycle hooks (e.g., DB connections, graceful shutdown) out of the box.

### Alternative DI approaches in Go (for context)

| Approach | When to use |
|----------|-------------|
| **Manual wiring** (pass dependencies explicitly) | Small projects, or when maximum simplicity is desired. |
| **Wire** (Google) | Compile‑time DI with no runtime reflection. |
| **FX** (Uber) | Runtime DI with powerful lifecycle and modularity features. |

---

## How an AI Agent Can Use This Skill

When a user asks:
- *“How do I set up FX with Echo?”* → Provide the code template and explain the lifecycle.
- *“How to inject a database connection?”* → Show adding a `*sql.DB` constructor and providing it to handlers.
- *“What’s the difference between Provide and Invoke?”* → Explain that Provide adds constructors, Invoke triggers execution.
- *“My Echo server doesn’t shut down gracefully with FX”* → Show the `OnStop` hook calling `e.Shutdown(ctx)`.
