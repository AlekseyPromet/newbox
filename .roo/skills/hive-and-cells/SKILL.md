---
name: hive-and-cells
description: Cilium’s Hive – a dependency injection framework for building modular Go applications.  
---

# Hive And Cells

## Instructions

Skill: Cilium Hive – Modular Dependency Injection in Go  

### Overview  

Cilium uses a dependency injection (DI) framework called **Hive** (package `pkg/hive`) to wire together the initialization, startup, and shutdown of its components. DI separates the *use* of objects from their *creation*: constructors only declare their dependencies as parameters, and the Hive library handles the rest. This encourages loosely coupled, testable, and inspectable architectures without global variables or hard‑to‑debug initialisation orders.  

**Core concepts**:  
- **Cell** – a building block that provides values (e.g. a server, a client) or performs side‑effects at startup/shutdown.  
- **Hive** – a container composed of cells; it can be run, inspected, and configured.  
- **Lifecycle** – a list of `Start` and `Stop` hooks executed in dependency order when the Hive runs.  

---

## 1. Creating a Hive  

A Hive is created with `hive.New()`, passing one or more cells:  

```go
var myHive = hive.New(foo.Cell, bar.Cell)
```

`New()` registers all providers but does **not** execute `Invoke` functions immediately – that happens when the Hive is started.  

To run the Hive (start → wait for signal → stop):  

```go
myHive.Run()
```

For tests, you can start/stop directly:  

```go
if err := myHive.Start(ctx); err != nil { /* ... */ }
if err := myHive.Stop(ctx); err != nil { /* ... */ }
```

---

## 2. Cells – the Building Blocks  

Cells are created using the `cell` package. The main cell constructors are:  

| Constructor      | Purpose                                                                                  |
|------------------|------------------------------------------------------------------------------------------|
| `cell.Module`    | A named set of cells, used for grouping related providers and configs.                   |
| `cell.Provide`   | Registers a constructor that returns one or more values (lazy – only invoked if needed). |
| `cell.ProvidePrivate` | Like `Provide`, but only visible inside the module and its sub‑modules.            |
| `cell.Invoke`    | Registers a function that is executed immediately (used to trigger instantiation).       |
| `cell.Config`    | Provides a configuration struct and automatically registers its flags.                   |
| `cell.Decorate`  | Wraps a set of cells, allowing you to augment the provided values.                       |
| `cell.Metric`    | Registers a Prometheus metrics collection struct.                                        |

### Example: A simple HTTP server cell  

```go
package server

import (
    "github.com/cilium/hive"
    "github.com/cilium/hive/cell"
    "github.com/spf13/pflag"
)

var Cell = cell.Module(
    "http-server",               // module identifier
    "HTTP Server",               // human‑readable title
    cell.Provide(New),           // provide the server constructor
    cell.Config(defaultServerConfig),
)

type Server interface {
    ListenAddress() string
    RegisterHandler(path string, fn http.HandlerFunc)
}

type ServerConfig struct {
    ServerPort uint16
}

var defaultServerConfig = ServerConfig{ServerPort: 8080}

func (def ServerConfig) Flags(flags *pflag.FlagSet) {
    flags.Uint16("server-port", def.ServerPort, "HTTP server listen port")
}

func New(lc cell.Lifecycle, cfg ServerConfig) Server {
    // ... initialise http.Server, register start/stop hooks ...
    return &serverImpl{...}
}
```

This cell can now be used in any Hive.  

---

## 3. Invoke – Triggering Execution  

`cell.Invoke` registers a function that is **always executed** when the Hive runs. Invoke functions can depend on any value provided by a `cell.Provide` constructor.  

### Example: Register a handler with the HTTP server  

```go
package hello

import "github.com/cilium/hive/cell"

var Cell = cell.Invoke(registerHelloHandler)

func helloHandler(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("hello"))
}

func registerHelloHandler(srv server.Server) {
    srv.RegisterHandler("/hello", helloHandler)
}
```

---

## 4. Lifecycle – Start / Stop Hooks  

The `cell.Lifecycle` type allows cells to schedule start and stop hooks. Hooks run in dependency order (start) and reverse order (stop).  

### Using the `Hook` struct  

```go
func New(lc cell.Lifecycle) *MyType {
    t := &MyType{}
    lc.Append(cell.Hook{
        OnStart: func(ctx cell.HookContext) error { return t.Start(ctx) },
        OnStop:  func(ctx cell.HookContext) error { return t.Stop(ctx) },
    })
    return t
}
```

Or implement the `HookInterface` directly.  

### Important guidelines  
- Constructors should **only do validation and allocation** – no goroutines or I/O.  
- Use `Start` hooks for launching goroutines and performing I/O.  
- `Stop` hooks **must block** until all resources are cleaned up (e.g. using `sync.WaitGroup`).  

---

## 5. Shutdowner – Graceful Termination  

If a component encounters a fatal error after startup, it can trigger a clean shutdown using the `hive.Shutdowner` interface:  

```go
type Example struct {
    Shutdowner hive.Shutdowner
}

func (e *Example) eventLoop() {
    // ...
    e.Shutdowner.Shutdown(hive.ShutdownWithError(err))
}
```

The Hive will then run all registered stop hooks in reverse order.  

---

## 6. Configuration  

Cells can provide a configuration struct using `cell.Config`. The struct must implement a `Flags(*pflag.FlagSet)` method. Hive automatically maps flags to struct fields by convention (e.g. `--server-port` → `ServerPort`).  

Configuration can be overridden in tests with `AddConfigOverride`:  

```go
h = hive.New(Cell)
h.AddConfigOverride(func(cfg *MyConfig) {
    cfg.MyOption = "test-override"
})
```

---

## 7. Metrics  

Use `cell.Metric` to define a metrics collection struct. All exported fields must implement `prometheus.Collector` (use types from `pkg/metrics/metric`). The struct is made available for injection, and all its collectors are automatically registered with the metrics registry.  

---

## 8. Inspecting a Hive  

After registering the Hive with a Cobra command, you can inspect it at runtime:  

```go
rootCmd.AddCommand(myHive.Command())
```

Then run:  

```bash
cilium$ go run ./daemon hive
```

This prints all modules, providers, configurations, and their dependencies. You can also generate a Graphviz dot‑graph:  

```bash
cilium$ go run ./daemon hive dot-graph | dot -Tx11
```

---

## 9. Testing with Hive Script  

Hive includes a scripting engine (`github.com/cilium/hive/script`) for writing integration tests.  

### Defining script commands in a cell  

```go
func ExampleCommands(e *Example) hive.ScriptCmdsOut {
    return hive.NewScriptCmds(map[string]script.Cmd{
        "example/hello": script.Command(
            script.CmdUsage{
                Summary: "Say hello",
                Args:    "name",
                Flags:   func(fs *pflag.FlagSet) { fs.String("greeting", "Hello,", "Greeting to use") },
            },
            func(s *script.State, args ...string) (script.WaitFunc, error) {
                // implementation
            },
        ),
    })
}
```

### Writing a test script  

Scripts are written as **txtar** files (extension `.txtar`) containing a sequence of commands:  

```
#! --enable-example=true
hive/start
example/hello foo
example/counts
```

Then run the test using `scripttest.Test()` – the Hive will start, execute the commands, and verify outputs.  

---

## 10. Best Practices & Guidelines  

| Rule                                                          | Why                                                                                  |
|---------------------------------------------------------------|--------------------------------------------------------------------------------------|
| Constructors do validation + allocation only                 | Side‑effectful constructors break inspection commands like `hive.PrintObjects`. |
| Use `Start` hooks for goroutines and I/O                      | Ensures predictable ordering and testability.                                       |
| `Stop` hooks must block until cleanup is complete             | Prevents test pollution and resource leaks.                            |
| Prefer interfaces over struct pointers in `Provide`           | Easier mocking, documentation, and dependency inspection.              |
| Use parameter (`cell.In`) and result (`cell.Out`) structs when constructors have many parameters | Improves readability and maintainability. |
| Utility cells should **not** contain `Invoke`                 | Keep cells lazy to avoid unwanted instantiation when the cell is not used. |
| Write a test for each non‑trivial cell                        | Tests serve as usage examples and validate dependencies.               |

---

## 11. Common Mistakes  

| Mistake                                       | Consequence                                                      |
|-----------------------------------------------|------------------------------------------------------------------|
| Starting goroutines in a constructor          | Breaks inspection and test predictability.                       |
| Forgetting to register `Config` flags         | Configuration options won’t be visible on the command line.      |
| Using `Provide` for side‑effectful functions   | May cause unexpected behaviour when the function is not needed.  |
| Not blocking in `Stop` hooks                  | Subsequent tests may interfere or resources may be leaked.       |
| Overusing `Invoke`                            | Forces eager instantiation, losing the benefits of lazy DI.      |

---

## 12. Further Reading  

- [Hive package documentation](https://pkg.go.dev/github.com/cilium/hive)  
- [Hive examples](https://github.com/cilium/hive/tree/main/example)  
- [Hive scripting engine](https://pkg.go.dev/github.com/cilium/hive/script)  

---
