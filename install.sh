#!/bin/bash
set -euo pipefail

REPO="shunsukeaihara/nanobanana_cli"
INSTALL_DIR="${HOME}/.local/bin"
SKILL_DIR="${HOME}/.claude/skills/nano-banana-pro"

# Detect OS and architecture
detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"

  case "$os" in
    linux)  os="linux" ;;
    darwin) os="darwin" ;;
    *)      echo "Unsupported OS: $os" >&2; exit 1 ;;
  esac

  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *)             echo "Unsupported architecture: $arch" >&2; exit 1 ;;
  esac

  echo "${os}_${arch}"
}

# Get latest release tag from GitHub API
get_latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

main() {
  echo "==> Detecting platform..."
  local platform
  platform="$(detect_platform)"
  echo "    Platform: ${platform}"

  echo "==> Fetching latest version..."
  local version
  version="$(get_latest_version)"
  echo "    Version: ${version}"

  local archive_url="https://github.com/${REPO}/releases/download/${version}/nanobanana_${platform}.tar.gz"

  echo "==> Downloading nanobanana..."
  local tmpdir
  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT

  curl -fsSL "$archive_url" | tar xz -C "$tmpdir"

  echo "==> Installing to ${INSTALL_DIR}/nanobanana..."
  mkdir -p "$INSTALL_DIR"
  mv "$tmpdir/nanobanana" "$INSTALL_DIR/nanobanana"
  chmod +x "$INSTALL_DIR/nanobanana"

  echo "==> Installing Claude Code skill..."
  mkdir -p "$SKILL_DIR"
  curl -fsSL "https://raw.githubusercontent.com/${REPO}/main/skill/SKILL.md" \
    -o "$SKILL_DIR/SKILL.md"

  echo ""
  echo "Installation complete!"
  echo ""
  echo "  Binary: ${INSTALL_DIR}/nanobanana"
  echo "  Skill:  ${SKILL_DIR}/SKILL.md"
  echo ""

  # Check if INSTALL_DIR is in PATH
  if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
    echo "NOTE: ${INSTALL_DIR} is not in your PATH."
    echo "Add it with:"
    echo ""
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    echo ""
  fi

  echo "Set your API key:"
  echo ""
  echo "  export GEMINI_API_KEY=\"your-api-key\""
  echo ""
  echo "Then use it:"
  echo ""
  echo "  nanobanana \"A cat playing piano in watercolor style\""
}

main
