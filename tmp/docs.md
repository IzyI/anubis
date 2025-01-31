```go

type Client interface {
    Close()
    Acquire(ctx context.Context) (*pgxpool.Conn, error)
    AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
    AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
    Stat() *pgxpool.Stat
    Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
    Begin(ctx context.Context) (pgx.Tx, error)
    BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
    }

```