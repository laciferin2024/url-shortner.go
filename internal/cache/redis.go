package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/laciferin2024/url-shortner.go/utils"
	"github.com/pkg/errors"
)

var ErrNotFound = fmt.Errorf("key not found")
var ErrInvalidObject = fmt.Errorf("object provided is nil")
var ErrKeyExists = fmt.Errorf("key exists")

const (
	DEFAULT = time.Duration(0)
	FOREVER = time.Duration(-1)
)

// Get - get the key from redis
func (s *redisStore) Get(key string, obj interface{}) (err error) {

	conn, cancel := s.open()
	defer cancel()

	raw, err := conn.Do("GET", key)
	if err != nil {
		return
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return utils.Deserialize(item, obj)
}

// Set - set key to redis
func (s *redisStore) Set(key string, value interface{}, expires time.Duration) error {
	conn, cancel := s.open()
	defer cancel()
	return s.invoke(conn.Do, key, value, expires)
}

// Add - add key to redis
func (s *redisStore) Add(key string, value interface{}, expires time.Duration) error {
	conn, cancel := s.open()
	defer cancel()
	if exists(conn, key) {
		return errors.Wrap(ErrKeyExists, fmt.Sprintf("key=%s", key))
	}
	return s.invoke(conn.Do, key, value, expires)
}

// Delete - delete kee
func (s *redisStore) Delete(key string) error {
	conn, cancel := s.open()
	defer cancel()

	if !exists(conn, key) {
		return ErrNotFound
	}
	_, err := conn.Do("DEL", key)
	return err
}

// Flush - clear key store
func (s *redisStore) Flush() error {
	conn, cancel := s.open()
	defer cancel()
	_, err := conn.Do("FLUSHALL")
	return err
}

// open-helper function
func (s *redisStore) open() (conn redis.Conn, cancel context.CancelFunc) {
	conn = s.pool.Get()
	ctx, cancel := utils.CreateContext()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	return
}

// exists -helper function
func exists(conn redis.Conn, key string) bool {
	exist, _ := redis.Bool(conn.Do("EXISTS", key))
	return exist
}

// invoke - helper function
func (s *redisStore) invoke(callback func(string, ...interface{}) (interface{}, error),
	key string, value interface{}, expires time.Duration) error {

	switch expires {
	case DEFAULT:
		expires = s.defaultExpiration
	case FOREVER:
		expires = time.Duration(0)
	}

	b, err := utils.Serialize(value)
	if err != nil {
		return err
	}

	if expires > 0 {
		_, err := callback("SETEX", key, int32(expires/time.Second), b)
		return err
	}

	_, err = callback("SET", key, b)
	return err

}

func (s *redisStore) GetPool() interface{} {
	return s.pool
}
