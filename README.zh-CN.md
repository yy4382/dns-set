# dns-set

[English](README.md) | 简体中文

一个用于在 DNS 服务商上自动管理 DNS 记录的命令行工具，适合运行在 VPS 服务器上。

## 功能特性

- **多来源域名**：手动输入域名，或从 Caddyfile 解析并交互式选择
- **灵活的 IP 检测**：支持从网络接口探测、外部 API（ip.sb）查询、或手动输入
- **DNS 服务商支持**：内置 Cloudflare，使用 API Token 认证
- **记录类型**：支持 A（IPv4）与 AAAA（IPv6），TTL 自动
- **代理开关**：可选择仅 DNS（灰云）或代理（黄云）
- **配置管理**：设置保存至 `~/.config/dns-set/`，并支持环境变量覆盖
- **交互式 CLI**：友好的命令行交互界面（计划升级为 TUI）

## 适用场景

非常适合需要以下能力的 VPS 管理员：
- 当服务器 IP 变化时自动更新 DNS 记录
- 在单台服务器上管理多个域名
- 将 DNS 记录与 Caddy 反向代理配置保持同步
- 在部署脚本中自动化 DNS 管理

## 安装

```bash
go install github.com/yy4382/dns-set@latest
```

## 快速开始

1. **设置 Cloudflare API Token**：
   ```bash
   export CLOUDFLARE_API_TOKEN="your-api-token-here"
   ```

2. **运行工具**：
   ```bash
   dns-set
   ```

3. **按交互提示操作**，包括：
   - 选择域名来源（手动输入或 Caddyfile）
   - 选择 IP 检测方式
   - 选择记录类型（A、AAAA 或两者）
   - 选择代理状态（仅 DNS 或 代理）
   - 选择要更新的域名
   - 确认 DNS 记录变更

## 配置

### 配置文件位置
- Linux/macOS：`~/.config/dns-set/config.yaml`
- Windows：`%APPDATA%/dns-set/config.yaml`
- 自定义位置：设置环境变量 `DNS_SET_CONFIG_DIR`

### 支持 .env 文件
工具会按照以下顺序自动加载 `.env` 文件中的环境变量：
1. 当前工作目录下的 `.env` 与 `.env.local`
2. 配置目录下的 `.env` 与 `.env.local`

`.env` 示例：
```bash
CLOUDFLARE_API_TOKEN=your-api-token-here
DNS_SET_CADDYFILE_PATH=/custom/path/Caddyfile
```

### 环境变量
所有配置项均可通过环境变量覆盖：

- `CLOUDFLARE_API_TOKEN`：Cloudflare API Token
- `DNS_SET_CADDYFILE_PATH`：自定义 Caddyfile 路径
- `DNS_SET_CONFIG_DIR`：自定义配置目录（覆盖默认 `~/.config/dns-set/`）

### 配置文件示例
```yaml
cloudflare:
  api_token: "your-token-here"
preferences:
  caddyfile_path: "/etc/caddy/Caddyfile"
  default_ttl: 300
```

## Cloudflare 配置

1. 打开 [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
2. 创建自定义 Token，包含：
   - **Permissions**：`Zone.DNS`
   - **Zone Resources**：包含全部或需要管理的 Zone
3. 在配置文件或环境变量 `CLOUDFLARE_API_TOKEN` 中使用该 Token

## 开发

### 从源码构建
```bash
git clone https://github.com/yy4382/dns-set.git
cd dns-set
go build -o dns-set ./cmd/dns-set
```

### 运行测试
```bash
go test ./...
```

## 路线图

- [x] 核心 DNS 管理功能
- [x] Cloudflare 提供商集成
- [x] 交互式 CLI 界面
- [ ] 终端 UI（TUI）
- [ ] 更多 DNS 服务商（规划中）

## 许可协议

MIT License - 详见 LICENSE 文件。


