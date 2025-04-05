package dtb

import (
	"anubis/tools/dtb"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func NewClientPgx(
	ctx context.Context,
	maxAttempts int,
	maxDelay time.Duration,
	databaseUrl string,
	maxpool int,
) (pool *pgxpool.Pool, err error) {

	pgxCfg, parseConfigErr := pgxpool.ParseConfig(databaseUrl)
	if parseConfigErr != nil {
		log.Printf("Unable to parse config: %v\n", parseConfigErr)
		return nil, parseConfigErr
	}
	pgxCfg.MaxConns = int32(maxpool)

	//if binary {
	//	pgxCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	//}

	pool, parseConfigErr = pgxpool.NewWithConfig(ctx, pgxCfg)
	if parseConfigErr != nil {
		log.Printf("Failed to parse PostgreSQL configuration due to error: %v\n", parseConfigErr)
		return nil, parseConfigErr
	}

	// check ping the pool
	err = dtb.DoWithAttempts(func() error {
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

	pool.Close()
	return pool, nil

}
