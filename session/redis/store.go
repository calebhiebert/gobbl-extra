package gobblredis

import (
	"time"

	"github.com/calebhiebert/gobbl/session"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

// RedisStore impliments the GOBBL session store to store data in a redis database
type RedisStore struct {
	client    *redis.Client
	keyExpiry time.Duration
	keyPrefix string
}

// New creates a new Redis Store object
func New(opts *redis.Options, keyExpiry time.Duration, keyPrefix string) *RedisStore {
	client := redis.NewClient(opts)

	return &RedisStore{client: client, keyExpiry: keyExpiry, keyPrefix: keyPrefix}
}

// Create creates a new session stored in redis
func (r *RedisStore) Create(id string, data *map[string]interface{}) error {

	b, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	_, err = r.client.Set(r.keyPrefix+id, b, r.keyExpiry).Result()
	if err != nil {
		return err
	}

	return nil
}

// Get returns the session data from the redis store
func (r *RedisStore) Get(id string) (map[string]interface{}, error) {
	b, err := r.client.Get(r.keyPrefix + id).Bytes()
	if err != nil {

		if err.Error() == "redis: nil" {
			return nil, sess.ErrSessionNonexistant
		}

		return nil, err
	}

	var sessionData = make(map[string]interface{})

	err = msgpack.Unmarshal(b, &sessionData)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

// Update overwrites an existing session value
func (r *RedisStore) Update(id string, data *map[string]interface{}) error {
	return r.Create(id, data)
}

// Destroy will completely delete the session
func (r *RedisStore) Destroy(id string) error {
	_, err := r.client.Del(r.keyPrefix + id).Result()
	return err
}
