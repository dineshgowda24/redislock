# go-redis

```
type RedisLockClient struct {
	client     *redis.Client
	luaRefresh *redis.Script
	luaPttl    *redis.Script
	luaRelease *redis.Script
}
```

`RedisLockClient` implements `RedisClient` interface from [redislock.go](../../../../redislock.go)

## Running the application

```
go run main.go
```