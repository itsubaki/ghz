package beta

import (
	"context"
	"time"

	"google.golang.org/appengine/v2/memcache"
)

var ErrCacheMiss = memcache.ErrCacheMiss

func MemGet(c context.Context, key string) ([]byte, error) {
	item, err := memcache.Get(c, key)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

func MemSet(c context.Context, key string, value []byte, expire time.Duration) error {
	return memcache.Set(c, &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expire,
	})
}
