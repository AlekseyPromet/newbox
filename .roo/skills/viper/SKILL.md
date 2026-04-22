---
name: viper
description: This skill provides a comprehensive guide to using **Viper** – a complete configuration solution for Go applications. Viper supports multiple configuration formats (JSON, TOML, YAML, HCL, envfile, Java properties), environment variables, command-line flags, remote systems (etcd, Consul), and live reloading.
---

# Viper

## Instructions

---

## 1. Installation

```bash
go get github.com/spf13/viper
```

---

## 2. Core Concepts

Viper follows a **priority order** for configuration values (highest to lowest):

1. **`Set()`** – explicit override
2. **Command-line flags** (via `pflag`)
3. **Environment variables**
4. **Configuration file**
5. **Defaults** (lowest)

This allows you to build 12‑factor apps where deployment settings override developer defaults.

---

## 3. Basic Usage Pattern

```go
import "github.com/spf13/viper"

func initConfig() error {
    // Set defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.host", "localhost")

    // Read config file
    viper.SetConfigName("config")   // name of config file (without extension)
    viper.SetConfigType("yaml")     // optional – inferred from extension if present
    viper.AddConfigPath(".")        // look for config in current directory
    viper.AddConfigPath("/etc/app/") // look in /etc/app

    // Auto‑read environment variables
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")       // e.g. APP_SERVER_PORT overrides server.port

    // Bind to specific environment variable for a key
    viper.BindEnv("server.port", "APP_SERVER_PORT")

    // Read the file
    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("failed to read config: %w", err)
    }
    return nil
}
```

---

## 4. Configuration Sources

### 4.1 Configuration Files

Supported formats: `json`, `toml`, `yaml`, `yml`, `properties`, `props`, `prop`, `hcl`, `env`, `dotenv`

**File example (config.yaml):**
```yaml
server:
  port: 8080
  timeout: 30s
database:
  host: "localhost"
  port: 5432
  name: "mydb"
```

**Reading values:**
```go
port := viper.GetInt("server.port")
timeout := viper.GetDuration("server.timeout")
dbHost := viper.GetString("database.host")
```

Viper supports **nested keys** using a delimiter (default `.`).

### 4.2 Environment Variables

```go
viper.AutomaticEnv()           // automatically match env vars to keys
viper.SetEnvPrefix("myapp")    // env var prefix: MYAPP_SERVER_PORT
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // convert dots to underscores
```

**Mapping rules:**
- `viper.Get("server.port")` looks for `MYAPP_SERVER_PORT` (if prefix set)
- If `SetEnvKeyReplacer` is used, `server.port` becomes `SERVER_PORT` (no prefix)
- Case‑insensitive by default

### 4.3 Command‑Line Flags (with pflag)

```go
import flag "github.com/spf13/pflag"

func init() {
    flag.Int("server.port", 8080, "Server port")
    flag.Parse()
    viper.BindPFlags(flag.CommandLine)
}
```

Flags have the **highest priority** after `viper.Set()`.

### 4.4 Remote Configuration (Consul, etcd)

```go
import _ "github.com/spf13/viper/remote"

viper.AddRemoteProvider("consul", "localhost:8500", "config/myapp")
viper.SetConfigType("yaml")
err := viper.ReadRemoteConfig()
```

### 4.5 Direct `Set()` Override

```go
viper.Set("server.port", 9000) // highest priority
```

---

## 5. Reading Values

| Method                      | Returns                     |
|-----------------------------|-----------------------------|
| `Get(key string) interface{}` | raw value                    |
| `GetString(key)`            | string                       |
| `GetInt(key)`               | int                          |
| `GetInt64(key)`             | int64                        |
| `GetBool(key)`              | bool                         |
| `GetFloat64(key)`           | float64                      |
| `GetDuration(key)`          | time.Duration (parses `30s`, `5m`, etc.) |
| `GetStringSlice(key)`       | []string                     |
| `GetStringMap(key)`         | map[string]interface{}       |
| `GetStringMapString(key)`   | map[string]string            |
| `GetTime(key)`              | time.Time                    |
| `IsSet(key)`                | bool – whether key exists    |
| `AllSettings()`             | map[string]interface{} of all settings |

**Example:**
```go
type ServerConfig struct {
    Port    int           `mapstructure:"port"`
    Timeout time.Duration `mapstructure:"timeout"`
}

var config ServerConfig
err := viper.UnmarshalKey("server", &config)
// or unmarshal all keys: viper.Unmarshal(&config)
```

**Note:** Use `mapstructure` tags to match Viper’s nested keys.

---

## 6. Default Values

```go
viper.SetDefault("cache.ttl", 300)           // seconds
viper.SetDefault("log.level", "info")
viper.SetDefault("features", []string{"metrics", "tracing"})
```

Defaults are **lowest priority** – overridden by file, env, flags, or `Set()`.

---

## 7. Watching Configuration Changes (Live Reload)

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
    // Reload internal structures
    var newConfig AppConfig
    viper.Unmarshal(&newConfig)
    // Apply newConfig atomically
})
```

Works with local files only (not remote). Supported on most OS except Windows (partial support).

---

## 8. Unmarshaling – Best Practice

Define your configuration struct with `mapstructure` tags:

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
    Port    int           `mapstructure:"port"`
    Timeout time.Duration `mapstructure:"timeout"`
}

var C Config
err := viper.Unmarshal(&C)
```

**Why use `mapstructure`?**  
Viper uses the `mapstructure` library to decode config maps into structs. Tags help with case‑insensitive matching and renaming.

If you prefer **default values** in the struct, use `mapstructure:",omitempty"` and set default values in code before unmarshaling.

---

## 9. Configuration File Locations

Viper searches in order of added paths:

```go
viper.AddConfigPath("/etc/app/")      // system config
viper.AddConfigPath("$HOME/.app")     // user config
viper.AddConfigPath(".")              // current directory
```

You can also set an explicit file:

```go
viper.SetConfigFile("/path/to/config.yaml")
```

---

## 10. Writing / Saving Configuration

Viper does **not** provide built‑in write methods. Use the `go-yaml` or `go-json` libraries to marshal and save.

---

## 11. Best Practices

| Practice                                  | Why                                                                 |
|-------------------------------------------|---------------------------------------------------------------------|
| Use **one configuration struct** and `Unmarshal` | Type safety, validation, IDE autocompletion.                     |
| Set defaults in code using `SetDefault`   | Avoids empty values when config file is missing.                   |
| Use **environment variables** for secrets | Never commit secrets to config files.                               |
| Prefer **nested keys** (`database.host`)  | Better organisation than flat keys.                                 |
| Call `viper.AutomaticEnv()` only once     | It scans env vars repeatedly – fine, but don’t call inside loops.   |
| Validate config after unmarshaling        | Use a `Validate()` method on your struct.                           |
| Use `viper.GetDuration()` for time values | Allows human‑readable `30s`, `5m` in config files.                  |

---

## 12. Common Mistakes & Pitfalls

| Mistake                                           | Consequence / Fix                                                      |
|---------------------------------------------------|------------------------------------------------------------------------|
| Forgetting to call `ReadInConfig()`               | Viper returns defaults only, no error.                                 |
| Mixing `SetConfigName` and `SetConfigFile`        | `SetConfigFile` overrides the name – use one.                          |
| Not handling `ReadInConfig` error                 | Silent failure: app runs with incomplete config.                       |
| Using `.` inside environment variable names       | Viper replaces `.` with `_` – use `SetEnvKeyReplacer`.                 |
| Assuming case‑sensitive matching with env vars   | Viper matches `APP_DB_HOST` to `db.host` case‑insensitively.           |
| Unmarshaling without `mapstructure` tags         | Works if field names match exactly, but brittle.                       |
| Watching config on Windows with default fsnotify | Some editors save as temp file → event not caught. Use polling?        |

---

## 13. Complete Example

**main.go**
```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    Server struct {
        Port    int           `mapstructure:"port"`
        Timeout time.Duration `mapstructure:"timeout"`
    } `mapstructure:"server"`
    Database struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"database"`
}

func main() {
    // Set defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.timeout", "30s")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", 5432)

    // Config file
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/myapp")

    // Environment vars
    viper.AutomaticEnv()
    viper.SetEnvPrefix("MYAPP")
    viper.BindEnv("server.port", "MYAPP_PORT")

    // Read config
    if err := viper.ReadInConfig(); err != nil {
        log.Printf("No config file found: %v", err)
    }

    // Unmarshal
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log.Fatalf("Failed to unmarshal: %v", err)
    }

    // Live reload
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Println("Config changed, reloading...")
        viper.Unmarshal(&cfg)
    })

    fmt.Printf("Running on :%d, timeout %v\n", cfg.Server.Port, cfg.Server.Timeout)
}
```

**config.yaml**
```yaml
server:
  port: 3000
database:
  host: "prod-db.example.com"
```

**Run with override:**  
```bash
MYAPP_PORT=9090 go run main.go
```

---

## 14. Further Reading

- [Viper GitHub](https://github.com/spf13/viper)
- [12‑Factor App – Config](https://12factor.net/config)
- [pflag package](https://github.com/spf13/pflag)

---
