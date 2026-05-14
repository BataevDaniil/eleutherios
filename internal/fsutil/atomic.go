// Package fsutil содержит утилиты файловой системы, общие для нескольких пакетов.
package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteAtomic пишет data в path атомарно: создаёт временный файл рядом,
// сбрасывает на диск и переименовывает. Это закрывает окно, когда читатель
// (например, NDM hook-loader) может прочитать полу-записанный файл.
//
// Файл получает mode после chmod (umask на rename не влияет). Если ошибка
// случится на любом этапе, временный файл будет удалён.
func WriteAtomic(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("временный файл рядом с %s: %w", path, err)
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("запись во временный %s: %w", tmpName, err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("fsync %s: %w", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close %s: %w", tmpName, err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		cleanup()
		return fmt.Errorf("chmod %s: %w", tmpName, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		cleanup()
		return fmt.Errorf("rename %s -> %s: %w", tmpName, path, err)
	}
	return nil
}
