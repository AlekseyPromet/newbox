---
name: slq-transaction
description: This SKILL describes approaches to working with transactions in sqlc — a type-safe Go code generator from SQL queries. Unlike traditional ORMs, sqlc allows you to write raw SQL while obtaining efficient Go functions with static typing.
---

# Using Transactions in sqlc

## Instructions

## Core Principle: The WithTx Method

In the generated sqlc code, there is a `WithTx` method that allows you to associate a `Queries` instance with a specific transaction. A typical implementation looks like this:

```go
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
    return &Queries{ db: tx }
}
```

Because `sql.Tx` implements the `DBTX` interface, this approach works perfectly.

## Basic Transaction Workflow Structure

The general pattern for using transactions in sqlc includes several mandatory steps:

1. **Begin the transaction** via the database driver
2. **Deferred Rollback** to guarantee rollback on error
3. **Create a copy of Queries** bound to the transaction using `WithTx`
4. **Execute operations** via the transactional queries instance
5. **Commit** upon successful execution of all operations

### Example with standard database/sql

```go
import (
    "context"
    "database/sql"
    _ "github.com/lib/pq"
)

func bumpCounter(ctx context.Context, db *sql.DB, queries *tutorial.Queries, id int32) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Guaranteed rollback on error

    qtx := queries.WithTx(tx) // Create a copy bound to the transaction

    r, err := qtx.GetRecord(ctx, id)
    if err != nil {
        return err
    }

    if err := qtx.UpdateRecord(ctx, tutorial.UpdateRecordParams{
        ID:      r.ID,
        Counter: r.Counter + 1,
    }); err != nil {
        return err
    }

    return tx.Commit() // Commit changes
}
```

This code ensures atomicity: either both operations (read and update) succeed, or neither is applied to the database.

### Example for pgx/v5 driver

When using the `github.com/jackc/pgx/v5` library, the syntax differs slightly because `Begin` and `Commit` methods accept a context:

```go
import (
    "context"
    "github.com/jackc/pgx/v5"
)

func bumpCounter(ctx context.Context, db *pgx.Conn, queries *tutorial.Queries, id int32) error {
    tx, err := db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    qtx := queries.WithTx(tx)

    r, err := qtx.GetRecord(ctx, id)
    if err != nil {
        return err
    }

    if err := qtx.UpdateRecord(ctx, tutorial.UpdateRecordParams{
        ID:      r.ID,
        Counter: r.Counter + 1,
    }); err != nil {
        return err
    }

    return tx.Commit(ctx)
}
```

## Advanced Techniques

### Proper Rollback Error Handling

The standard `defer tx.Rollback()` approach has a nuance: the error returned by `Rollback()` may be lost. To capture it correctly, you can use a named return parameter and `errors.Join()`:

```go
func bumpCounter(ctx context.Context, db *sql.DB, queries *tutorial.Queries, id int32) (err error) {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, tx.Rollback())
    }()

    qtx := queries.WithTx(tx)

    // ... execute operations

    return tx.Commit()
}
```

However, note that after a successful `Commit`, the subsequent `Rollback` in `defer` may cause an error like "transaction has already been committed or rolled back". Therefore, a safer approach is to explicitly call `Rollback` only when errors occur.

### Row Locking to Prevent Race Conditions

To prevent race conditions, you can use the `SELECT ... FOR UPDATE` construct, which locks selected rows until the transaction completes:

```sql
-- name: GetCourtSlotForUpdate :one
SELECT * FROM court_slots WHERE id = $1 FOR UPDATE;
```

This is especially relevant for booking systems, payment processing, and other scenarios with high data integrity requirements.

## Best Practices

1. **Always use transactions for grouped operations**: If a process involves multiple logically related queries, wrap them in a transaction to ensure atomicity.

2. **Use `defer tx.Rollback()` as a safety net**: Even if you plan to call `Commit`, a deferred `Rollback` does no harm — after `Commit` it will be ignored, but it protects against leaks if you return early.

3. **Keep transactions short**: Long transactions increase the likelihood of deadlocks and hold locks longer than necessary.

4. **Avoid user input inside a transaction**: Do not wait for user response while holding an open transaction — this will lead to resource locking.

5. **Handle errors at every step**: Any operation inside a transaction can fail and must be handled properly with subsequent rollback.

6. **Use locks (`FOR UPDATE`) for critical sections**: In high‑concurrency scenarios, this prevents race conditions.

## Common Mistakes

| Mistake | Consequence | Solution |
|---------|-------------|----------|
| Omitting `defer tx.Rollback()` | On error, the transaction remains open, resource leak | Always add a deferred Rollback |
| Ignoring the Commit error | Changes may not be applied without notice | Check the error returned by `Commit()` |
| Mixing transactional and non‑transactional Queries instances | Some operations may run outside the transaction | Use `WithTx` to create a separate instance |
| Transactions that are too long | Locking other operations, reduced performance | Move non‑atomic operations outside the transaction |

## Typical Use Cases

- **Financial transactions**: Transferring funds between accounts, guaranteeing no double spending
- **Batch insertion**: When inserting many records, a failure in one should roll back the entire operation
- **Order processing**: Atomically updating inventory and creating an order record
- **User registration**: Inserting data into several related tables must be either fully successful or fully cancelled
- **Booking systems**: Preventing double booking of resources in a concurrent environment

## Conclusion

Transactions in sqlc are implemented via the `WithTx` method, which creates a copy of the `Queries` struct bound to the driver’s transaction object. This approach combines sqlc’s type safety with full control over transactional logic provided by standard Go facilities. When applying the described patterns correctly, you can build reliable systems with ACID‑level data integrity guarantees.
