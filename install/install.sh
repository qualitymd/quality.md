#!/bin/sh
set -eu

repo="qualitymd/quality.md"
version="${QUALITYMD_VERSION:-latest}"
install_root="${QUALITYMD_HOME:-"$HOME/.qualitymd"}"
bin_dir="$install_root/bin"
non_interactive="${QUALITYMD_NO_INPUT:-0}"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      version="$2"
      shift 2
      ;;
    --install-dir)
      install_root="$2"
      bin_dir="$install_root/bin"
      shift 2
      ;;
    --yes|--non-interactive)
      non_interactive=1
      shift
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2
      ;;
  esac
done

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 70 ;;
esac
case "$os" in
  darwin) os="darwin" ;;
  linux) os="linux" ;;
  *) echo "unsupported operating system: $os" >&2; exit 70 ;;
esac

if [ "$version" = "latest" ]; then
  latest_json="$(curl -fsSL -H "User-Agent: qualitymd-installer" "https://api.github.com/repos/$repo/releases/latest" 2>/dev/null || true)"
  version="$(printf '%s\n' "$latest_json" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
  if [ -z "$version" ]; then
    latest_url="$(curl -fsIL -o /dev/null -w '%{url_effective}' "https://github.com/$repo/releases/latest" 2>/dev/null || true)"
    case "$latest_url" in
      */releases/tag/*) version="${latest_url##*/}" ;;
    esac
  fi
fi
if [ -z "$version" ]; then
  echo "could not resolve qualitymd version" >&2
  exit 70
fi

libc_suffix=""
if [ "$os" = "linux" ]; then
  if (ldd --version 2>&1 || true) | grep -qi musl || ls /lib/ld-musl-*.so.1 >/dev/null 2>&1; then
    libc_suffix="_musl"
  fi
fi
archive="qualitymd_${os}_${arch}${libc_suffix}.tar.gz"
base_url="https://github.com/$repo/releases/download/$version"
tmp="${TMPDIR:-/tmp}/qualitymd-install.$$"
stage="$install_root/releases/$version"
mkdir -p "$tmp" "$stage" "$bin_dir"
trap 'rm -rf "$tmp"' EXIT INT TERM

# Print the SHA-256 hex digest of "$1" using whichever tool is available, or
# return non-zero when none is. Stock Linux ships sha256sum, macOS ships shasum,
# and openssl is a common fallback; gating on shasum alone verified nothing on a
# typical Linux host.
sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  elif command -v openssl >/dev/null 2>&1; then
    openssl dgst -sha256 "$1" | awk '{print $NF}'
  else
    return 1
  fi
}

curl -fsSL "$base_url/$archive" -o "$tmp/$archive"
curl -fsSL "$base_url/checksums.txt" -o "$tmp/checksums.txt"
expected="$(grep " $archive\$" "$tmp/checksums.txt" | awk '{print $1}')"
if [ -z "$expected" ]; then
  echo "$archive is not listed in checksums.txt" >&2
  exit 70
fi
if ! actual="$(sha256_of "$tmp/$archive")"; then
  echo "no SHA-256 tool (sha256sum, shasum, or openssl) found" >&2
  exit 70
fi
if [ "$expected" != "$actual" ]; then
  echo "checksum mismatch for $archive" >&2
  exit 70
fi

tar -xzf "$tmp/$archive" -C "$stage"
binary="$stage/qualitymd"
if [ ! -x "$binary" ] && [ -x "$stage/bin/qualitymd" ]; then
  binary="$stage/bin/qualitymd"
fi
if [ ! -x "$binary" ]; then
  echo "archive did not contain an executable qualitymd binary" >&2
  exit 70
fi

ln -sfn "$binary" "$bin_dir/qualitymd"
cat > "$install_root/.qualitymd-managed-install" <<EOF
layoutVersion=1
version=$version
channel=github
EOF

"$bin_dir/qualitymd" --version >/dev/null

echo "Installed qualitymd $version to $bin_dir/qualitymd"

# Print the PATH line for interactive installs when the bin directory is not
# already reachable. We never edit shell profiles from a piped install.
case ":${PATH}:" in
  *":$bin_dir:"*) on_path=1 ;;
  *) on_path=0 ;;
esac
if [ "$non_interactive" != "1" ] && [ "$on_path" != "1" ]; then
  echo "Add $bin_dir to your PATH (for example in ~/.profile or your shell rc):"
  echo "  export PATH=\"$bin_dir:\$PATH\""
fi
