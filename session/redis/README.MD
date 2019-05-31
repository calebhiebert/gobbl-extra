# GOBBL Redis Store

This is a session store for GOBBL that stores session data in a redis database

## Usage

```bash
go get -u github.com/calebhiebert/gobbl-redis-store
```

```go
redisStore := gobblredis.New(&redis.Options{
	Addr: "localhost:6379",
}, 0, "session:")
  
bot.Use(sess.Middleware(redisStore))
```