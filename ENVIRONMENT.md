# StreamForge Environment Configuration

This document explains how to configure and deploy StreamForge using environment variables.

## Configuration Files

- `.env` - Your local environment configuration (created by you, should not be committed)
- `.env.example` - Template showing all available configuration variables (committed to repo)

## Quick Start

1. **Copy the example configuration:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your settings (optional):**
   ```bash
   nano .env
   ```

3. **Build and push Docker images:**
   ```bash
   bash build-and-push.sh
   ```

4. **Deploy with Docker Compose:**
   ```bash
   cd docker
   docker-compose up -d
   ```

## Available Variables

### Docker Registry
- `DOCKER_REGISTRY` - Docker Hub username/registry (default: `iamrahulgupta`)
- `DOCKER_IMAGE_NAME` - Image name (default: `github-project-streamforge`)

### Broker Configuration
- `BROKER_LISTEN_ADDR` - Address brokers bind to (default: `0.0.0.0`)
- `BROKER_DATA_DIR` - Where broker stores data (default: `/data/streamforge`)
- `BROKER_LOG_LEVEL` - Logging level: `trace`, `debug`, `info`, `warn`, `error` (default: `info`)

### Per-Broker Settings
- `BROKER_1_ID` / `BROKER_2_ID` / `BROKER_3_ID` - Broker identifiers
- `BROKER_X_PORT` - Port each broker listens on (X=1,2,3)
- `BROKER_X_METRICS_PORT` - Prometheus metrics port for each broker

### Admin Configuration
- `ADMIN_PORT` - Admin REST API port (default: `8080`)
- `ADMIN_LOG_LEVEL` - Admin logging level (default: `info`)

### Monitoring
- `PROMETHEUS_PORT` - Prometheus scraping interface (default: `9090`)
- `GRAFANA_ADMIN_USER` - Grafana login username (default: `admin`)
- `GRAFANA_ADMIN_PASSWORD` - Grafana login password (default: `admin123`)
- `GRAFANA_PORT` - Grafana web interface (default: `3000`)

### Network
- `NETWORK_DRIVER` - Docker network driver (default: `bridge`)

## Example Configurations

### Development Environment
```bash
BROKER_LOG_LEVEL=debug
ADMIN_LOG_LEVEL=debug
GRAFANA_ADMIN_PASSWORD=dev123
```

### Production Environment
```bash
DOCKER_REGISTRY=your-company-registry
DOCKER_IMAGE_NAME=streamforge-prod
BROKER_LOG_LEVEL=warn
GRAFANA_ADMIN_PASSWORD=$(openssl rand -base64 32)
```

### Custom Ports
```bash
BROKER_1_PORT=10092
BROKER_2_PORT=10093
BROKER_3_PORT=10094
ADMIN_PORT=9000
PROMETHEUS_PORT=8090
GRAFANA_PORT=3001
```

## How build-and-push.sh Uses Environment Variables

The `build-and-push.sh` script automatically:
1. Sources the `.env` file: `source .env`
2. Uses variables for Docker registry: `${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}`
3. Generates timestamp tags: `YYYYMMDDHHMMSS` format
4. Pushes images to: `${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${TAG}`

Example execution:
```bash
# With .env sourced, build-and-push.sh uses:
# DOCKER_REGISTRY=iamrahulgupta
# DOCKER_IMAGE_NAME=github-project-streamforge
# Resulting images: iamrahulgupta/github-project-streamforge:20260303212918
```

## How docker-compose.yml Uses Environment Variables

The `docker-compose.yml` file references all variables using `${VAR_NAME:-default}` syntax:
- If variable is set in `.env`, it uses that value
- If variable is not set, it uses the default value specified after `:-`

Example service configuration:
```yaml
environment:
  BROKER_ID: ${BROKER_1_ID:-1}
  PORT: ${BROKER_1_PORT:-9092}
  LOG_LEVEL: ${BROKER_LOG_LEVEL:-info}
```

## Security Best Practices

1. **Never commit `.env`** - Add to `.gitignore`:
   ```bash
   echo ".env" >> .gitignore
   ```

2. **Use `.env.example` for documentation** - Commit this file with default/placeholder values

3. **Never hardcode secrets** - Keep passwords in `.env` only:
   ```bash
   # Bad: GRAFANA_ADMIN_PASSWORD=admin123 in scripts
   # Good: GRAFANA_ADMIN_PASSWORD=admin123 in .env only
   ```

4. **Rotate passwords in production** - Generate strong passwords:
   ```bash
   GRAFANA_ADMIN_PASSWORD=$(openssl rand -base64 32)
   ```

5. **Different .env per environment** - Use separate files:
   ```bash
   .env.dev
   .env.staging
   .env.prod
   ```
   Load with: `source .env.prod`

## Accessing Services

After deployment with `docker-compose up -d`:

### Admin API
- URL: `http://localhost:${ADMIN_PORT}`
- Example: `http://localhost:8080`

### Prometheus Metrics
- URL: `http://localhost:${PROMETHEUS_PORT}`
- Example: `http://localhost:9090`
- View broker metrics: `/graph`

### Grafana Dashboards
- URL: `http://localhost:${GRAFANA_PORT}`
- Example: `http://localhost:3000`
- Login: `${GRAFANA_ADMIN_USER}` / `${GRAFANA_ADMIN_PASSWORD}`

### Broker Ports
```
Broker 1: localhost:${BROKER_1_PORT} (default 9092)
Broker 2: localhost:${BROKER_2_PORT} (default 9093)
Broker 3: localhost:${BROKER_3_PORT} (default 9094)
```

## Troubleshooting

### Variables not loading in docker-compose
```bash
# Ensure .env is in docker/ directory or specify explicitly:
docker-compose --env-file ../.env -f docker-compose.yml up -d
```

### Variables not loading in build-and-push.sh
```bash
# Make sure you're in project root:
cd /path/to/project-streamforge
bash build-and-push.sh
```

### Check loaded variables
```bash
# View all variables that will be used:
set -a; source .env; set +a; env | grep -E "DOCKER_|BROKER_|ADMIN_|PROMETHEUS_|GRAFANA_"
```

## Migration from Hardcoded Values

If you previously had hardcoded values in scripts:

1. **Identify all hardcoded values:**
   ```bash
   grep -r "iamrahulgupta" . --include="*.sh" --include="*.yml"
   ```

2. **Add to .env:**
   ```bash
   DOCKER_REGISTRY=iamrahulgupta
   DOCKER_IMAGE_NAME=github-project-streamforge
   ```

3. **Update scripts/configs to use `${VAR_NAME}`:**
   ```bash
   # Before:
   docker build -t iamrahulgupta/github-project-streamforge:tag .
   
   # After:
   docker build -t ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:tag .
   ```

4. **Source .env before running scripts:**
   ```bash
   source .env
   bash build-and-push.sh
   ```

## Additional Resources

- [Docker Compose Environment Variables](https://docs.docker.com/compose/environment-variables/)
- [Twelve-Factor App Configuration](https://12factor.net/config)
- [StreamForge README](./README.md)
- [StreamForge Setup Guide](./DEVELOPER.md)
