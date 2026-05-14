package dns

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "обновить golden-файлы в testdata/")

func TestRenderOverlay(t *testing.T) {
	got := renderOverlay()
	assertGolden(t, "overlay.dnsmasq", got)
}

func TestRenderBaseConfig(t *testing.T) {
	got := renderBaseConfig(Port)
	assertGolden(t, "base.dnsmasq.conf", got)
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
		t.Errorf("несовпадение с golden %s\nПерезапустите тест с -update если изменение намеренное.\n--- got ---\n%s\n--- want ---\n%s", path, got, string(want))
	}
}
