package postgresql

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pgmatcher/internal/model"
)

type PgStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dbURL string) (*PgStorage, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect:%w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("connect:%w", err)
	}

	return &PgStorage{
		pool: pool,
	}, nil
}

func (r *PgStorage) Close() {
	r.pool.Close()
}

func (r *PgStorage) Store(ctx context.Context, bids []*model.Rule) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("pool conn:%w", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}
	for _, bid := range bids {
		batch.Queue("insert into rules(id, price, conds) values($1, $2, $3)", uuid.New().String(), bid.Price, bid.Conditions)
	}

	br := conn.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("exec batch:%w", err)
	}

	return nil
}

func (r *PgStorage) GetAllConditions(ctx context.Context) ([]string, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool conn:%w", err)
	}
	defer conn.Release()

	var conds []string

	const q = `select distinct unnest(conds) as  cond from rules`
	if err := pgxscan.Select(ctx, conn, &conds, q); err != nil {
		return nil, err
	}

	return conds, nil
}

func (r *PgStorage) GetAllowRules(ctx context.Context, price int, allowConds []string) ([]model.Rule, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	var x []model.Rule
	const query = `SELECT id, price,conds FROM rules where price between $1 and $2 and conds  <@ $3 limit 5`

	if err := pgxscan.Select(ctx, conn, &x, query, price-10, price+10, allowConds); err != nil {
		return nil, err
	}

	return x, nil
}
