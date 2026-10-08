#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
fixture=$(mktemp -d)
trap 'rm -rf "$fixture"' 0
mkdir -p "$fixture/home" "$fixture/release" "$fixture/fakebin"

case "$(uname -m)" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo 'Unsupported test architecture' >&2; exit 1 ;;
esac
archive=dellkvm_1.2.3_linux_$arch.tar.gz
printf '#!/bin/sh\necho fixture\n' > "$fixture/release/dellkvm"
tar -czf "$fixture/release/$archive" -C "$fixture/release" dellkvm
(cd "$fixture/release" && sha256sum "$archive" > checksums.txt)

cat > "$fixture/fakebin/curl" <<'EOF'
#!/bin/sh
destination=
last=
while [ "$#" -gt 0 ]; do
  case "$1" in
    -o) destination=$2; shift 2 ;;
    *) last=$1; shift ;;
  esac
done
case "$last" in
  */releases/latest)
    printf '%s' 'https://github.com/andmetoo/dellkvm/releases/tag/v1.2.3'
    exit 0
    ;;
esac
cp "$INSTALL_TEST_RELEASE/${last##*/}" "$destination"
EOF
chmod +x "$fixture/fakebin/curl"

export HOME="$fixture/home"
export XDG_DATA_HOME="$fixture/home/data"
export XDG_CONFIG_HOME="$fixture/home/config"
export INSTALL_TEST_RELEASE="$fixture/release"
PATH="$fixture/fakebin:$PATH"
export PATH

sh "$root/install.sh"
test -x "$HOME/.local/bin/dellkvm"
grep -Fqx "Exec=\"$HOME/.local/bin/dellkvm\" tray" "$XDG_CONFIG_HOME/autostart/dellkvm.desktop"
test -f "$XDG_DATA_HOME/applications/dellkvm.desktop"

printf 'my config\n' > "$XDG_CONFIG_HOME/config.toml"
printf 'bad checksum  %s\n' "$archive" > "$fixture/release/checksums.txt"
if sh "$root/install.sh" --version v1.2.3 >/dev/null 2>&1; then
  echo 'Accepted an invalid checksum' >&2
  exit 1
fi
test "$("$HOME/.local/bin/dellkvm")" = fixture

(cd "$fixture/release" && sha256sum "$archive" > checksums.txt)
sh "$root/install.sh" --version v1.2.3 --no-autostart
test ! -e "$XDG_CONFIG_HOME/autostart/dellkvm.desktop"
sh "$root/install.sh" --uninstall
test ! -e "$HOME/.local/bin/dellkvm"
test ! -e "$XDG_DATA_HOME/applications/dellkvm.desktop"
test -f "$XDG_CONFIG_HOME/config.toml"
echo 'Linux installer checks passed'
