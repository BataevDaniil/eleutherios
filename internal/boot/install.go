package boot

import (
	"fmt"
	"os"
)

func Install(wgName, netName string) error {
	bin, err := os.Executable()
	if err != nil {
		return fmt.Errorf("путь текущего бинарника: %w", err)
	}
	if err := writeInit(bin, wgName, netName); err != nil {
		return err
	}
	if err := writeFSHook(bin); err != nil {
		return err
	}
	fmt.Println("  автозапуск: init.d и fs-hook установлены")
	return nil
}

func Remove() {
	_ = os.Remove(InitFile)
	_ = os.Remove(FSHook)
}

func writeInit(bin, wgName, netName string) error {
	if err := os.MkdirAll("/opt/etc/init.d", 0755); err != nil {
		return fmt.Errorf("создание init.d: %w", err)
	}
	data := fmt.Sprintf(`#!/bin/sh
case "$1" in
	start|restart)
		exec %q start --wg %q --net %q
	;;
esac
exit 0
`, bin, wgName, netName)
	if err := os.WriteFile(InitFile, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", InitFile, err)
	}
	return nil
}
