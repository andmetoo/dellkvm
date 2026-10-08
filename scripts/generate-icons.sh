#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
package_dir=$project_root/cmd/dellkvm
asset_dir=$package_dir/assets
scratch=$(mktemp -d)
trap 'rm -rf "$scratch"' 0

for tool in rsvg-convert python3 llvm-rc llvm-cvtres; do
  command -v "$tool" >/dev/null 2>&1 || { echo "Missing tool: $tool" >&2; exit 1; }
done

for size in 16 32 48 64 256; do
  rsvg-convert -w "$size" -h "$size" -o "$scratch/$size.png" "$asset_dir/icon.svg"
done
cp "$scratch/32.png" "$asset_dir/icon-32.png"
cp "$scratch/256.png" "$asset_dir/icon.png"

python3 - "$scratch" "$asset_dir/icon.ico" <<'PY'
from pathlib import Path
import struct
import sys

scratch, output = map(Path, sys.argv[1:])
sizes = (16, 32, 48, 64, 256)
images = [(size, (scratch / f'{size}.png').read_bytes()) for size in sizes]
offset = 6 + 16 * len(images)
data = bytearray(struct.pack('<HHH', 0, 1, len(images)))
for size, image in images:
    dimension = size if size < 256 else 0
    data.extend(struct.pack('<BBBBHHII', dimension, dimension, 0, 0, 1, 32, len(image), offset))
    offset += len(image)
for _, image in images:
    data.extend(image)
output.write_bytes(data)
PY

(cd "$package_dir" && llvm-rc /FO "$scratch/dellkvm.res" dellkvm.rc)
llvm-cvtres /MACHINE:X64 /TIMESTAMP:0 "/OUT:$package_dir/dellkvm_windows_amd64.syso" "$scratch/dellkvm.res"
llvm-cvtres /MACHINE:ARM64 /TIMESTAMP:0 "/OUT:$package_dir/dellkvm_windows_arm64.syso" "$scratch/dellkvm.res"
