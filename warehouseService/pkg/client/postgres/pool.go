package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	//Begin(ctx context.Context) (pgx.Tx, error)
	//BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Db struct {
	pool *pgxpool.Pool
}

type Cfg struct {
	MaxAttempts int
	Login       string
	Password    string
	Host        string
	Port        string
	Database    string
}

func NewClient(ctx context.Context, cfg Cfg) (*Db, error) {
	var pool *pgxpool.Pool
	var err error
	//todo: build dns from config
	//dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", "postgres", "localhost", "5432", "postgres")
	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	log.Println(dns)
	err = again(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dns)
		if err != nil {
			return err
		}

		err = pool.Ping(ctx)
		if err != nil {
			return err
		}
		return nil
	}, cfg.MaxAttempts, 5*time.Second)
	if err != nil {
		log.Fatal("tries limit exceeded")
	}
	return &Db{pool: pool}, nil
}

func (d *Db) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	txx := fromContext(ctx)
	if txx != nil && txx.IsActive {
		return txx.Tx.Exec(ctx, sql, arguments...)
	}

	return d.pool.Exec(ctx, sql, arguments...)
}

func (d *Db) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	txx := fromContext(ctx)
	if txx != nil && txx.IsActive {
		rows, err := txx.Tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		//defer rows.Close()

		return rows, nil
	}
	rows, err := d.pool.Query(ctx, query, args...)
	//defer rows.Close()
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (d *Db) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	tx := fromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.QueryRow(ctx, query, args...)
	}

	return d.pool.QueryRow(ctx, query, args...)
}
