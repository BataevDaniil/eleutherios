package ipset

import (
	"fmt"
	"os/exec"
)

const (
	SetRU       = "ELEUTHERIOS_RU"
	SetExcluded = "ELEUTHERIOS_EXCLUDED"
	TTL         = "86400" // 24 часа
)

func CreateSets() error {
	for _, set := range []struct {
		name        string
		ttl         string
		storageType string
	}{
		{SetRU, TTL, "hash:ip"},
		{SetExcluded, "", "hash:net"},
	} {
		args := []string{"create", set.name, set.storageType, "family", "inet", "-exist"}
		if set.ttl != "" {
			args = append(args, "timeout", set.ttl)
		}
		out, err := exec.Command("ipset", args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("ipset create %s: %w (%s)", set.name, err, out)
		}
	}
	if err := FillExcluded(); err != nil {
		return err
	}
	fmt.Printf("  ipset %s и %s созданы\n", SetRU, SetExcluded)
	return nil
}

func FillExcluded() error {
	reserved := []string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"78.47.125.180",
	}
	for _, ip := range reserved {
		out, err := exec.Command("ipset", "-exist", "add", SetExcluded, ip).CombinedOutput()
		if err != nil {
			return fmt.Errorf("ipset add %s %s: %w (%s)", SetExcluded, ip, err, out)
		}
	}
	return nil
}

func DestroySets() {
	for _, name := range []string{SetRU, SetExcluded} {
		exec.Command("ipset", "destroy", name).Run()
	}
	fmt.Printf("  ipset %s и %s удалены\n", SetRU, SetExcluded)
}
