package boot

import (
	"context"
	"fmt"
	"os"

	"github.com/BataevDaniil/eleutherios/internal/fsutil"
	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/logging"
)

func RunFSHook(ctx context.Context) error {
	logging.Logger().Info("fs hook запущен", "component", "hook", "hook", "fs")
	if err := ipset.CreateSets(ctx); err != nil {
		return err
	}
	logging.Logger().Info("fs hook выполнен", "component", "hook", "hook", "fs")
	return nil
}

func writeFSHook(bin string) error {
	if err := os.MkdirAll("/opt/etc/ndm/fs.d", 0755); err != nil {
		return fmt.Errorf("создание fs.d: %w", err)
	}
	data := renderFSHook(bin, logging.LogFilePath())
	if err := fsutil.WriteAtomic(FSHook, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", FSHook, err)
	}
	logging.Logger().Info("fs hook установлен", "component", "hook", "hook", "fs", "file", FSHook)
	return nil
}

// renderFSHook генерирует тело NDM fs-hook, который вызывает eleutherios --fs-hook
// при старте файловой системы opt-раздела.
func renderFSHook(bin, logPath string) string {
	logFileArg := ""
	if logPath != "" {
		logFileArg = fmt.Sprintf(" --log-file %q", logPath)
	}
	return fmt.Sprintf(`#!/bin/sh
[ "$1" = "start" ] || exit 0
exec %q --fs-hook%s
`, bin, logFileArg)
}
