package network

import "fmt"

func List() string {
	if bridges := Bridges(); len(bridges) > 0 {
		return formatBridges(bridges)
	}
	out := ""
	for _, name := range ListNetworks() {
		ip, err := GetNetIP(name)
		if err != nil {
			ip = "-"
		}
		out += fmt.Sprintf("%s\t%s\t%s\t%s\n", name, name, ip, "-")
	}
	return out
}

func formatBridges(bridges []Bridge) string {
	out := ""
	for _, b := range bridges {
		desc := b.Description
		if desc == "" {
			desc = "-"
		}
		out += fmt.Sprintf("%s\t%s\t%s\t%s\n", b.Name, b.LinuxName, b.Address, desc)
	}
	return out
}
