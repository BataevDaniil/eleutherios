#!/bin/sh
# Сборка eleutherios под все архитектуры Keenetic
# mips, mipsel, aarch64

set -e

APP="eleutherios"
SRC="$(dirname "$0")/.."

cd "$SRC"

for arch in mips mipsel aarch64; do
    echo "=== Сборка $APP для $arch ==="
    case $arch in
        mips)
            export GOOS=linux
            export GOARCH=mips
            export GOMIPS=softfloat
            ;;
        mipsel)
            export GOOS=linux
            export GOARCH=mipsle
            export GOMIPS=softfloat
            ;;
        aarch64)
            export GOOS=linux
            export GOARCH=arm64
            ;;
    esac

    go build -ldflags="-s -w" -o "build/$APP-$arch" .
    echo "  -> build/$APP-$arch готов"
done

echo "=== Готово ==="
ls -lh build/$APP-*
