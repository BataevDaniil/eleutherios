package iptables

import (
	"fmt"
	"os"
)

const hookFile = "/opt/etc/ndm/netfilter.d/100-eleutherios"

func InstallHook(iface string) error {
	bin, err := os.Executable()
	if err != nil {
		return fmt.Errorf("путь текущего бинарника: %w", err)
	}
	if err := os.MkdirAll("/opt/etc/ndm/netfilter.d", 0755); err != nil {
		return fmt.Errorf("создание netfilter.d: %w", err)
	}
	data := fmt.Sprintf(`#!/bin/sh
[ "$type" = "iptables" ] || exit 0
exec %q --iptables-hook --net %q
`, bin, iface)
	if err := os.WriteFile(hookFile, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", hookFile, err)
	}
	return nil
}

func RemoveHook() {
	_ = os.Remove(hookFile)
}
