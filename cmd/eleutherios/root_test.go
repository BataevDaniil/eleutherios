package eleutherios

import (
	"bytes"
	"strings"
	"testing"
)

// TestRootHelpDoesNotError ловит сломанные флаги/команды на этапе init():
// regressions в MarkHidden, дубликаты команд, неправильные annotations.
func TestRootHelpDoesNotError(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd --help: %v\nOutput:\n%s", err, buf.String())
	}
	out := buf.String()
	for _, want := range []string{"start", "stop", "status", "--timeout"} {
		if !strings.Contains(out, want) {
			t.Errorf("help не содержит %q, output:\n%s", want, out)
		}
	}
}

func TestNeedsRoot(t *testing.T) {
	if !needsRoot(startCmd) {
		t.Error("startCmd должна требовать root")
	}
	if !needsRoot(stopCmd) {
		t.Error("stopCmd должна требовать root")
	}
	if !needsRoot(statusCmd) {
		t.Error("statusCmd должна требовать root")
	}
	if needsRoot(vpnLsCmd) {
		t.Error("vpnLsCmd не должна требовать root")
	}
	if needsRoot(netLsCmd) {
		t.Error("netLsCmd не должна требовать root")
	}
}
