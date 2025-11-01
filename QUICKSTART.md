# Glance Quick Start Guide

## 5-Minute Setup

### 1. Prerequisites
```bash
# Check Go version (need 1.19+)
go version

# Check SQLite
sqlite3 --version
```

### 2. Clone & Setup
```bash
git clone https://github.com/glance-project/glance.git
cd glance
go mod download
npm install
```

### 3. Initialize Database
```bash
sqlite3 glance.db < scripts/init.sql
```

### 4. Run Application
```bash
go run ./cmd/main.go
```

### 5. Open Browser
```
http://localhost:8080
```

## Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /` | Dashboard UI |
| `GET /api/v1/health` | Health check |
| `GET /api/v1/metrics` | System metrics |
| `GET /api/v1/search` | Search widgets |
| `GET /api/ws` | WebSocket |
| `GET /api/v1/activity` | Activity logs |

## API Examples

### Get Metrics
```bash
curl http://localhost:8080/api/v1/metrics
```

### Search Widgets
```bash
curl "http://localhost:8080/api/v1/search?q=weather&type=weather"
```

### WebSocket Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws');
ws.onmessage = (event) => {
  console.log(JSON.parse(event.data));
};
```

## Common Commands

### Build
```bash
# Development
go build -o bin/glance ./cmd/main.go

# Production
go build -ldflags="-s -w" -o bin/glance-prod ./cmd/main.go
```

### Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific test
go test -run TestHubRegister ./test/...
```

### Lint
```bash
go vet ./...
golint ./...
```

## Configuration

Create `.env` file:
```env
PORT=8080
DATABASE_URL=sqlite:glance.db
RATE_LIMIT_REQUESTS=60
CACHE_TTL=300
```

## Docker Quick Start

```bash
# Build image
docker build -t glance:latest .

# Run container
docker run -p 8080:8080 glance:latest

# Or with Docker Compose
docker-compose up -d
```

## Project Structure Quick Reference

```
glance/
├── cmd/main.go                    # Application entry
├── internal/
│   ├── api/                       # HTTP handlers
│   ├── websocket/                 # Real-time updates
│   ├── metrics/                   # Performance metrics
│   ├── database/                  # Data storage
│   ├── cache/                     # Caching layer
│   └── search/                    # Search feature
├── assets/js/                     # Frontend code
├── test/                          # Tests
└── docs/                          # Documentation
```

## Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080
kill -9 <PID>
```

### Database Issues
```bash
# Check database
sqlite3 glance.db ".tables"

# Reinitialize
rm glance.db
sqlite3 glance.db < scripts/init.sql
```

### Module Not Found
```bash
go mod tidy
go mod download
```

## Performance Tips

- Increase `RATE_LIMIT_REQUESTS` for higher throughput
- Adjust `CACHE_TTL` based on data freshness needs
- Use `WS_WRITE_BUFFER_SIZE` for large messages
- Monitor memory with: `go tool pprof http://localhost:8080/debug/pprof`

## Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/amazing-feature

# 2. Make changes
# ... edit files ...

# 3. Run tests
go test ./...

# 4. Run linter
go vet ./...

# 5. Commit changes
git commit -m "Add amazing feature"

# 6. Push and create PR
git push origin feature/amazing-feature
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific test file
go test -v ./test/websocket_hub_test.go

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Deployment Checklist

- [ ] Build production binary
- [ ] Set `ENV=production`
- [ ] Configure environment variables
- [ ] Initialize database
- [ ] Set up reverse proxy (nginx)
- [ ] Configure SSL/TLS
- [ ] Set up systemd service
- [ ] Enable monitoring
- [ ] Configure backups
- [ ] Test failover

## Next Steps

1. **Read Documentation**: `docs/API.md`, `docs/DEPLOYMENT.md`
2. **Explore Examples**: Check `test/e2e_test.go` for API usage
3. **Customize**: Add widgets and features as needed
4. **Deploy**: Follow `docs/DEPLOYMENT.md`

## Need Help?

- 📖 **API Docs**: `docs/API.md`
- 🚀 **Deployment**: `docs/DEPLOYMENT.md`
- 🤝 **Contributing**: `CONTRIBUTING.md`
- 📝 **Issues**: GitHub Issues
- 💬 **Discussions**: GitHub Discussions

## Key Features at a Glance

✨ **Real-time Updates**: WebSocket-powered live updates
📊 **Metrics**: Built-in performance monitoring
🔒 **Rate Limiting**: Protect your API
📱 **PWA**: Install as app, works offline
🔍 **Search**: Advanced filtering
📝 **Logging**: Comprehensive audit trail
⚡ **Fast**: <100ms response time

---

**Ready to dive deeper?** Check out the full documentation in the `docs/` folder.
