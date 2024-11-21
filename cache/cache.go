package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
)

type CacheService interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Delete(key string) error
}

type Cache struct {
	CachedData *bigcache.BigCache
	Mutex      *sync.RWMutex
	Expiry     time.Time
}

func NewCache() CacheService {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

	return &Cache{
		CachedData: cache,
		Mutex:      &sync.RWMutex{},
		Expiry:     time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, error) {
	log.Println("Getting cache key: ", key)
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	val, err := c.CachedData.Get(key)
	if err != nil {
		log.Println("Error while getting cache key: ", err)
		return nil, err
	}

	if c.Expiry.Before(time.Now()) {
		log.Println("Cache expired")
		return nil, err
	}
	log.Println("Cache not expired")
	return val, err
}

func (c *Cache) Set(key string, value []byte) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	log.Println("Setting cache key: ", key)
	err := c.CachedData.Set(key, value)
	if err != nil {
		log.Println("Error while setting cache key: ", key)
		return err
	}
	log.Println("Cache key set: ", key)
	c.Expiry = time.Now().Add(5 * time.Minute)
	return nil
}

func (c *Cache) Delete(key string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	log.Println("Deleting cache key: ", key)
	return nil
}
