---
name: sqlc
description: Here is a skill guide for using sqlc with PostgreSQL, intended for an AI agent. It covers the complete workflow from installation and configuration to generating type-safe code and managing a real-world project.
---

# Sqlc

## Instructions

sqlc is a code generation tool that compiles SQL queries into type-safe, idiomatic Go code. It bridges the gap between raw SQL and ORMs by offering compile-time query validation, full SQL feature support, and better performance than traditional ORMs. sqlc itself has no runtime dependencies, but you will need the Go toolchain to build and run programs that use the generated code.

## 2. Installation

### 2.1 Install sqlc
Follow the official [installation guide](https://docs.sqlc.dev/en/latest/overview/install.html). For macOS and Linux, the quickest method is often using `go install`:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### 2.2 Verify Installation
```bash
sqlc version
```

## 3. Project Setup

Create a new directory for your project and initialize a Go module:
```bash
mkdir sqlc-tutorial
cd sqlc-tutorial
go mod init tutorial.sqlc.dev/app
```

## 4. Configuration (`sqlc.yaml`)

Create a `sqlc.yaml` file in the root of your project. The configuration file uses YAML format and consists of a version declaration, SQL configuration, and code generation settings.

A minimal PostgreSQL configuration looks like this:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
        sql_package: "pgx/v5"
```

### 4.1 Key Configuration Fields

| Field | Description |
|-------|-------------|
| `version` | Configuration version (use "2") |
| `engine` | Database engine: `postgresql`, `mysql`, or `sqlite` |
| `schema` | Directory of SQL migrations or path to single SQL file |
| `queries` | Directory of SQL queries or path to single SQL file |
| `gen.go.package` | Go package name for generated code |
| `gen.go.out` | Output directory for generated code |
| `gen.go.sql_package` | SQL driver package (`pgx/v4`, `pgx/v5`, or `database/sql`) |

### 4.2 Advanced Configuration Options

```yaml
version: "2"
sql:
  - schema: "db/migrations"
    queries: "db/queries"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true      # Add JSON tags to generated structs
        emit_prepared_queries: true
        emit_interface: true       # Generate Querier interface
        emit_exact_table_names: false
    database:
      uri: "postgresql://postgres:${PG_PASSWORD}@localhost:5432/app"  # Environment variable support
```

The `uri` string can contain references to environment variables using the `${...}` syntax.

## 5. Schema Definition

Create a `schema.sql` file with your database schema:
```sql
CREATE TABLE authors (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    bio text
);
```

### 5.1 Schema Modifications
sqlc parses `CREATE TABLE` and `ALTER TABLE` statements to generate the necessary code. For example:
```sql
ALTER TABLE authors ADD COLUMN birth_year int NOT NULL;
ALTER TABLE authors DROP COLUMN birth_year;
ALTER TABLE authors RENAME TO writers;
```

### 5.2 Multi-Schema Environments
When working with multiple schemas, always specify the schema name explicitly in your SQL queries:
```sql
-- Correct: Explicit schema prefix
INSERT INTO schema_name.table_name (id, name) VALUES ($1, $2);
```
This is particularly important when migrating from local development to cloud environments, where sqlc may behave differently.

## 6. Writing Queries

Create a `query.sql` file. Each query requires a special comment annotation that tells sqlc the query name and command type:

### 6.1 Query Annotations Format
```sql
-- name: <QueryName> <CommandType>
<SQL statement>
```

### 6.2 Command Types

| Command | Description | Generated Method Return |
|---------|-------------|------------------------|
| `:one` | Returns a single row | `(Model, error)` |
| `:many` | Returns multiple rows | `([]Model, error)` |
| `:exec` | Executes without returning rows | `error` |
| `:execrows` | Executes and returns number of affected rows | `(int64, error)` |
| `:batch` | Batches multiple queries together | Custom batch type |

### 6.3 Example Queries
```sql
-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (name, bio) VALUES ($1, $2) RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors SET name = $2, bio = $3 WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;
```

### 6.4 Returning Updated Records
To return the updated record, use the `:one` command with `RETURNING *`:
```sql
-- name: UpdateAuthor :one
UPDATE authors SET name = $2, bio = $3 WHERE id = $1 RETURNING *;
```

## 7. Code Generation
для запуска slqc generate ипользуй команду `docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate` где %cd% - заменить на реальный путь к файлу миграции
Run the generate command:
```bash
docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate
```

Upon success, sqlc creates the output directory containing:
- `db.go` - Database connection and Querier interface
- `models.go` - Type-safe model structs for your tables
- `query.sql.go` - Generated query methods

## 8. Using Generated Code

### 8.1 Basic Usage with pgx
```go
package main

import (
    "context"
    "log"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgtype"
    "tutorial.sqlc.dev/app/tutorial"
)

func run() error {
    ctx := context.Background()
    conn, err := pgx.Connect(ctx, "postgresql://user:pass@localhost:5432/dbname")
    if err != nil {
        return err
    }
    defer conn.Close(ctx)

    queries := tutorial.New(conn)

    // Create an author
    insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
        Name: "Brian Kernighan",
        Bio: pgtype.Text{
            String: "Co-author of The C Programming Language",
            Valid:  true,
        },
    })
    if err != nil {
        return err
    }

    // Get the author
    fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
    if err != nil {
        return err
    }

    log.Println(fetchedAuthor)
    return nil
}
```

### 8.2 Using with database/sql
```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

conn, err := sql.Open("postgres", "postgresql://...")
queries := db.New(conn)
```

## 9. SQL Migrations

**Important:** sqlc does not perform database migrations for you. It only reads schema files to generate code.

### 9.1 Migration Tools Integration
sqlc can parse migrations from popular tools by pointing the `schema` field to the migration directory instead of a single file:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "db/migrations"  # Directory containing migration files
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
```

### 9.2 Supported Migration Tools
- **Goose**: Uses `-- +goose Up` and `-- +goose Down` annotations
- **golang-migrate**: Uses `-- +migrate Up` and `-- +migrate Down`
- **tern**: Uses `---- create above / drop below ----` separator
- **dbmate**: Uses `-- migrate:up` and `-- migrate:down`

### 9.3 Important Note for golang-migrate
sqlc parses migration files in lexicographic order, not numeric order. To avoid unexpected behavior, prepend enough zeroes to your migration filenames:
- ✅ Good: `001_initial.up.sql`, `002_add_column.up.sql`
- ❌ Bad: `1_initial.up.sql`, `10_add_column.up.sql`

## 10. Error Handling Best Practices

PostgreSQL errors contain structured information via SQLSTATE codes. Use type assertion to extract detailed error information:

```go
import "github.com/jackc/pgx/v5/pgconn"

res, err := queries.SomeQuery(ctx, params)
if err != nil {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case "23505":  // Unique violation
            return fmt.Errorf("duplicate entry: %s", pgErr.ConstraintName)
        case "23503":  // Foreign key violation
            return fmt.Errorf("referenced record not found")
        default:
            return err
        }
    }
    return err
}
```

### 10.1 Common PostgreSQL Error Codes
| Code | Meaning |
|------|---------|
| 23505 | Unique violation |
| 23503 | Foreign key violation |
| 23502 | Not null violation |
| 42P01 | Undefined table |

Always consider using error codes rather than message text for logic branching, and avoid exposing raw database errors to end users for security reasons.

## 11. CI/CD Integration

For projects with multiple developers, run sqlc as part of your CI/CD pipeline using these four subcommands:

| Command | Purpose |
|---------|---------|
| `sqlc diff` | Ensures generated code is up to date with queries/schema |
| `sqlc vet` | Runs lint rules against SQL queries |
| `sqlc verify` | Checks that schema changes don't break existing queries |
| `sqlc push` | Pushes schema and queries to sqlc Cloud (after merge on main branch) |

### 11.1 GitHub Actions Example
```yaml
name: sqlc
on: [push]
jobs:
  diff:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: sqlc-dev/setup-sqlc@v3
        with:
          sqlc-version: '1.30.0'
      - run: sqlc diff
```

For `sqlc vet` with the built-in `sqlc/db-prepare` lint rule, you need a PostgreSQL server:
```yaml
- uses: sqlc-dev/action-setup-postgres@master
  with:
    postgres-version: "16"
- run: sqlc vet
  env:
    POSTGRESQL_SERVER_URI: ${{ steps.postgres.outputs.connection-uri }}?sslmode=disable
```


## 12. Type Overrides

To customize type mappings, add an `overrides` section to your configuration:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "schema.sql"
    queries: "query.sql"
    gen:
      go:
        package: "db"
        out: "db"
    overrides:
      - db_type: "bigint"
        go_type: "int64"
      - db_type: "timestamp"
        go_type: "time.Time"
```

This is particularly useful when working with cloud databases where sqlc may map types differently.

## 13. Best Practices

### 13.1 Query Writing
- **Avoid `SELECT *`**: Always explicitly list columns for clarity and safety. Cloud environments may prohibit wildcard queries entirely.
- **Use PostgreSQL placeholders**: Use `$1`, `$2`, etc., for parameters
- **Add comments**: Document complex queries with SQL comments

### 13.2 Project Structure Recommendation
```
project/
├── db/
│   ├── migrations/      # Migration files
│   ├── queries/         # .sql query files
│   │   └── authors.sql
│   ├── schema/          # Schema definition
│   │   └── schema.sql
│   └── sqlc.yaml        # sqlc configuration
├── internal/
│   └── db/              # Generated code (out directory)
├── go.mod
└── main.go
```


### 13.3 Development Workflow
1. Write or update schema in `schema.sql`
2. Write or update queries in `query.sql`
3. Run `sqlc generate` to regenerate code
4. Use the generated type-safe methods in your application
5. Run `sqlc diff` before committing to ensure generated code is up to date

### 13.4 PostgreSQL Engine Characteristics
- Most mature and feature-rich engine in sqlc
- Supports PostgreSQL-specific features like enums, arrays, and composite types
- Uses CGO and only available on Linux and macOS (not Windows)
- Full support for advanced features like CTEs and window functions

## 14. Troubleshooting Common Issues

| Issue | Solution |
|-------|----------|
| `sqlc generate` produces no output | Check that schema and query paths are correct; sqlc fails silently on parse errors |
| Type mismatches with pgx | Ensure `sql_package: "pgx/v5"` is set in configuration |
| Cloud database connection issues | URL-encode special characters in passwords; verify pg_catalog access |
| Multi-schema table not found | Explicitly prefix table names with schema name in queries |
| Migration files parsed incorrectly | Ensure filenames have consistent zero-padding for lexicographic ordering |

## 15. Additional Resources
- [Official sqlc Documentation](https://docs.sqlc.dev)
- [sqlc Cloud Dashboard](https://dashboard.sqlc.dev)
- [GitHub Repository](https://github.com/sqlc-dev/sqlc)
