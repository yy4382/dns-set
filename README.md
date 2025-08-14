# dns-set

English | [简体中文](README.zh-CN.md)

A command-line tool for automatically managing DNS records on DNS providers, designed to run on VPS servers.

## Features

- **Multiple domain sources**: Manually input domains or parse from Caddyfile with interactive selection
- **Flexible IP detection**: Choose from network interface detection, external API queries (ip.sb), or manual input
- **DNS provider support**: Cloudflare integration with API token authentication
- **Record types**: Supports both A (IPv4) and AAAA (IPv6) records with TTL auto
- **Proxy control**: Choose between DNS-only (grey cloud) or proxied (yellow cloud) status
- **Configuration management**: Settings saved to `~/.config/dns-set/` with environment variable overrides
- **Interactive CLI**: User-friendly command-line interface with planned TUI upgrade

## Use Cases

Perfect for VPS administrators who need to:
- Automatically update DNS records when server IP changes
- Manage multiple domains from a single server
- Keep DNS records in sync with Caddy reverse proxy configurations
- Automate DNS management in deployment scripts

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/yy4382/dns-set/master/install.sh | bash
```

Or for user-only installation:
```bash
curl -sSL https://raw.githubusercontent.com/yy4382/dns-set/master/install.sh | bash -s -- --user
```

Then you should have access to the `dns-set` cli command. If not command is not found, check if `/usr/local/bin` (global install) or `$HOME/local/bin` (user install) is in your system `$PATH`.

### Alternative: Go Install

```bash
go install github.com/yy4382/dns-set@latest
```

## Quick Start

1. **Set up Cloudflare API token**:
   ```bash
   export CLOUDFLARE_API_TOKEN="your-api-token-here"
   ```

2. **Run the tool**:
   ```bash
   dns-set
   ```

3. **Follow interactive prompts** to:
   - Choose domain source (manual input or Caddyfile)
   - Select IP detection method
   - Choose record types (A, AAAA, or both)
   - Select proxy status (DNS-only or proxied)
   - Choose which domains to update
   - Confirm DNS record changes

## Configuration

### Config File Location
- Linux/macOS: `~/.config/dns-set/config.yaml`
- Windows: `%APPDATA%/dns-set/config.yaml`
- Custom location: Set `DNS_SET_CONFIG_DIR` environment variable

### .env File Support
The tool automatically loads environment variables from `.env` files in the following locations (in order):
1. `.env` and `.env.local` in current working directory
2. `.env` and `.env.local` in config directory

Example `.env` file:
```bash
CLOUDFLARE_API_TOKEN=your-api-token-here
DNS_SET_CADDYFILE_PATH=/custom/path/Caddyfile
```

### Environment Variables
All configuration options can be overridden with environment variables:

- `CLOUDFLARE_API_TOKEN`: Cloudflare API token
- `DNS_SET_CADDYFILE_PATH`: Custom Caddyfile location
- `DNS_SET_CONFIG_DIR`: Custom config directory location (overrides default `~/.config/dns-set/`)

### Example Config File
```yaml
cloudflare:
  api_token: "your-token-here"
preferences:
  caddyfile_path: "/etc/caddy/Caddyfile"
  default_ttl: 300
```

## Cloudflare Setup

1. Go to [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
2. Create a custom token with:
   - **Permissions**: `Zone.DNS`
   - **Zone Resources**: Include all zones or specific zones you want to manage
3. Use the token in config file or `CLOUDFLARE_API_TOKEN` environment variable

## Development

### Building from source
```bash
git clone https://github.com/yy4382/dns-set.git
cd dns-set
go build -o dns-set ./cmd/dns-set
```

### Running tests
```bash
go test ./...
```

## Roadmap

- [x] Core DNS management functionality
- [x] Cloudflare provider integration
- [x] Interactive CLI interface
- [ ] Terminal UI (TUI) interface
- [ ] Additional DNS providers (planned)

## License

MIT License - see LICENSE file for details.
