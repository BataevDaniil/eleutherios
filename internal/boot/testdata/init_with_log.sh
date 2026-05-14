#!/bin/sh
case "$1" in
	start|restart)
		exec "/opt/root/eleutherios" start --wg "Wireguard0" --net "br0" --log-file "/tmp/eleutherios.log"
	;;
esac
exit 0
