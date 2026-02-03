package cache

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// ConfigCache provides local caching with Redis backing and Kafka-driven invalidation
// This cache implements a two-layer strategy:
// 1. In-memory cache for fastest access
// 2. Redis cache for distributed consistency
type ConfigCache struct {
	redis       *redis.Client
	localCache  map[string]CachedConfig
	mu          sync.RWMutex
	serviceName string
	ttl         time.Duration
}

// CachedConfig represents a cached configuration item
type CachedConfig struct {
	Key              string    `json:"key"`
	Name             string    `json:"name,omitempty"`
	Type             string    `json:"type,omitempty"`
	Value            string    `json:"value,omitempty"`
	FixedAmount      float64   `json:"fixed_amount,omitempty"`
	PercentageAmount float64   `json:"percentage_amount,omitempty"`
	Currency         string    `json:"currency,omitempty"`
	IsEnabled        bool      `json:"is_enabled"`
	CachedAt         time.Time `json:"cached_at"`
}

// NewConfigCache creates a new configuration cache
func NewConfigCache(redisClient *redis.Client, serviceName string) *ConfigCache {
	return &ConfigCache{
		redis:       redisClient,
		localCache:  make(map[string]CachedConfig),
		serviceName: serviceName,
		ttl:         5 * time.Minute,
	}
}

// NewConfigCacheWithTTL creates a new configuration cache with custom TTL
func NewConfigCacheWithTTL(redisClient *redis.Client, serviceName string, ttl time.Duration) *ConfigCache {
	return &ConfigCache{
		redis:       redisClient,
		localCache:  make(map[string]CachedConfig),
		serviceName: serviceName,
		ttl:         ttl,
	}
}

// Get retrieves a config from local cache, then Redis
// Returns nil if not found in either cache
func (c *ConfigCache) Get(ctx context.Context, key string) (*CachedConfig, error) {
	// 1. Check local in-memory cache first (fastest)
	c.mu.RLock()
	if cached, ok := c.localCache[key]; ok {
		if time.Since(cached.CachedAt) < c.ttl {
			c.mu.RUnlock()
			return &cached, nil
		}
	}
	c.mu.RUnlock()

	// 2. Check Redis (distributed cache)
	if c.redis != nil {
		redisKey := c.redisKey(key)
		data, err := c.redis.Get(ctx, redisKey).Bytes()
		if err == nil {
			var config CachedConfig
			if json.Unmarshal(data, &config) == nil {
				// Update local cache
				c.mu.Lock()
				c.localCache[key] = config
				c.mu.Unlock()
				return &config, nil
			}
		}
	}

	// 3. Not found - caller should fetch from source
	return nil, nil
}

// GetBool is a convenience method for boolean configs (system settings)
func (c *ConfigCache) GetBool(ctx context.Context, key string) (bool, bool) {
	config, err := c.Get(ctx, key)
	if err != nil || config == nil {
		return false, false // not found
	}
	return config.IsEnabled, true
}

// GetFloat is a convenience method for numeric configs (fees, limits)
func (c *ConfigCache) GetFloat(ctx context.Context, key string) (float64, bool) {
	config, err := c.Get(ctx, key)
	if err != nil || config == nil {
		return 0, false
	}
	// Return fixed amount if set, otherwise percentage
	if config.FixedAmount > 0 {
		return config.FixedAmount, true
	}
	return config.PercentageAmount, true
}

// Set updates both local cache and Redis
func (c *ConfigCache) Set(ctx context.Context, config CachedConfig) error {
	config.CachedAt = time.Now()

	// Update local cache
	c.mu.Lock()
	c.localCache[config.Key] = config
	c.mu.Unlock()

	// Update Redis
	if c.redis != nil {
		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		if err := c.redis.Set(ctx, c.redisKey(config.Key), data, c.ttl).Err(); err != nil {
			log.Printf("[ConfigCache] Warning: failed to set Redis cache for %s: %v", config.Key, err)
		}
	}

	return nil
}

// SetMultiple updates multiple configs at once
func (c *ConfigCache) SetMultiple(ctx context.Context, configs []CachedConfig) error {
	now := time.Now()

	c.mu.Lock()
	for i := range configs {
		configs[i].CachedAt = now
		c.localCache[configs[i].Key] = configs[i]
	}
	c.mu.Unlock()

	// Batch update Redis using pipeline
	if c.redis != nil {
		pipe := c.redis.Pipeline()
		for _, config := range configs {
			data, _ := json.Marshal(config)
			pipe.Set(ctx, c.redisKey(config.Key), data, c.ttl)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("[ConfigCache] Warning: failed to batch set Redis cache: %v", err)
		}
	}

	return nil
}

// Invalidate removes a config from both caches (triggered by Kafka event)
func (c *ConfigCache) Invalidate(ctx context.Context, key string) {
	c.mu.Lock()
	delete(c.localCache, key)
	c.mu.Unlock()

	if c.redis != nil {
		c.redis.Del(ctx, c.redisKey(key))
	}

	log.Printf("[ConfigCache] Invalidated config: %s", key)
}

// InvalidateByPrefix removes all configs matching a prefix
func (c *ConfigCache) InvalidateByPrefix(ctx context.Context, prefix string) {
	c.mu.Lock()
	for key := range c.localCache {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.localCache, key)
		}
	}
	c.mu.Unlock()

	if c.redis != nil {
		iter := c.redis.Scan(ctx, 0, c.redisKey(prefix+"*"), 100).Iterator()
		for iter.Next(ctx) {
			c.redis.Del(ctx, iter.Val())
		}
	}

	log.Printf("[ConfigCache] Invalidated configs with prefix: %s", prefix)
}

// InvalidateAll clears all cached configs
func (c *ConfigCache) InvalidateAll(ctx context.Context) {
	c.mu.Lock()
	c.localCache = make(map[string]CachedConfig)
	c.mu.Unlock()

	if c.redis != nil {
		iter := c.redis.Scan(ctx, 0, c.redisKey("*"), 100).Iterator()
		for iter.Next(ctx) {
			c.redis.Del(ctx, iter.Val())
		}
	}

	log.Printf("[ConfigCache] Invalidated all configs for service: %s", c.serviceName)
}

// Warmup preloads configs into the cache
func (c *ConfigCache) Warmup(ctx context.Context, configs []CachedConfig) {
	log.Printf("[ConfigCache] Warming up cache with %d configs", len(configs))
	c.SetMultiple(ctx, configs)
}

// Stats returns cache statistics
func (c *ConfigCache) Stats() map[string]interface{} {
	c.mu.RLock()
	localCount := len(c.localCache)
	c.mu.RUnlock()

	return map[string]interface{}{
		"service_name":     c.serviceName,
		"local_cache_size": localCount,
		"ttl_seconds":      c.ttl.Seconds(),
	}
}

// redisKey generates a namespaced Redis key
func (c *ConfigCache) redisKey(key string) string {
	return "config:" + c.serviceName + ":" + key
}
