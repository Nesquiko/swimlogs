package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

var ErrRowsNotFound = errors.New("didn't find row")

const (
	SmallIntMax = 32767
)

type PostgresDbPool struct {
	conStr        string
	migrationsDir string
	*pgxpool.Pool
}

func NewPostgresPool(conStr, migrationsDir string) (*PostgresDbPool, error) {
	dbConfig, err := pgxpool.ParseConfig(conStr)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresPool parse config: %w", err)
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresPool new: %w", err)
	}

	return &PostgresDbPool{conStr, migrationsDir, dbPool}, nil
}

func (p *PostgresDbPool) Reconnect(ctx context.Context) error {
	dbConfig, err := pgxpool.ParseConfig(p.conStr)
	if err != nil {
		return fmt.Errorf("Reconnect parse config: %w", err)
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return fmt.Errorf("Reconnect new: %w", err)
	}

	p.Pool = dbPool
	return nil
}

func (psql *PostgresDbPool) Close() {
	psql.Pool.Close()
}

func Sql(pool *PostgresDbPool, sql string, args ...any) error {
	_, err := pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("Sql: %w", err)
	}
	return nil
}

func SqlWithResult(pool *PostgresDbPool, sql string, args, dest []any) error {
	err := pool.QueryRow(context.Background(), sql, args...).Scan(dest...)
	if err != nil {
		return fmt.Errorf("SqlWithResult: %w", err)
	}
	return nil
}

func Tx(ctx context.Context, pool *PostgresDbPool, f func(context.Context, pgx.Tx) error) error {
	_, err := TxWithResult(ctx, pool, func(ctx context.Context, tx pgx.Tx) (struct{}, error) {
		return struct{}{}, f(ctx, tx)
	})
	if err != nil {
		return fmt.Errorf("Tx: %w", err)
	}
	return nil
}

func TxWithResult[R any](
	ctx context.Context,
	pool *PostgresDbPool,
	f func(context.Context, pgx.Tx) (R, error),
) (R, error) {
	var res R
	tx, err := pool.Begin(ctx)
	if err != nil {
		return res, fmt.Errorf("TxWithResult init: %w", err)
	}
	defer tx.Rollback(ctx)

	res, err = f(ctx, tx)
	if err != nil {
		return res, fmt.Errorf("TxWithResult: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return res, fmt.Errorf("TxWithResult commit: %w", err)
	}

	return res, nil
}

func ConnectionString(user, pass, host, db, port string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
}
