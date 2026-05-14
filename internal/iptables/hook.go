package iptables

import (
	"fmt"
	"os"

	"github.com/BataevDaniil/eleutherios/internal/logging"
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
	logFileArg := ""
	if path := logging.LogFilePath(); path != "" {
		logFileArg = fmt.Sprintf(" --log-file %q", path)
	}
	data := fmt.Sprintf(`#!/bin/sh
[ "$type" = "iptables" ] || exit 0
exec %q --iptables-hook --net %q%s
`, bin, iface, logFileArg)
	if err := os.WriteFile(hookFile, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", hookFile, err)
	}
	logging.Logger().Info("iptables hook установлен", "component", "hook", "hook", "iptables", "file", hookFile, "iface", iface)
	return nil
}

func RemoveHook() {
	if err := os.Remove(hookFile); err != nil && !os.IsNotExist(err) {
		logging.Logger().Warn("не удалось удалить hook файл", "component", "hook", "hook", "iptables", "file", hookFile, "error", err)
	}
	logging.Logger().Info("iptables hook удалён", "component", "hook", "hook", "iptables", "file", hookFile)
}
