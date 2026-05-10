package dns

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Cleanup() {
	os.Remove(ConfFile)
	if readPID(PIDFile) != "" || pidOf("dnsmasq") != "" {
		_ = restart()
	}
	fmt.Printf("  dnsmasq: конфиг %s удалён\n", ConfFile)
}

func Status() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("init: %s\n", pathState(InitFile)))
	b.WriteString(fmt.Sprintf("pid-file: %s (%s)\n", PIDFile, readPID(PIDFile)))
	b.WriteString(fmt.Sprintf("pidof dnsmasq: %s\n", valueOr(pidOf("dnsmasq"), "нет")))
	out, err := exec.Command(InitFile, "status").CombinedOutput()
	b.WriteString(fmt.Sprintf("$ %s status: %v\n%s\n", InitFile, err, strings.TrimSpace(string(out))))
	b.WriteString(fmt.Sprintf("config: %s\n", pathState(ConfFile)))
	if data, err := os.ReadFile(ConfFile); err == nil {
		b.WriteString(strings.TrimSpace(string(data)) + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func pathState(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "нет (" + err.Error() + ")"
	}
	return fmt.Sprintf("есть, %d байт", info.Size())
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
