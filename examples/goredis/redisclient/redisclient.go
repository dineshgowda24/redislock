package goredis

import (
	"time"

	"github.com/dineshgowda24/redislock"
	"github.com/go-redis/redis/v7"
)

type RedisLockClient struct {
	client     *redis.Client
	luaRefresh *redis.Script
	luaPttl    *redis.Script
	luaRelease *redis.Script
}

func NewRedisLockClient(client *redis.Client) *RedisLockClient {
	return &RedisLockClient{
		client:     client,
		luaRefresh: redis.NewScript(redislock.LuaRefreshScript),
		luaPttl:    redis.NewScript(redislock.LuaPTTLScript),
		luaRelease: redis.NewScript(redislock.LuaReleaseScript),
	}
}

func (r *RedisLockClient) SetNX(key, value string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(key, value, ttl).Result()
}

func (r *RedisLockClient) Refresh(key, value string, ttl string) error {
	status, err := r.luaRefresh.Run(r.client, []string{key}, value, ttl).Result()
	if err != nil {
		return err
	} else if status == int64(1) {
		return nil
	}
	return redislock.ErrNotObtained

}

func (r *RedisLockClient) Release(key, value string) error {
	res, err := r.luaRelease.Run(r.client, []string{key}, value).Result()
	if err == redis.Nil {
		return redislock.ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return redislock.ErrLockNotHeld
	}
	return nil
}

func (r *RedisLockClient) TTL(key, value string) (int64, error) {
	res, err := r.luaPttl.Run(r.client, []string{key}, value).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return res.(int64), nil

}
