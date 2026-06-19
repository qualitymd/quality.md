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
  version="$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
fi
if [ -z "$version" ]; then
  echo "could not resolve qualitymd version" >&2
  exit 70
fi

archive="qualitymd_${os}_${arch}.tar.gz"
base_url="https://github.com/$repo/releases/download/$version"
tmp="${TMPDIR:-/tmp}/qualitymd-install.$$"
stage="$install_root/releases/$version"
mkdir -p "$tmp" "$stage" "$bin_dir"
trap 'rm -rf "$tmp"' EXIT INT TERM

curl -fsSL "$base_url/$archive" -o "$tmp/$archive"
if curl -fsSL "$base_url/checksums.txt" -o "$tmp/checksums.txt"; then
  if command -v shasum >/dev/null 2>&1; then
    expected="$(grep " $archive\$" "$tmp/checksums.txt" | awk '{print $1}')"
    actual="$(shasum -a 256 "$tmp/$archive" | awk '{print $1}')"
    if [ -n "$expected" ] && [ "$expected" != "$actual" ]; then
      echo "checksum mismatch for $archive" >&2
      exit 70
    fi
  fi
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

if [ "$non_interactive" != "1" ]; then
  echo "Installed qualitymd $version to $bin_dir/qualitymd"
  echo "Add $bin_dir to PATH if qualitymd is not already visible."
else
  echo "Installed qualitymd $version to $bin_dir/qualitymd"
fi
