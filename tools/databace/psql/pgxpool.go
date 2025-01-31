package psql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func NewClient(
	ctx context.Context,
	maxAttempts int,
	maxDelay time.Duration,
	databaseUrl string,
	binary bool,
) (pool *pgxpool.Pool, err error) {
	// check the configs
	pgxCfg, parseConfigErr := pgxpool.ParseConfig(databaseUrl)
	if parseConfigErr != nil {
		log.Printf("Unable to parse settings: %v\n", parseConfigErr)
		return nil, parseConfigErr
	}
	// a special mod that disables caching and all sorts of perks.
	if binary {
		pgxCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	}
	// create pool
	pool, parseConfigErr = pgxpool.NewWithConfig(ctx, pgxCfg)
	if parseConfigErr != nil {
		log.Printf("Failed to parse PostgreSQL configuration due to error: %v\n", parseConfigErr)
		return nil, parseConfigErr
	}

	// check ping the pool
	err = DoWithAttempts(func() error {
		pingErr := pool.Ping(ctx)
		if pingErr != nil {
			log.Printf("Failed to connect to postgres due to error %v... Going to do the next attempt\n", pingErr)
			return pingErr
		}

		return nil
	}, maxAttempts, maxDelay)
	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to PostgreSQL")
	}

	return pool, nil

}

func DoWithAttempts(callback func() error, maxAttempts int, delay time.Duration) error {
	// function executes callback several times
	var err error
	for maxAttempts > 0 {
		if err = callback(); err != nil {
			time.Sleep(delay)
			maxAttempts--
			continue
		}
		return nil
	}
	return err
}
