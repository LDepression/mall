package query

import "github.com/go-redis/redis/v8"

type Queries struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Queries {
	return &Queries{rdb: rdb}
}
