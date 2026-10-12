# SubNode Docker Compose

## Development

```bash
# Start all services
cd compose
cp .env.example .env
docker-compose up --build

# Or use Makefile
make dev
```

## Production

```bash
# Create production environment
cd compose
cp .env.example .env
# Edit .env with your configuration

# Start production services
docker-compose -f docker-compose.prod.yaml up -d

# Or use Makefile
make dev
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| POSTGRES_USER | Database user | subnode |
| POSTGRES_PASSWORD | Database password | subnode |
| POSTGRES_DB | Database name | subnode |
| JWT_SECRET | JWT signing secret | change-this-in-production |
| JWT_EXPIRY | JWT token expiry (minutes) | 60 |
| CORS_ALLOWED_ORIGINS | Allowed CORS origins | * |
| RATE_LIMIT_ENABLED | Enable rate limiting | true |
| RATE_LIMIT_REQUESTS | Max requests per window | 100 |
| RATE_LIMIT_WINDOW | Rate limit window (minutes) | 15 |

## Services

- **postgres**: PostgreSQL database
- **backend**: Go backend API
- **frontend**: React frontend (Nginx)

## Access

- Frontend: http://localhost
- Backend API: http://localhost:8080
- Health Check: http://localhost:8080/health

## Data Persistence

PostgreSQL data is stored in a Docker volume named `postgres-data`.

## Cleanup

```bash
docker-compose down
docker-compose down -v  # Remove volumes
```
