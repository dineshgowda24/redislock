package garyburd

import (
	"time"

	"github.com/dineshgowda24/redislock"
	"github.com/garyburd/redigo/redis"
)

type RedisLockClient struct {
	pool       *redis.Pool
	luaRefresh *redis.Script
	luaPttl    *redis.Script
	luaRelease *redis.Script
}

func NewRedisLockClient(pool *redis.Pool) *RedisLockClient {
	return &RedisLockClient{
		pool:       pool,
		luaRefresh: redis.NewScript(1, redislock.LuaRefreshScript),
		luaPttl:    redis.NewScript(1, redislock.LuaPTTLScript),
		luaRelease: redis.NewScript(1, redislock.LuaReleaseScript),
	}
}

func (r *RedisLockClient) SetNX(key, value string, ttl time.Duration) (bool, error) {
	con := r.pool.Get()
	defer con.Close()
	_, err := redis.String(con.Do("SET", key, value, "PX", ttl.Milliseconds(), "NX"))
	//Redigo returns nil so that means lock is not obtained so mask and return error
	if err == redis.ErrNil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisLockClient) Refresh(key, value string, ttl string) error {
	con := r.pool.Get()
	defer con.Close()

	status, err := redis.Int64(r.luaRefresh.Do(con, key, value, ttl))
	if err != nil {
		return err
	} else if status == 1 {
		return nil
	}
	//either the value did not match or key does not exist
	return redislock.ErrNotObtained
}

func (r *RedisLockClient) Release(key, value string) error {
	con := r.pool.Get()
	defer con.Close()

	res, err := redis.Int64(r.luaRelease.Do(con, key, value))
	if err == redis.ErrNil {
		return redislock.ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if res != 1 {
		return redislock.ErrLockNotHeld
	}
	return nil
}

func (r *RedisLockClient) TTL(key, value string) (int64, error) {
	con := r.pool.Get()
	defer con.Close()

	res, err := redis.Int64(r.luaPttl.Do(con, key, value))
	if err == redis.ErrNil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return res, nil
}
