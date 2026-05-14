package boot

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "обновить golden-файлы в testdata/")

func TestRenderInit(t *testing.T) {
	got := renderInit("/opt/root/eleutherios", "Wireguard0", "br0", "")
	assertGolden(t, "init.sh", got)
}

func TestRenderInitWithLog(t *testing.T) {
	got := renderInit("/opt/root/eleutherios", "Wireguard0", "br0", "/tmp/eleutherios.log")
	assertGolden(t, "init_with_log.sh", got)
}

func TestRenderFSHook(t *testing.T) {
	got := renderFSHook("/opt/root/eleutherios", "")
	assertGolden(t, "fs_hook.sh", got)
}

func TestRenderFSHookWithLog(t *testing.T) {
	got := renderFSHook("/opt/root/eleutherios", "/tmp/eleutherios.log")
	assertGolden(t, "fs_hook_with_log.sh", got)
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.WriteFile(path, []byte(got), 0644); err != nil {
			t.Fatalf("обновить golden %s: %v", path, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("читать golden %s (запустите с -update для генерации): %v", path, err)
	}
	if string(want) != got {
		t.Errorf("несовпадение с golden %s\n--- got ---\n%s\n--- want ---\n%s", path, got, string(want))
	}
}
