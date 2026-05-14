package eleutherios

import (
	"regexp"
	"strings"
	"testing"
)

func TestIndent(t *testing.T) {
	tests := []struct {
		in     string
		prefix string
		want   string
	}{
		{"", "  ", "  (нет данных)"},
		{"line1\nline2", "  ", "  line1\n  line2"},
		{"single", "  ", "  single"},
	}
	for _, tt := range tests {
		got := indent(tt.in, tt.prefix)
		if got != tt.want {
			t.Errorf("indent(%q, %q) = %q, want %q", tt.in, tt.prefix, got, tt.want)
		}
	}
}

func TestFilter(t *testing.T) {
	input := "alpha\nbeta\ngamma\nalpha2"
	got := filter(input, "alpha|gamma")
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	for _, l := range lines {
		if !strings.Contains(l, "alpha") && !strings.Contains(l, "gamma") {
			t.Errorf("unexpected line in filter result: %q", l)
		}
	}
}

func TestRelevantInterfaces(t *testing.T) {
	re := regexp.MustCompile(`^(br[0-9]+|nwg[0-9]+)$`)
	_ = re

	input := `1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500
3: br0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500
4: nwg0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420
5: wlan0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500`

	got := relevantInterfaces(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 ifaces, got %d: %v", len(got), got)
	}
	if got[0] != "br0" || got[1] != "nwg0" {
		t.Errorf("unexpected ifaces: %v", got)
	}
}
