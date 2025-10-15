# Cloudflare DDNS Updater
A lightweight, self-hosted Dynamic DNS (DDNS) service that automatically updates your Cloudflare DNS records with your current public IP address.
Features
- 🔄 Automatic IP detection and DNS record updates
- ⏰ Configurable update intervals using cron expressions
- 🐳 Docker support with multi-stage builds
- 🔒 Secure API token-based authentication
- 📦 Minimal footprint (~10MB Docker image)
- 🚀 Runs as non-root user for enhanced security
- ✨ Immediate update on startup + scheduled updates

## Prerequisites
- A Cloudflare account with DNS records to manage
- Cloudflare API Token with DNS edit permissions
- Docker and Docker Compose (for containerized deployment)
- Go 1.24+ (for local development)

## Configuration
The service is configured via environment variables:

| Variable     | Required | Default    | Description                                    |
|--------------|----------|------------|------------------------------------------------|
| CF_API_TOKEN | ✅ Yes    | -          | Cloudflare API Token with DNS edit permissions |
| CF_ZONE_ID   | ✅ Yes    | -          | Cloudflare Zone ID for your domain             |
| CRON_SPEC    | ❌ No     | @every 30m | Cron expression for the update interval        |

## Getting Cloudflare Credentials
1. API Token: Go to Cloudflare Dashboard → Create Token → Edit zone DNS template
2. Zone ID: Found in your domain's Overview page under "API" section

## Deployment

### Using Pre-built Image from GitHub Container Registry (Easiest)
Pre-built multi-platform Docker images are automatically published to GitHub Container Registry on every release and push to main.

Pull and run the latest image:
```bash
docker pull ghcr.io/nico-iaco/cloudflarednsupdater:latest

docker run -d \
  --name cf-ddns \
  --restart unless-stopped \
  -e CF_API_TOKEN=your_token \
  -e CF_ZONE_ID=your_zone_id \
  -e CRON_SPEC="@every 30m" \
  ghcr.io/nico-iaco/cloudflarednsupdater:latest
```

Or use with Docker Compose:
```yaml
services:
  cf-ddns:
    image: ghcr.io/nico-iaco/cloudflarednsupdater:latest
    environment:
      CF_API_TOKEN: "${CF_API_TOKEN}"
      CF_ZONE_ID: "${CF_ZONE_ID}"
      CRON_SPEC: "@every 30m"
    restart: unless-stopped
```

Available tags:
- `latest` - Latest build from main branch
- `v1.0.0` - Specific version tags
- `1.0` - Major.minor version tags
- `1` - Major version tags
- `sha-<commit>` - Specific commit builds

### Docker Compose (Build from Source)
1. Clone the repository:
```bash
git clone <repository-url>
cd cloudflareDnsUpdater
```
2. Create a .env file:
```bash
CF_API_TOKEN=your_cloudflare_api_token_here
CF_ZONE_ID=your_zone_id_here
CRON_SPEC=@every 30m
```
3. Start the service:
```bash
docker-compose up -d
```
4. View logs:
```bash
docker-compose logs -f
```

### Standalone Docker
```bash
docker build -t cf-ddns:latest .

docker run -d \
--name cf-ddns \
--restart unless-stopped \
-e CF_API_TOKEN=your_token \
-e CF_ZONE_ID=your_zone_id \
-e CRON_SPEC="@every 30m" \
cf-ddns:latest
### Local Development
# Install dependencies
go mod download

# Set environment variables
export CF_API_TOKEN=your_token
export CF_ZONE_ID=your_zone_id
export CRON_SPEC="@every 5m"

# Run
go run main.go
```
## Cron Schedule Examples
The CRON_SPEC variable accepts standard cron expressions or predefined schedules:
- `@every 30m` - Every 30 minutes (default)
- `@every 1h` - Every hour
- `@hourly - Every` hour
- `0 */6 * * *` - Every 6 hours
- `0 0 * * *` - Daily at midnight

## How It Works
1. The service starts and immediately checks your current public IP
2. It retrieves all A records from your Cloudflare zone
3. For each record, it compares the configured IP with your current public IP
4. If they differ, the DNS record is automatically updated
5. The process repeats on the configured schedule

## Logging
The service provides clear, emoji-enhanced logging:
- 🚀 Job start
- ℹ️ Information messages
- ✅ Success confirmations
- ❌ Error messages
- 🔄 Update operations
- 🎉 Successful updates
- ✨ Job completion

## Security
- Runs as non-root user (app:app) in Docker
- Uses multi-stage builds to minimize attack surface
- API tokens are never logged or exposed
- HTTPS is used for all external API calls

## Troubleshooting
### Service fails to start
- Verify CF_API_TOKEN and CF_ZONE_ID are correctly set
- Check API token permissions include DNS edit access
### IP not updating
- Check logs for error messages
- Verify network connectivity
- Ensure Cloudflare API is accessible
### Invalid cron expression
- Test your cron expression syntax
- Use crontab.guru for validation

## CI/CD and Container Registry

### Automated Builds
This project uses GitHub Actions to automatically build and publish Docker images to the GitHub Container Registry (ghcr.io).

**Workflow triggers:**
- **Push to main branch**: Builds and tags image as `latest`
- **Release creation**: Builds and tags with version numbers (e.g., `v1.0.0`, `1.0`, `1`)
- **Manual dispatch**: Can be triggered manually from the Actions tab

**Multi-platform support:**
- The workflow builds images for both `linux/amd64` and `linux/arm64` platforms
- Images are automatically pushed to `ghcr.io/nico-iaco/cloudflarednsupdater`

**Image tags:**
- `latest` - Always points to the most recent build from main
- `v1.0.0` - Exact version tags for releases
- `1.0` and `1` - Convenience tags for major/minor versions
- `sha-<commit>` - Specific commit identifiers for reproducibility

**No additional secrets required:**
- The workflow uses the built-in `GITHUB_TOKEN` for authentication
- No manual configuration needed - works out of the box

**Viewing published images:**
Visit the [packages page](https://github.com/nico-iaco/cloudflareDnsUpdater/pkgs/container/cloudflarednsupdater) to see all available versions.

## License
This project is provided as-is for personal and commercial use.

## Contributing
Contributions are welcome! Please feel free to submit issues or pull requests.