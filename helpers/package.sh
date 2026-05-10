#!/bin/sh
# Упаковка в .ipk для Entware

set -e

APP="eleutherios"
VER="0.1.0"
ARCH="${1:-mipsel}"
PKG_DIR="build/${APP}_${VER}_${ARCH}"

mkdir -p "$PKG_DIR/opt/bin"
mkdir -p "$PKG_DIR/opt/etc"

cp "build/$APP-$ARCH" "$PKG_DIR/opt/bin/$APP"
chmod 755 "$PKG_DIR/opt/bin/$APP"

# control-файл для opkg
mkdir -p "$PKG_DIR/CONTROL"
cat > "$PKG_DIR/CONTROL/control" <<EOF
Package: $APP
Version: $VER
Architecture: $ARCH
Maintainer: eleutherios@zeleza.ru
Description: WireGuard split-tunnel: *.ru → ISP, остальное → WG
Depends: wireguard-go, dnsmasq-full, ipset, iptables
EOF

# упаковка
cd "$PKG_DIR"
tar czf ../data.tar.gz opt
tar czf ../control.tar.gz CONTROL
echo "2.0" > ../debian-binary
cd ..
ar r "${APP}_${VER}_${ARCH}.ipk" debian-binary control.tar.gz data.tar.gz

rm -f debian-binary control.tar.gz data.tar.gz
rm -rf "${APP}_${VER}_${ARCH}"

echo "Готово: build/${APP}_${VER}_${ARCH}.ipk"
