// Package cache provides caching mechanisms for expensive operations in the Teapot CLI.
// It includes caching for project structure rendering to improve performance during
// real-time UI updates.
package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"teapot/internal/models"
)

// StructureCache provides thread-safe caching for project structure rendering.
// It caches rendered structures based on project configuration and terminal dimensions.
type StructureCache struct {
	// cache stores the cached rendered structures
	cache map[string]cacheEntry
	// mutex protects concurrent access to the cache
	mutex sync.RWMutex
	// maxSize limits the maximum number of cached entries
	maxSize int
	// ttl defines how long cache entries remain valid
	ttl time.Duration
}

// cacheEntry represents a single cached structure with metadata
type cacheEntry struct {
	// content is the rendered structure string
	content string
	// timestamp tracks when the entry was created
	timestamp time.Time
	// hits tracks how many times this entry has been accessed
	hits int
}

// cacheKey represents the key used for caching structure renders
type cacheKey struct {
	// projectHash is a hash of the project configuration
	ProjectHash string `json:"projectHash"`
	// terminalWidth is the terminal width used for rendering
	TerminalWidth int `json:"terminalWidth"`
	// terminalHeight is the terminal height used for rendering
	TerminalHeight int `json:"terminalHeight"`
}

// NewStructureCache creates a new structure cache with the specified configuration.
func NewStructureCache(maxSize int, ttl time.Duration) *StructureCache {
	return &StructureCache{
		cache:   make(map[string]cacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// GetStructure retrieves a cached structure or returns an empty string if not found.
// It also tracks cache hits and automatically evicts expired entries.
func (sc *StructureCache) GetStructure(project models.ProjectConfig, terminalWidth, terminalHeight int) (string, bool) {
	key := sc.generateKey(project, terminalWidth, terminalHeight)
	
	sc.mutex.RLock()
	entry, exists := sc.cache[key]
	sc.mutex.RUnlock()
	
	if !exists {
		return "", false
	}
	
	// Check if entry has expired
	if time.Since(entry.timestamp) > sc.ttl {
		sc.mutex.Lock()
		delete(sc.cache, key)
		sc.mutex.Unlock()
		return "", false
	}
	
	// Update hit count
	sc.mutex.Lock()
	entry.hits++
	sc.cache[key] = entry
	sc.mutex.Unlock()
	
	return entry.content, true
}

// SetStructure stores a rendered structure in the cache.
// It handles cache eviction when the maximum size is reached.
func (sc *StructureCache) SetStructure(project models.ProjectConfig, terminalWidth, terminalHeight int, content string) {
	key := sc.generateKey(project, terminalWidth, terminalHeight)
	
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	
	// Check if we need to evict entries
	if len(sc.cache) >= sc.maxSize {
		sc.evictLeastUsed()
	}
	
	// Store the new entry
	sc.cache[key] = cacheEntry{
		content:   content,
		timestamp: time.Now(),
		hits:      0,
	}
}

// Clear removes all cached entries.
func (sc *StructureCache) Clear() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	
	sc.cache = make(map[string]cacheEntry)
}

// GetStats returns cache statistics for monitoring and debugging.
func (sc *StructureCache) GetStats() CacheStats {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	
	totalHits := 0
	expiredCount := 0
	now := time.Now()
	
	for _, entry := range sc.cache {
		totalHits += entry.hits
		if now.Sub(entry.timestamp) > sc.ttl {
			expiredCount++
		}
	}
	
	return CacheStats{
		TotalEntries:  len(sc.cache),
		TotalHits:     totalHits,
		ExpiredCount:  expiredCount,
		MaxSize:       sc.maxSize,
		TTL:           sc.ttl,
	}
}

// CacheStats provides statistics about cache usage
type CacheStats struct {
	TotalEntries  int
	TotalHits     int
	ExpiredCount  int
	MaxSize       int
	TTL           time.Duration
}

// generateKey creates a unique key for caching based on project config and terminal dimensions
func (sc *StructureCache) generateKey(project models.ProjectConfig, terminalWidth, terminalHeight int) string {
	// Generate a hash of the project configuration
	projectHash := sc.hashProject(project)
	
	// Create cache key
	key := cacheKey{
		ProjectHash:    projectHash,
		TerminalWidth:  terminalWidth,
		TerminalHeight: terminalHeight,
	}
	
	// Convert to JSON string for map key
	keyBytes, _ := json.Marshal(key)
	return string(keyBytes)
}

// hashProject generates a hash of the project configuration for cache key generation
func (sc *StructureCache) hashProject(project models.ProjectConfig) string {
	// Create a simplified struct for hashing that includes only relevant fields
	hashData := struct {
		Name         string                   `json:"name"`
		Architecture models.ArchitectureType `json:"architecture"`
		Applications []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"applications"`
		Infrastructure struct {
			Docker        bool   `json:"docker"`
			DockerCompose bool   `json:"dockerCompose"`
			Pulumi        bool   `json:"pulumi"`
			Terraform     bool   `json:"terraform"`
			CloudProvider string `json:"cloudProvider"`
		} `json:"infrastructure"`
		CIPipeline struct {
			Provider string   `json:"provider"`
			Features []string `json:"features"`
		} `json:"ciPipeline"`
	}{
		Name:         project.Name,
		Architecture: project.Architecture,
		Infrastructure: struct {
			Docker        bool   `json:"docker"`
			DockerCompose bool   `json:"dockerCompose"`
			Pulumi        bool   `json:"pulumi"`
			Terraform     bool   `json:"terraform"`
			CloudProvider string `json:"cloudProvider"`
		}{
			Docker:        project.Infrastructure.Docker,
			DockerCompose: project.Infrastructure.DockerCompose,
			Pulumi:        project.Infrastructure.Pulumi,
			Terraform:     project.Infrastructure.Terraform,
			CloudProvider: project.Infrastructure.CloudProvider,
		},
		CIPipeline: struct {
			Provider string   `json:"provider"`
			Features []string `json:"features"`
		}{
			Provider: project.CIPipeline.Provider,
			Features: project.CIPipeline.Features,
		},
	}
	
	// Convert applications to simplified format
	for _, app := range project.Applications {
		hashData.Applications = append(hashData.Applications, struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}{
			Type: string(app.Type),
			Name: app.Name,
		})
	}
	
	// Generate JSON and hash it
	jsonData, _ := json.Marshal(hashData)
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}

// evictLeastUsed removes the least recently used entry from the cache
func (sc *StructureCache) evictLeastUsed() {
	if len(sc.cache) == 0 {
		return
	}
	
	var oldestKey string
	var oldestEntry cacheEntry
	isFirst := true
	
	// Find the entry with the lowest hits and oldest timestamp
	for key, entry := range sc.cache {
		if isFirst || entry.hits < oldestEntry.hits || 
			(entry.hits == oldestEntry.hits && entry.timestamp.Before(oldestEntry.timestamp)) {
			oldestKey = key
			oldestEntry = entry
			isFirst = false
		}
	}
	
	// Remove the oldest entry
	if oldestKey != "" {
		delete(sc.cache, oldestKey)
	}
}

// CleanupExpired removes all expired entries from the cache
func (sc *StructureCache) CleanupExpired() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	
	now := time.Now()
	for key, entry := range sc.cache {
		if now.Sub(entry.timestamp) > sc.ttl {
			delete(sc.cache, key)
		}
	}
}