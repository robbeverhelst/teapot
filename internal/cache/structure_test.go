package cache

import (
	"testing"
	"time"

	"teapot/internal/models"
)

func TestStructureCache_BasicOperations(t *testing.T) {
	cache := NewStructureCache(10, 5*time.Minute)
	
	// Test empty cache
	project := models.ProjectConfig{Name: "test-project"}
	content, found := cache.GetStructure(project, 80, 24)
	if found {
		t.Error("Expected cache miss for empty cache")
	}
	if content != "" {
		t.Error("Expected empty content for cache miss")
	}
	
	// Test cache set and get
	expectedContent := "test structure content"
	cache.SetStructure(project, 80, 24, expectedContent)
	
	content, found = cache.GetStructure(project, 80, 24)
	if !found {
		t.Error("Expected cache hit after setting")
	}
	if content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, content)
	}
	
	// Test cache miss with different dimensions
	content, found = cache.GetStructure(project, 100, 30)
	if found {
		t.Error("Expected cache miss with different dimensions")
	}
}

func TestStructureCache_DifferentProjects(t *testing.T) {
	cache := NewStructureCache(10, 5*time.Minute)
	
	// Create two different projects
	project1 := models.ProjectConfig{
		Name:         "project1",
		Architecture: models.ArchitectureTurborepo,
		Applications: []models.Application{
			{Type: models.AppTypeReact, Name: "web"},
		},
	}
	
	project2 := models.ProjectConfig{
		Name:         "project2",
		Architecture: models.ArchitectureSingle,
		Applications: []models.Application{
			{Type: models.AppTypeNest, Name: "api"},
		},
	}
	
	// Cache both projects
	cache.SetStructure(project1, 80, 24, "project1 structure")
	cache.SetStructure(project2, 80, 24, "project2 structure")
	
	// Test that both are cached separately
	content1, found1 := cache.GetStructure(project1, 80, 24)
	content2, found2 := cache.GetStructure(project2, 80, 24)
	
	if !found1 || !found2 {
		t.Error("Expected both projects to be cached")
	}
	
	if content1 != "project1 structure" {
		t.Errorf("Expected 'project1 structure', got '%s'", content1)
	}
	
	if content2 != "project2 structure" {
		t.Errorf("Expected 'project2 structure', got '%s'", content2)
	}
}

func TestStructureCache_TTLExpiration(t *testing.T) {
	cache := NewStructureCache(10, 100*time.Millisecond)
	
	project := models.ProjectConfig{Name: "test-project"}
	cache.SetStructure(project, 80, 24, "test content")
	
	// Should be found immediately
	content, found := cache.GetStructure(project, 80, 24)
	if !found {
		t.Error("Expected cache hit immediately after setting")
	}
	if content != "test content" {
		t.Errorf("Expected 'test content', got '%s'", content)
	}
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should be expired now
	content, found = cache.GetStructure(project, 80, 24)
	if found {
		t.Error("Expected cache miss after TTL expiration")
	}
	if content != "" {
		t.Error("Expected empty content after expiration")
	}
}

func TestStructureCache_MaxSizeEviction(t *testing.T) {
	cache := NewStructureCache(2, 5*time.Minute) // Small cache size
	
	project1 := models.ProjectConfig{Name: "project1"}
	project2 := models.ProjectConfig{Name: "project2"}
	project3 := models.ProjectConfig{Name: "project3"}
	
	// Fill cache to capacity
	cache.SetStructure(project1, 80, 24, "content1")
	cache.SetStructure(project2, 80, 24, "content2")
	
	// Both should be cached
	_, found1 := cache.GetStructure(project1, 80, 24)
	_, found2 := cache.GetStructure(project2, 80, 24)
	if !found1 || !found2 {
		t.Error("Expected both projects to be cached")
	}
	
	// Add third project, should evict least used
	cache.SetStructure(project3, 80, 24, "content3")
	
	// project1 should be evicted (least recently used)
	_, found1 = cache.GetStructure(project1, 80, 24)
	_, found2 = cache.GetStructure(project2, 80, 24)
	_, found3 := cache.GetStructure(project3, 80, 24)
	
	if found1 {
		t.Error("Expected project1 to be evicted")
	}
	if !found2 {
		t.Error("Expected project2 to still be cached")
	}
	if !found3 {
		t.Error("Expected project3 to be cached")
	}
}

func TestStructureCache_HitCounting(t *testing.T) {
	cache := NewStructureCache(10, 5*time.Minute)
	
	project := models.ProjectConfig{Name: "test-project"}
	cache.SetStructure(project, 80, 24, "test content")
	
	// Get initial stats
	stats := cache.GetStats()
	initialHits := stats.TotalHits
	
	// Access the cached item multiple times
	for i := 0; i < 5; i++ {
		cache.GetStructure(project, 80, 24)
	}
	
	// Check that hits increased
	stats = cache.GetStats()
	if stats.TotalHits != initialHits+5 {
		t.Errorf("Expected %d total hits, got %d", initialHits+5, stats.TotalHits)
	}
}

func TestStructureCache_Clear(t *testing.T) {
	cache := NewStructureCache(10, 5*time.Minute)
	
	// Add some entries
	project1 := models.ProjectConfig{Name: "project1"}
	project2 := models.ProjectConfig{Name: "project2"}
	
	cache.SetStructure(project1, 80, 24, "content1")
	cache.SetStructure(project2, 80, 24, "content2")
	
	// Verify they're cached
	stats := cache.GetStats()
	if stats.TotalEntries != 2 {
		t.Errorf("Expected 2 entries, got %d", stats.TotalEntries)
	}
	
	// Clear cache
	cache.Clear()
	
	// Verify cache is empty
	stats = cache.GetStats()
	if stats.TotalEntries != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", stats.TotalEntries)
	}
	
	// Verify entries are not found
	_, found1 := cache.GetStructure(project1, 80, 24)
	_, found2 := cache.GetStructure(project2, 80, 24)
	
	if found1 || found2 {
		t.Error("Expected no entries to be found after clear")
	}
}

func TestStructureCache_Stats(t *testing.T) {
	cache := NewStructureCache(10, 5*time.Minute)
	
	// Test empty cache stats
	stats := cache.GetStats()
	if stats.TotalEntries != 0 {
		t.Errorf("Expected 0 entries, got %d", stats.TotalEntries)
	}
	if stats.TotalHits != 0 {
		t.Errorf("Expected 0 hits, got %d", stats.TotalHits)
	}
	if stats.MaxSize != 10 {
		t.Errorf("Expected max size 10, got %d", stats.MaxSize)
	}
	if stats.TTL != 5*time.Minute {
		t.Errorf("Expected TTL 5m, got %v", stats.TTL)
	}
	
	// Add some entries and hits
	project := models.ProjectConfig{Name: "test-project"}
	cache.SetStructure(project, 80, 24, "test content")
	cache.GetStructure(project, 80, 24) // Generate hit
	cache.GetStructure(project, 80, 24) // Generate another hit
	
	stats = cache.GetStats()
	if stats.TotalEntries != 1 {
		t.Errorf("Expected 1 entry, got %d", stats.TotalEntries)
	}
	if stats.TotalHits != 2 {
		t.Errorf("Expected 2 hits, got %d", stats.TotalHits)
	}
}

func TestStructureCache_CleanupExpired(t *testing.T) {
	cache := NewStructureCache(10, 100*time.Millisecond)
	
	project1 := models.ProjectConfig{Name: "project1"}
	project2 := models.ProjectConfig{Name: "project2"}
	
	// Add entries
	cache.SetStructure(project1, 80, 24, "content1")
	time.Sleep(50 * time.Millisecond) // Wait half TTL
	cache.SetStructure(project2, 80, 24, "content2")
	
	// Wait for first entry to expire
	time.Sleep(75 * time.Millisecond)
	
	// Cleanup expired entries
	cache.CleanupExpired()
	
	// Check that expired entry is gone but recent entry remains
	_, found1 := cache.GetStructure(project1, 80, 24)
	_, found2 := cache.GetStructure(project2, 80, 24)
	
	if found1 {
		t.Error("Expected expired entry to be cleaned up")
	}
	if !found2 {
		t.Error("Expected recent entry to remain after cleanup")
	}
}