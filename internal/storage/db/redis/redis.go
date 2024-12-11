package redis

import "github.com/redis/go-redis/v9"

type Cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) Client() *redis.Client {
	return c.rdb
}
