package ports

import "time"

type CacheRepository interface {
	Get(key string) ([]byte, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
}
