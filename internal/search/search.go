package search

import (
	"strings"
	"time"
)

type SearchFilter struct {
	Query       string
	WidgetType  string
	DateFrom    time.Time
	DateTo      time.Time
	SortBy      string
	SortOrder   string
	Limit       int
	Offset      int
}

type SearchResult struct {
	ID        string
	Title     string
	Type      string
	Relevance float64
	Data      map[string]interface{}
}

type Searcher struct {
	data []map[string]interface{}
}

func NewSearcher(data []map[string]interface{}) *Searcher {
	return &Searcher{
		data: data,
	}
}

func (s *Searcher) Search(filter SearchFilter) []SearchResult {
	var results []SearchResult

	for _, item := range s.data {
		if s.matches(item, filter) {
			result := SearchResult{
				ID:   s.getStringField(item, "id"),
				Type: s.getStringField(item, "type"),
				Data: item,
			}
			if s.getStringField(item, "name") != "" {
				result.Title = s.getStringField(item, "name")
			}
			result.Relevance = s.calculateRelevance(item, filter)
			results = append(results, result)
		}
	}

	// Sort results
	if filter.SortBy != "" {
		s.sortResults(results, filter)
	}

	// Apply pagination
	start := filter.Offset
	end := start + filter.Limit

	if start >= len(results) {
		return []SearchResult{}
	}

	if end > len(results) {
		end = len(results)
	}

	return results[start:end]
}

func (s *Searcher) matches(item map[string]interface{}, filter SearchFilter) bool {
	// Query matching
	if filter.Query != "" {
		if !s.searchFields(item, filter.Query) {
			return false
		}
	}

	// Widget type matching
	if filter.WidgetType != "" {
		if s.getStringField(item, "type") != filter.WidgetType {
			return false
		}
	}

	// Date range matching
	if !filter.DateFrom.IsZero() || !filter.DateTo.IsZero() {
		createdAt := s.getTimeField(item, "created_at")
		if !filter.DateFrom.IsZero() && createdAt.Before(filter.DateFrom) {
			return false
		}
		if !filter.DateTo.IsZero() && createdAt.After(filter.DateTo) {
			return false
		}
	}

	return true
}

func (s *Searcher) searchFields(item map[string]interface{}, query string) bool {
	query = strings.ToLower(query)
	searchableFields := []string{"id", "name", "type", "description"}

	for _, field := range searchableFields {
		value := s.getStringField(item, field)
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}

	return false
}

func (s *Searcher) calculateRelevance(item map[string]interface{}, filter SearchFilter) float64 {
	relevance := 0.0

	if filter.Query != "" {
		query := strings.ToLower(filter.Query)
		name := strings.ToLower(s.getStringField(item, "name"))
		id := strings.ToLower(s.getStringField(item, "id"))

		if name == query {
			relevance += 1.0
		} else if strings.HasPrefix(name, query) {
			relevance += 0.8
		} else if strings.Contains(name, query) {
			relevance += 0.5
		} else if strings.Contains(id, query) {
			relevance += 0.3
		}
	}

	return relevance
}

func (s *Searcher) sortResults(results []SearchResult, filter SearchFilter) {
	// Implement sorting based on SortBy and SortOrder
}

func (s *Searcher) getStringField(item map[string]interface{}, field string) string {
	if value, ok := item[field]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Searcher) getTimeField(item map[string]interface{}, field string) time.Time {
	if value, ok := item[field]; ok {
		if t, ok := value.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}
