# Deployment Guide

## Prerequisites

- Go 1.19+
- Node.js 18+ (for frontend development)
- SQLite 3
- Docker (optional)

## Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/glance-project/glance.git
cd glance
```

### 2. Install Dependencies

**Backend:**
```bash
go mod download
go mod tidy
```

**Frontend:**
```bash
npm install
```

### 3. Initialize Database

```bash
sqlite3 glance.db < scripts/init.sql
```

### 4. Build the Project

**Backend:**
```bash
go build -o bin/glance ./cmd/main.go
```

**Frontend (Optional, for production build):**
```bash
npm run build
```

### 5. Run the Application

```bash
./bin/glance
```

The application will start on `http://localhost:8080`

## Configuration

Create a `.env` file in the project root:

```env
# Server
PORT=8080
HOST=0.0.0.0
ENV=development

# Database
DATABASE_URL=sqlite:glance.db

# WebSocket
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024

# Metrics
METRICS_ENABLED=true
METRICS_INTERVAL=5000

# Caching
CACHE_TTL=300
CACHE_CLEANUP_INTERVAL=3600

# Rate Limiting
RATE_LIMIT_REQUESTS=60
RATE_LIMIT_WINDOW=60
```

## Docker Deployment

### Build Docker Image

```bash
docker build -t glance:latest .
```

### Run Docker Container

```bash
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e DATABASE_URL=sqlite:/app/data/glance.db \
  glance:latest
```

### Docker Compose

```bash
docker-compose up -d
```

See `docker-compose.yml` for configuration options.

## Production Deployment

### 1. Build Release Binary

```bash
go build -ldflags="-s -w" -o bin/glance-prod ./cmd/main.go
```

### 2. Set Environment to Production

```bash
export ENV=production
```

### 3. Configure Reverse Proxy (nginx)

```nginx
upstream glance {
    server localhost:8080;
}

server {
    listen 80;
    server_name yourdomain.com;

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    location / {
        proxy_pass http://glance;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api/ws {
        proxy_pass http://glance;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }
}
```

### 4. Set Up SSL/TLS

```bash
# Using Let's Encrypt with Certbot
sudo certbot certonly --standalone -d yourdomain.com
```

Update nginx configuration to use SSL:

```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    # ... rest of configuration
}

server {
    listen 80;
    return 301 https://$host$request_uri;
}
```

### 5. Set Up Process Manager

**Using systemd:**

Create `/etc/systemd/system/glance.service`:

```ini
[Unit]
Description=Glance Dashboard
After=network.target

[Service]
Type=simple
User=glance
WorkingDirectory=/opt/glance
ExecStart=/opt/glance/bin/glance-prod
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable glance
sudo systemctl start glance
```

### 6. Database Backup

```bash
# Daily backup
0 2 * * * /usr/bin/sqlite3 /opt/glance/glance.db ".backup /backups/glance-$(date +\%Y\%m\%d).db"
```

### 7. Monitoring

Install monitoring tools:

```bash
# Using Prometheus
wget https://github.com/prometheus/prometheus/releases/download/v2.40.0/prometheus-2.40.0.linux-amd64.tar.gz
tar xvfz prometheus-2.40.0.linux-amd64.tar.gz
```

Configure Prometheus scrape config:

```yaml
scrape_configs:
  - job_name: 'glance'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

## Troubleshooting

### Port Already in Use

```bash
lsof -i :8080
kill -9 <PID>
```

### Database Lock

```bash
# Close any open connections and restart
sqlite3 glance.db "PRAGMA integrity_check;"
```

### WebSocket Connection Issues

- Ensure WebSocket upgrade is allowed in reverse proxy
- Check firewall rules for port 8080
- Verify proxy configuration supports WebSocket

### Memory Issues

Monitor memory usage:
```bash
# Increase cache TTL to reduce memory usage
CACHE_TTL=600
```

## Performance Tuning

### Database

```sql
-- Optimize for common queries
ANALYZE;
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
```

### Application

- Increase `RATE_LIMIT_REQUESTS` for higher throughput
- Adjust `WS_READ_BUFFER_SIZE` and `WS_WRITE_BUFFER_SIZE` for your needs
- Set appropriate `CACHE_TTL` based on data freshness requirements

## Scaling

For horizontal scaling:

1. Use a shared database (PostgreSQL recommended)
2. Set up a load balancer (nginx, HAProxy)
3. Use a session store (Redis)
4. Implement distributed metrics collection

## Backup and Recovery

### Full Backup

```bash
tar -czf glance-backup-$(date +%Y%m%d).tar.gz data/
```

### Restore

```bash
tar -xzf glance-backup-20240115.tar.gz
```

## Support

For issues and questions, refer to the documentation or submit an issue on GitHub.
