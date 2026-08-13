#!/bin/sh
set -e

# xyops-integrations installer — downloads release binaries from GitHub.
# Usage: curl -fsSL https://github.com/gumieri/xyops-integrations/releases/latest/download/install.sh | sh
# Options:
#   -b, --binary-only   install binaries only (no systemd units / config)
#   -v, --version=X     install a specific version (default: latest)
#       --dry-run       show what would be installed without installing
#   -h, --help          show this help

VERSION=""
BINARY_ONLY=0
DRY_RUN=0
REPO="gumieri/xyops-integrations"
INSTALL_DIR="/usr/bin"
ETC_DIR="/etc/xyops"

while [ "$#" -gt 0 ]; do
    case "$1" in
        -b|--binary-only) BINARY_ONLY=1 ;;
        --dry-run) DRY_RUN=1 ;;
        -v|--version)
            shift
            VERSION="$1"
            ;;
        -v=*|--version=*) VERSION="${1#*=}" ;;
        -h|--help)
            grep '^#' "$0" | sed 's/^# \{0,1\}//'
            exit 0
            ;;
        *) echo "Unknown option: $1" >&2; exit 1 ;;
    esac
    shift
done

detect_os_arch() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
    esac
    case "$OS" in
        linux|darwin) ;;
        *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
    esac
}

resolve_latest() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' | head -n1 | sed 's/.*"v//; s/".*//'
}

install_archive() {
    name="$1"
    arch_name="$2"
    os_arch="${OS}_${ARCH}"
    archive="xyops-integrations_${arch_name}_${VERSION}_${os_arch}.tar.gz"
    url="https://github.com/${REPO}/releases/download/v${VERSION}/${archive}"

    echo "Downloading ${url}"
    curl -fsSL -o "/tmp/${archive}" "${url}"

    tmpdir="/tmp/xyops-integrations-${name}-${VERSION}"
    rm -rf "${tmpdir}"
    mkdir -p "${tmpdir}"
    tar -xzf "/tmp/${archive}" -C "${tmpdir}"
    echo "${tmpdir}"
}

detect_os_arch

if [ -z "$VERSION" ]; then
    VERSION="$(resolve_latest)"
fi
echo "Installing xyops-integrations v${VERSION} (${OS}/${ARCH})"

if [ "$DRY_RUN" = "1" ]; then
    echo "Dry run: would install xyops-mcp to ${INSTALL_DIR} and units to /usr/lib/systemd/system"
    exit 0
fi

# xyops-mcp ships the systemd unit inside its tarball (linux only).
dir="$(install_archive "mcp" "mcp")"
install -Dm755 "${dir}/xyops-mcp" "${INSTALL_DIR}/xyops-mcp"
if [ "$BINARY_ONLY" = "1" ] || [ "$OS" != "linux" ]; then
    rm -rf "${dir}"
    echo "Installed xyops-mcp to ${INSTALL_DIR}/xyops-mcp"
    exit 0
fi

install -Dm644 "${dir}/xyops-mcp.service" /usr/lib/systemd/system/xyops-mcp.service
mkdir -p "${ETC_DIR}"
cat > "${ETC_DIR}/opencode-config.json.example" <<'EOF'
{
  "mcp": {
    "servers": {
      "xyops": {
        "command": "/usr/bin/xyops-mcp --stdio",
        "env": {
          "XYOPS_BASE_URL": "https://xyops.example.com",
          "XYOPS_API_KEY": "${XYOPS_API_KEY}"
        }
      }
    }
  }
}
EOF
rm -rf "${dir}"

systemctl daemon-reload
systemctl enable --now xyops-mcp

echo "Installed xyops-mcp to ${INSTALL_DIR}/xyops-mcp"
echo "Set XYOPS_BASE_URL and XYOPS_API_KEY in ${ETC_DIR}/xyops-mcp.env and restart xyops-mcp"
