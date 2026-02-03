package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func setupTestCache(t *testing.T) (*ConfigCache, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := NewConfigCache(client, "test-service")
	return cache, mr
}

func TestConfigCache_SetAndGet(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	config := CachedConfig{
		Key:       "test_key",
		Name:      "Test Config",
		IsEnabled: true,
	}

	// Set
	err := cache.Set(ctx, config)
	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Get
	result, err := cache.Get(ctx, "test_key")
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if result == nil {
		t.Fatal("Expected config, got nil")
	}

	if result.Key != "test_key" {
		t.Errorf("Expected key 'test_key', got '%s'", result.Key)
	}

	if !result.IsEnabled {
		t.Error("Expected IsEnabled to be true")
	}
}

func TestConfigCache_GetBool(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	// Set a boolean config
	cache.Set(ctx, CachedConfig{
		Key:       "feature_enabled",
		IsEnabled: true,
	})

	value, found := cache.GetBool(ctx, "feature_enabled")
	if !found {
		t.Error("Expected to find config")
	}
	if !value {
		t.Error("Expected value to be true")
	}

	// Test not found
	_, found = cache.GetBool(ctx, "nonexistent")
	if found {
		t.Error("Expected not to find nonexistent config")
	}
}

func TestConfigCache_GetFloat(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	// Set a fee config
	cache.Set(ctx, CachedConfig{
		Key:         "transfer_fee",
		FixedAmount: 2.5,
	})

	value, found := cache.GetFloat(ctx, "transfer_fee")
	if !found {
		t.Error("Expected to find config")
	}
	if value != 2.5 {
		t.Errorf("Expected 2.5, got %f", value)
	}
}

func TestConfigCache_Invalidate(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	// Set
	cache.Set(ctx, CachedConfig{Key: "to_delete", IsEnabled: true})

	// Verify it exists
	result, _ := cache.Get(ctx, "to_delete")
	if result == nil {
		t.Fatal("Expected config to exist")
	}

	// Invalidate
	cache.Invalidate(ctx, "to_delete")

	// Verify it's gone
	result, _ = cache.Get(ctx, "to_delete")
	if result != nil {
		t.Error("Expected config to be invalidated")
	}
}

func TestConfigCache_InvalidateByPrefix(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	// Set multiple configs
	cache.Set(ctx, CachedConfig{Key: "transfer_internal", IsEnabled: true})
	cache.Set(ctx, CachedConfig{Key: "transfer_external", IsEnabled: true})
	cache.Set(ctx, CachedConfig{Key: "card_fee", IsEnabled: true})

	// Invalidate by prefix
	cache.InvalidateByPrefix(ctx, "transfer_")

	// transfer_ configs should be gone
	result, _ := cache.Get(ctx, "transfer_internal")
	if result != nil {
		t.Error("Expected transfer_internal to be invalidated")
	}

	// card_fee should still exist
	result, _ = cache.Get(ctx, "card_fee")
	if result == nil {
		t.Error("Expected card_fee to still exist")
	}
}

func TestConfigCache_TTLExpiration(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Create cache with very short TTL
	cache := NewConfigCacheWithTTL(client, "test-service", 100*time.Millisecond)

	ctx := context.Background()

	cache.Set(ctx, CachedConfig{Key: "short_lived", IsEnabled: true})

	// Should be available immediately
	result, _ := cache.Get(ctx, "short_lived")
	if result == nil {
		t.Fatal("Expected config to exist")
	}

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Local cache should be expired
	cache.mu.RLock()
	cached, _ := cache.localCache["short_lived"]
	cache.mu.RUnlock()

	if time.Since(cached.CachedAt) < cache.ttl {
		t.Error("Expected local cache TTL to be expired")
	}
}

func TestConfigCache_Stats(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	// Add some configs
	cache.Set(ctx, CachedConfig{Key: "config1"})
	cache.Set(ctx, CachedConfig{Key: "config2"})

	stats := cache.Stats()

	if stats["service_name"] != "test-service" {
		t.Errorf("Expected service name 'test-service', got %v", stats["service_name"])
	}

	if stats["local_cache_size"].(int) != 2 {
		t.Errorf("Expected local cache size 2, got %v", stats["local_cache_size"])
	}
}
