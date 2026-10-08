#!/bin/sh
set -eu

repo=https://github.com/andmetoo/dellkvm
version=latest
autostart=yes
uninstall=no

usage() {
  cat <<'EOF'
Usage: sh install.sh [--version vX.Y.Z] [--no-autostart] [--uninstall]

Installs dellkvm for the current Linux user in ~/.local/bin. Autostart is
enabled by default. No sudo is used and the monitor configuration is preserved.
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      [ "$#" -ge 2 ] || { usage >&2; exit 2; }
      version=$2
      shift 2
      ;;
    --no-autostart)
      autostart=no
      shift
      ;;
    --uninstall)
      uninstall=yes
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      usage >&2
      exit 2
      ;;
  esac
done

[ "$(uname -s)" = Linux ] || { echo 'Linux is required' >&2; exit 1; }
bin_dir=$HOME/.local/bin
data_dir=${XDG_DATA_HOME:-$HOME/.local/share}/applications
autostart_dir=${XDG_CONFIG_HOME:-$HOME/.config}/autostart
desktop_file=$data_dir/dellkvm.desktop
autostart_file=$autostart_dir/dellkvm.desktop

remove_managed_desktop() {
  if [ -f "$1" ] && grep -qx 'X-dellkvm-managed=true' "$1"; then
    rm -f "$1"
  fi
}

if [ "$uninstall" = yes ]; then
  rm -f "$bin_dir/dellkvm"
  remove_managed_desktop "$desktop_file"
  remove_managed_desktop "$autostart_file"
  echo 'dellkvm removed; your config.toml was kept'
  exit 0
fi

case "$(uname -m)" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo 'Unsupported Linux architecture' >&2; exit 1 ;;
esac

for tool in curl sha256sum tar install awk sed grep mktemp; do
  command -v "$tool" >/dev/null 2>&1 || { echo "Missing dependency: $tool" >&2; exit 1; }
done

if [ "$version" = latest ]; then
  latest_url=$(curl -fsSL -o /dev/null -w '%{url_effective}' "$repo/releases/latest")
  version=${latest_url##*/}
fi
printf '%s\n' "$version" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$' || {
  echo "Invalid release version: $version" >&2
  exit 1
}

archive=dellkvm_${version#v}_linux_$arch.tar.gz
release_url=$repo/releases/download/$version
tmp=$(mktemp -d)
staged=
cleanup() {
  [ -z "$staged" ] || rm -f "$staged"
  rm -rf "$tmp"
}
trap cleanup 0
trap 'exit 1' 1 2 3 15

curl -fsSL --retry 3 -o "$tmp/$archive" "$release_url/$archive"
curl -fsSL --retry 3 -o "$tmp/checksums.txt" "$release_url/checksums.txt"
awk -v target="$archive" '$2 == target { print; found=1 } END { if (!found) exit 1 }' \
  "$tmp/checksums.txt" > "$tmp/archive.sha256" || {
    echo "No checksum found for $archive" >&2
    exit 1
  }
(cd "$tmp" && sha256sum -c archive.sha256)
tar -xzf "$tmp/$archive" -C "$tmp" dellkvm

mkdir -p "$bin_dir" "$data_dir"
staged=$bin_dir/.dellkvm.new.$$
install -m 0755 "$tmp/dellkvm" "$staged"
mv -f "$staged" "$bin_dir/dellkvm"
staged=

desktop_exec=$(printf '%s' "$bin_dir/dellkvm" | sed 's/\\/\\\\/g; s/"/\\"/g; s/`/\\`/g; s/\$/\\$/g')
cat > "$tmp/dellkvm.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=dellkvm
Comment=Switch monitor inputs from the system tray
Exec="$desktop_exec" tray
Icon=video-display
Terminal=false
Categories=Utility;HardwareSettings;
StartupNotify=false
X-dellkvm-managed=true
EOF
install -m 0644 "$tmp/dellkvm.desktop" "$desktop_file"

if [ "$autostart" = yes ]; then
  mkdir -p "$autostart_dir"
  install -m 0644 "$tmp/dellkvm.desktop" "$autostart_file"
else
  remove_managed_desktop "$autostart_file"
fi

echo "Installed dellkvm $version to $bin_dir/dellkvm"
if [ "$autostart" = yes ]; then
  echo 'Tray autostart is enabled for the next login'
else
  echo 'Tray autostart is disabled'
fi
if ! command -v ddcutil >/dev/null 2>&1; then
  echo 'Install ddcutil and grant access to /dev/i2c-* before using the monitor' >&2
fi
