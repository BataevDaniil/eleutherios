package boot

import (
	"fmt"
	"os"

	"github.com/BataevDaniil/eleutherios/internal/fsutil"
	"github.com/BataevDaniil/eleutherios/internal/logging"
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
	logging.Logger().Info("Автозапуск установлен", "component", "boot", "init_file", InitFile, "fs_hook", FSHook)
	return nil
}

func Remove() {
	for _, path := range []string{InitFile, FSHook} {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			logging.Logger().Warn("не удалось удалить файл", "component", "boot", "file", path, "error", err)
		}
	}
	logging.Logger().Info("Автозапуск удалён", "component", "boot", "init_file", InitFile, "fs_hook", FSHook)
}

func writeInit(bin, wgName, netName string) error {
	if err := os.MkdirAll("/opt/etc/init.d", 0755); err != nil {
		return fmt.Errorf("создание init.d: %w", err)
	}
	data := renderInit(bin, wgName, netName, logging.LogFilePath())
	if err := fsutil.WriteAtomic(InitFile, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", InitFile, err)
	}
	return nil
}

// renderInit генерирует тело init.d-скрипта. Все аргументы попадают в shell
// через %q — это безопасно квотирует пробелы и спецсимволы.
func renderInit(bin, wgName, netName, logPath string) string {
	logFileArg := ""
	if logPath != "" {
		logFileArg = fmt.Sprintf(" --log-file %q", logPath)
	}
	return fmt.Sprintf(`#!/bin/sh
case "$1" in
	start|restart)
		exec %q start --wg %q --net %q%s
	;;
esac
exit 0
`, bin, wgName, netName, logFileArg)
}
