#!/bin/sh
[ "$1" = "start" ] || exit 0
exec "/opt/root/eleutherios" --fs-hook
