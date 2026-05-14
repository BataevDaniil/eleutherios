#!/bin/sh
case "$1" in
	start|restart)
		exec "/opt/root/eleutherios" start --wg "Wireguard0" --net "br0"
	;;
esac
exit 0
