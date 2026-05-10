package boot

import (
	"fmt"
	"os"

	"github.com/BataevDaniil/eleutherios/ipset"
)

func RunFSHook() error {
	return ipset.CreateSets()
}

func writeFSHook(bin string) error {
	if err := os.MkdirAll("/opt/etc/ndm/fs.d", 0755); err != nil {
		return fmt.Errorf("создание fs.d: %w", err)
	}
	data := fmt.Sprintf(`#!/bin/sh
[ "$1" = "start" ] || exit 0
exec %q --fs-hook
`, bin)
	if err := os.WriteFile(FSHook, []byte(data), 0755); err != nil {
		return fmt.Errorf("запись %s: %w", FSHook, err)
	}
	return nil
}
