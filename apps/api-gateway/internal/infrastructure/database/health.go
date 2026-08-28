package database

import (
	"context"
	"time"
)

func (d *Database) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return d.Pool.Ping(ctx)
}
