package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"scaffold/internal/pkg/logger"
	"scaffold/internal/repo/redis"

	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// TagInfo contains tag name and annotation information
type TagInfo struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// CachedTags cached tag information
type CachedTags struct {
	TemplateID string    `json:"template_id"`
	Tags       []TagInfo `json:"tags"`
	Timestamp  time.Time `json:"timestamp"`
	ExpireTime time.Time `json:"expire_time"`
}

// TemplateCache manages template tag caching with Redis primary and file fallback
type TemplateCache struct {
	cacheDir     string
	redisEnabled bool
	mu           sync.RWMutex
	memCache     map[string]*CachedTags
}

const (
	cacheKeyPrefix = "scaffold:tags:"
	cacheTTL       = time.Hour
)

// NewTemplateCache creates a new template cache manager
func NewTemplateCache(cacheDir string, redisEnabled bool) *TemplateCache {
	if cacheDir == "" {
		homeDir, _ := os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".scaffold", "template_cache")
	}

	tc := &TemplateCache{
		cacheDir:     cacheDir,
		redisEnabled: redisEnabled,
		memCache:     make(map[string]*CachedTags),
	}

	os.MkdirAll(cacheDir, 0755)
	tc.loadFileCache()

	return tc
}

// GetTags retrieves tags with caching (Redis -> memory -> file -> remote)
func (tc *TemplateCache) GetTags(templateID string, fetchFn func() ([]TagInfo, error)) ([]TagInfo, error) {
	// Try Redis first
	if tc.redisEnabled && redis.Client != nil {
		tags, err := tc.getFromRedis(templateID)
		if err == nil && tags != nil {
			return tags, nil
		}
	}

	// Try memory cache
	tc.mu.RLock()
	cached, exists := tc.memCache[templateID]
	tc.mu.RUnlock()

	if exists && cached.ExpireTime.After(time.Now()) {
		return cached.Tags, nil
	}

	// Fetch from remote
	tags, err := fetchFn()
	if err != nil {
		// If fetch fails but we have old cache, use it
		if exists {
			logger.Warn("Failed to fetch tags, using cached data",
				zap.String("templateID", templateID),
				zap.Error(err))
			return cached.Tags, nil
		}
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	// Update all caches
	tc.setCache(templateID, tags)

	return tags, nil
}

func (tc *TemplateCache) getFromRedis(templateID string) ([]TagInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := cacheKeyPrefix + templateID
	data, err := redis.Get(ctx, key)
	if err != nil {
		if err == redisclient.Nil {
			return nil, nil
		}
		return nil, err
	}

	var tags []TagInfo
	if err := json.Unmarshal([]byte(data), &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (tc *TemplateCache) setCache(templateID string, tags []TagInfo) {
	now := time.Now()
	cached := &CachedTags{
		TemplateID: templateID,
		Tags:       tags,
		Timestamp:  now,
		ExpireTime: now.Add(cacheTTL),
	}

	// Update memory cache
	tc.mu.Lock()
	tc.memCache[templateID] = cached
	tc.mu.Unlock()

	// Update Redis cache
	if tc.redisEnabled && redis.Client != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			data, err := json.Marshal(tags)
			if err != nil {
				return
			}
			key := cacheKeyPrefix + templateID
			redis.Set(ctx, key, string(data), cacheTTL)
		}()
	}

	// Save to file as fallback
	go tc.saveFileEntry(templateID, cached)
}

// Refresh forces a cache refresh for a template
func (tc *TemplateCache) Refresh(templateID string, fetchFn func() ([]TagInfo, error)) error {
	tags, err := fetchFn()
	if err != nil {
		return fmt.Errorf("failed to refresh tags: %w", err)
	}

	tc.setCache(templateID, tags)
	return nil
}

// Clear removes all cached data
func (tc *TemplateCache) Clear() error {
	tc.mu.Lock()
	tc.memCache = make(map[string]*CachedTags)
	tc.mu.Unlock()

	entries, err := os.ReadDir(tc.cacheDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			os.Remove(filepath.Join(tc.cacheDir, entry.Name()))
		}
	}

	return nil
}

// ClearTemplate removes cached data for a specific template
func (tc *TemplateCache) ClearTemplate(templateID string) error {
	tc.mu.Lock()
	delete(tc.memCache, templateID)
	tc.mu.Unlock()

	if tc.redisEnabled && redis.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		redis.Del(ctx, cacheKeyPrefix+templateID)
	}

	filename := filepath.Join(tc.cacheDir, templateID+".json")
	return os.Remove(filename)
}

func (tc *TemplateCache) saveFileEntry(templateID string, cached *CachedTags) {
	filename := filepath.Join(tc.cacheDir, templateID+".json")
	data, err := json.Marshal(cached)
	if err != nil {
		return
	}
	os.WriteFile(filename, data, 0644)
}

func (tc *TemplateCache) loadFileCache() {
	entries, err := os.ReadDir(tc.cacheDir)
	if err != nil {
		return
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		templateID := entry.Name()[:len(entry.Name())-5]
		filename := filepath.Join(tc.cacheDir, entry.Name())

		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var cached CachedTags
		if err := json.Unmarshal(data, &cached); err != nil {
			continue
		}

		if cached.ExpireTime.After(now) {
			tc.memCache[templateID] = &cached
		} else {
			os.Remove(filename)
		}
	}
}
