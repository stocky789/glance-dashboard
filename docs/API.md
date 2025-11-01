# Glance API Documentation

## Overview

The Glance API provides real-time widget data, metrics, and WebSocket connections for dashboard functionality.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. Authentication will be added in future releases.

## Endpoints

### Health Check

**GET** `/health`

Returns the health status of the API.

**Response:**
```json
{
  "status": "healthy",
  "uptime": 3600
}
```

### Metrics

**GET** `/metrics`

Returns real-time performance metrics for the system and widgets.

**Query Parameters:**
- `widget_id` (optional): Filter metrics for a specific widget

**Response:**
```json
{
  "system": {
    "memory": 52428800,
    "goroutines": 42,
    "uptime": 3600000
  },
  "widgets": {
    "widget_1": {
      "updateCount": 150,
      "errorCount": 2,
      "avgUpdateDuration": 125.5
    }
  }
}
```

### Widget Metrics

**GET** `/metrics/widgets/{id}`

Returns detailed metrics for a specific widget.

**Response:**
```json
{
  "widgetId": "widget_1",
  "updateCount": 150,
  "errorCount": 2,
  "avgUpdateDuration": 125.5,
  "lastUpdate": "2024-01-15T10:30:00Z"
}
```

### Search

**GET** `/search`

Search widgets by query, type, and date range.

**Query Parameters:**
- `q` (optional): Search query
- `type` (optional): Widget type
- `from` (optional): Start date (YYYY-MM-DD)
- `to` (optional): End date (YYYY-MM-DD)
- `limit` (optional): Result limit (default: 10)
- `offset` (optional): Result offset (default: 0)

**Response:**
```json
{
  "results": [
    {
      "id": "widget_1",
      "title": "Weather Widget",
      "type": "weather",
      "relevance": 0.95,
      "data": {}
    }
  ]
}
```

### Activity Log

**GET** `/activity`

Retrieve activity logs.

**Query Parameters:**
- `limit` (optional): Number of logs (default: 50)
- `types` (optional): Comma-separated event types to filter

**Response:**
```json
{
  "logs": [
    {
      "id": 1,
      "eventType": "widget_update",
      "widgetId": "widget_1",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## WebSocket

**GET** `/ws`

Upgrade HTTP connection to WebSocket for real-time updates.

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws');

ws.onopen = () => {
  console.log('Connected');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Update:', data);
};

ws.onerror = (error) => {
  console.error('Error:', error);
};
```

### Message Format

**Widget Update:**
```json
{
  "type": "widget_update",
  "widgetId": "widget_1",
  "data": {}
}
```

**Metrics Update:**
```json
{
  "type": "metrics_update",
  "system": {},
  "widgets": {}
}
```

**Error:**
```json
{
  "type": "error",
  "message": "Error message",
  "widgetId": "widget_1"
}
```

## Rate Limiting

API endpoints are rate-limited to 60 requests per minute per client IP.

**Rate Limit Headers:**
- `X-RateLimit-Limit`: Maximum requests per window
- `X-RateLimit-Remaining`: Requests remaining in current window
- `X-RateLimit-Reset`: Unix timestamp when limit resets

**Status Code:** `429 Too Many Requests`

## Error Handling

### Error Response Format

```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

### Common Error Codes

- `INVALID_REQUEST`: Request parameters are invalid
- `NOT_FOUND`: Resource not found
- `INTERNAL_ERROR`: Server error
- `RATE_LIMIT_EXCEEDED`: Too many requests

## CORS

All endpoints support CORS. Requests from any origin are allowed.

## Examples

### Get System Metrics

```bash
curl http://localhost:8080/api/v1/metrics
```

### Search Widgets

```bash
curl "http://localhost:8080/api/v1/search?q=weather&type=weather&limit=10"
```

### WebSocket Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws');
ws.onmessage = (event) => {
  console.log(JSON.parse(event.data));
};
```

## Versioning

The current API version is v1. Future versions will be available under different version prefixes.

## Changelog

### v1.0.0
- Initial release
- WebSocket support
- Metrics API
- Rate limiting
- CORS support
