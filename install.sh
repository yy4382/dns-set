#!/bin/bash

set -e

REPO="yy4382/dns-set"
BINARY_NAME="dns-set"
INSTALL_DIR_SYSTEM="/usr/local/bin"
INSTALL_DIR_USER="$HOME/.local/bin"

print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --version VERSION    Install specific version (default: latest)"
    echo "  --user               Install to ~/.local/bin instead of /usr/local/bin"
    echo "  --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                   # Install latest version system-wide"
    echo "  $0 --user            # Install latest version for current user"
    echo "  $0 --version v0.1.0  # Install specific version"
}

log_info() {
    echo "INFO: $1"
}

log_error() {
    echo "ERROR: $1" >&2
}

log_success() {
    echo "SUCCESS: $1"
}

detect_platform() {
    local os arch
    
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        darwin*)
            os="darwin"
            ;;
        linux*)
            os="linux"
            ;;
        mingw*|msys*|cygwin*)
            os="windows"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

get_latest_version() {
    log_info "Fetching latest release version..."
    local version
    version=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    
    if [ -z "$version" ]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi
    
    echo "$version"
}

download_binary() {
    local version="$1"
    local platform="$2"
    local binary_name="$BINARY_NAME"
    
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="${BINARY_NAME}.exe"
    fi
    
    local asset_name="${BINARY_NAME}-${platform}"
    if [[ "$platform" == *"windows"* ]]; then
        asset_name="${asset_name}.exe"
    fi
    
    local download_url="https://github.com/$REPO/releases/download/$version/$asset_name"
    local temp_file="/tmp/$asset_name"
    
    log_info "Downloading $asset_name from $version..."
    
    if ! curl -L -o "$temp_file" "$download_url"; then
        log_error "Failed to download binary from $download_url"
        exit 1
    fi
    
    echo "$temp_file"
}

verify_checksum() {
    local file="$1"
    local version="$2"
    local filename=$(basename "$file")
    
    log_info "Verifying checksum..."
    
    local checksums_url="https://github.com/$REPO/releases/download/$version/checksums.txt"
    local checksums_file="/tmp/checksums.txt"
    
    if curl -L -o "$checksums_file" "$checksums_url" 2>/dev/null; then
        local expected_checksum=$(grep "$filename" "$checksums_file" | awk '{print $1}')
        
        if [ -n "$expected_checksum" ]; then
            if command -v sha256sum >/dev/null 2>&1; then
                local actual_checksum=$(sha256sum "$file" | awk '{print $1}')
            elif command -v shasum >/dev/null 2>&1; then
                local actual_checksum=$(shasum -a 256 "$file" | awk '{print $1}')
            else
                log_info "No SHA256 utility found, skipping checksum verification"
                return 0
            fi
            
            if [ "$expected_checksum" = "$actual_checksum" ]; then
                log_info "Checksum verification passed"
            else
                log_error "Checksum verification failed"
                log_error "Expected: $expected_checksum"
                log_error "Actual:   $actual_checksum"
                exit 1
            fi
        else
            log_info "Checksum not found for $filename, skipping verification"
        fi
        
        rm -f "$checksums_file"
    else
        log_info "Could not download checksums, skipping verification"
    fi
}

install_binary() {
    local temp_file="$1"
    local install_dir="$2"
    local use_sudo="$3"
    
    local target_path="$install_dir/$BINARY_NAME"
    
    if [ ! -d "$install_dir" ]; then
        log_info "Creating directory $install_dir..."
        if [ "$use_sudo" = "true" ]; then
            sudo mkdir -p "$install_dir"
        else
            mkdir -p "$install_dir"
        fi
    fi
    
    log_info "Installing binary to $target_path..."
    
    if [ "$use_sudo" = "true" ]; then
        sudo cp "$temp_file" "$target_path"
        sudo chmod +x "$target_path"
    else
        cp "$temp_file" "$target_path"
        chmod +x "$target_path"
    fi
    
    rm -f "$temp_file"
    
    echo "$target_path"
}

check_installation() {
    local binary_path="$1"
    
    if [ -x "$binary_path" ]; then
        log_success "Installation completed successfully!"
        log_info "Binary installed at: $binary_path"
        
        if command -v "$BINARY_NAME" >/dev/null 2>&1; then
            local version_output
            version_output=$("$BINARY_NAME" --version 2>/dev/null || echo "Version info not available")
            log_info "Installed version: $version_output"
        else
            log_info "Note: $binary_path is not in your PATH"
            log_info "You may need to add $(dirname "$binary_path") to your PATH or use the full path"
        fi
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

main() {
    local version=""
    local user_install=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                version="$2"
                shift 2
                ;;
            --user)
                user_install=true
                shift
                ;;
            --help)
                print_usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                print_usage
                exit 1
                ;;
        esac
    done
    
    log_info "Starting dns-set installation..."
    
    local platform
    platform=$(detect_platform)
    log_info "Detected platform: $platform"
    
    if [ -z "$version" ]; then
        version=$(get_latest_version)
    fi
    log_info "Installing version: $version"
    
    local temp_file
    temp_file=$(download_binary "$version" "$platform")
    
    verify_checksum "$temp_file" "$version"
    
    local install_dir use_sudo
    if [ "$user_install" = "true" ]; then
        install_dir="$INSTALL_DIR_USER"
        use_sudo=false
        log_info "Installing for current user to $install_dir"
    else
        install_dir="$INSTALL_DIR_SYSTEM"
        use_sudo=true
        log_info "Installing system-wide to $install_dir (requires sudo)"
    fi
    
    local binary_path
    binary_path=$(install_binary "$temp_file" "$install_dir" "$use_sudo")
    
    check_installation "$binary_path"
}

if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi