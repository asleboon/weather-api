package services

// In Memory Cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data       interface{}
	Expiration time.Time
}

type InMemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		items: make(map[string]CacheItem),
	}
	go cache.cleanup()
	return cache
}

func (cache *InMemoryCache) Set(key string, data interface{}, expiration time.Duration) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.items[key] = CacheItem{
		Data:       data,
		Expiration: time.Now().Add(expiration),
	}
}

func (cache *InMemoryCache) Get(key string) (interface{}, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	item, exists := cache.items[key]
	if !exists || item.Expiration.Before(time.Now()) {
		return nil, false
	}

	return item.Data, true
}

func (cache *InMemoryCache) cleanup() {
	for {
		time.Sleep(10 * time.Minute) // Adjust the cleanup interval as needed

		cache.mu.Lock()
		now := time.Now()
		for k, v := range cache.items {
			if v.Expiration.Before(now) {
				delete(cache.items, k)
			}
		}
		cache.mu.Unlock()
	}
}

// Redis Cache stuff

// CACHEING STUFF

// ctx := context.Background()
//
// rdb := redis.NewClient(&redis.Options{
// 	Addr:     "redis:6379",
// 	Password: "", // no password set
// 	DB:       0,  // use default DB
// })

// Get value from cache
// cacheKey := fmt.Sprintf("weather:%s", city)
// log.Println(cacheKey)
// cachedData, err := rdb.Get(ctx, cacheKey).Result()
// if err == nil {
// 	log.Println("Found in cache")
// 	c.String(200, cachedData)
// 	return
// }
// log.Println("NOT IN CACHE")

// Store new data in cache for 1 hour
// err = rdb.Set(ctx, cacheKey, weatherData, time.Hour).Err()
// if err != nil {
// 	c.JSON(500, gin.H{"error": err.Error()})
