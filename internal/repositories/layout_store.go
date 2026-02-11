package repositories

import (
	"crypto-aggregator-service/internal/models"
	"sync"
)

type LayoutStore struct {
	mu     sync.RWMutex
	layout []models.Component
}

// NewLayoutStore initializes the store with the config layout.
// This ensures the order is fixed from startup.
func NewLayoutStore(initialLayout []models.Component) *LayoutStore {
	// Deep copy to ensure safety
	safeLayout := make([]models.Component, len(initialLayout))
	copy(safeLayout, initialLayout)
	return &LayoutStore{layout: safeLayout}
}

// GetLayout returns a safe copy of the current state.
func (s *LayoutStore) GetLayout() []models.Component {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Component, len(s.layout))
	copy(result, s.layout)
	return result
}

// UpdateModel updates the model of a specific component by index.
// Access by index is O(1) and preserves order.
func (s *LayoutStore) UpdateModel(index int, model interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index >= 0 && index < len(s.layout) {
		s.layout[index].Model = model
	}
}
