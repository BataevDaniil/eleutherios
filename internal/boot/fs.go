package boot

import (
	"context"
	"fmt"
	"os"

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
	logFileArg := ""
	if path := logging.LogFilePath(); path != "" {
		logFileArg = fmt.Sprintf(" --log-file %q", path)
	}
	data := fmt.Sprintf(`#!/bin/sh
[ "$1" = "start" ] || exit 0
exec %q --fs-hook%s
`, bin, logFileArg)
	if err := os.WriteFile(FSHook, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", FSHook, err)
	}
	logging.Logger().Info("fs hook установлен", "component", "hook", "hook", "fs", "file", FSHook)
	return nil
}
