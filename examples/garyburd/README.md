# Garyburd redigo

```
type RedisLockClient struct {
	pool       *redis.Pool
	luaRefresh *redis.Script
	luaPttl    *redis.Script
	luaRelease *redis.Script
}
```

`RedisLockClient` implements `RedisClient` interface from `redislock.go`

## Installing the dependencies

```
dep ensure -v
```

## Running the application

```
go run main.go
```