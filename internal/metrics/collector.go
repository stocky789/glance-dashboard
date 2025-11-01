package metrics

import (
	"runtime"
	"sync"
	"time"
)

type WidgetMetrics struct {
	WidgetID         string        `json:"widget_id"`
	UpdateCount      int64         `json:"update_count"`
	ErrorCount       int64         `json:"error_count"`
	LastUpdateTime   time.Time     `json:"last_update_time"`
	AverageUpdateMS  float64       `json:"average_update_ms"`
	LastError        string        `json:"last_error,omitempty"`
	updateTimes      []time.Duration
	mu               sync.RWMutex
}

type Collector struct {
	widgets map[string]*WidgetMetrics
	mu      sync.RWMutex
}

func NewCollector() *Collector {
	return &Collector{
		widgets: make(map[string]*WidgetMetrics),
	}
}

func (c *Collector) RecordUpdate(widgetID string, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics, exists := c.widgets[widgetID]
	if !exists {
		metrics = &WidgetMetrics{
			WidgetID:    widgetID,
			updateTimes: make([]time.Duration, 0, 100),
		}
		c.widgets[widgetID] = metrics
	}

	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	metrics.UpdateCount++
	metrics.LastUpdateTime = time.Now()
	metrics.updateTimes = append(metrics.updateTimes, duration)

	// Keep only last 100 update times for average calculation
	if len(metrics.updateTimes) > 100 {
		metrics.updateTimes = metrics.updateTimes[1:]
	}

	// Calculate average
	var total time.Duration
	for _, d := range metrics.updateTimes {
		total += d
	}
	metrics.AverageUpdateMS = float64(total) / float64(len(metrics.updateTimes)) / float64(time.Millisecond)
}

func (c *Collector) RecordError(widgetID string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics, exists := c.widgets[widgetID]
	if !exists {
		metrics = &WidgetMetrics{
			WidgetID:    widgetID,
			updateTimes: make([]time.Duration, 0, 100),
		}
		c.widgets[widgetID] = metrics
	}

	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	metrics.ErrorCount++
	if err != nil {
		metrics.LastError = err.Error()
	}
}

func (c *Collector) GetMetrics(widgetID string) *WidgetMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics, exists := c.widgets[widgetID]
	if !exists {
		return nil
	}

	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	// Return a copy
	return &WidgetMetrics{
		WidgetID:        metrics.WidgetID,
		UpdateCount:     metrics.UpdateCount,
		ErrorCount:      metrics.ErrorCount,
		LastUpdateTime:  metrics.LastUpdateTime,
		AverageUpdateMS: metrics.AverageUpdateMS,
		LastError:       metrics.LastError,
	}
}

func (c *Collector) GetAllMetrics() map[string]*WidgetMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*WidgetMetrics)
	for id := range c.widgets {
		result[id] = c.GetMetrics(id)
	}
	return result
}

func (c *Collector) GetSystemMetrics() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"memory": map[string]interface{}{
			"alloc_mb":       m.Alloc / 1024 / 1024,
			"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
			"sys_mb":         m.Sys / 1024 / 1024,
			"num_gc":         m.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
	}
}
