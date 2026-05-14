package wg

import (
	"testing"
)

func TestFindIfaceByIP(t *testing.T) {
	ipAddrOut := `1: lo: <LOOPBACK,UP> mtu 65536 qdisc noqueue state UNKNOWN
    inet 127.0.0.1/8 scope host lo
2: nwg0: <POINTOPOINT,NOARP,UP> mtu 1420 qdisc noqueue state UNKNOWN
    inet 10.8.0.2/24 scope global nwg0`

	tests := []struct {
		ip   string
		want string
	}{
		{"10.8.0.2", "nwg0"},
		{"127.0.0.1", "lo"},
		{"1.2.3.4", ""},
	}
	for _, tt := range tests {
		got := findIfaceByIP(ipAddrOut, tt.ip)
		if got != tt.want {
			t.Errorf("findIfaceByIP(_, %q) = %q, want %q", tt.ip, got, tt.want)
		}
	}
}

func TestFindWireguard(t *testing.T) {
	ifaces := []ifaceRecord{
		{ID: "Wireguard0", Type: "Wireguard", Description: "My VPN", State: "up"},
		{ID: "Wireguard1", Type: "Wireguard", Description: "Work VPN", State: "down"},
		{ID: "GigabitEthernet0", Type: "GigabitEthernet"},
	}

	tests := []struct {
		name    string
		wantID  string
		wantNil bool
	}{
		{"Wireguard0", "Wireguard0", false},
		{"My VPN", "Wireguard0", false},
		{"Work VPN", "Wireguard1", false},
		{"Unknown", "", true},
	}
	for _, tt := range tests {
		got, id := findWireguard(ifaces, tt.name)
		if tt.wantNil && got != nil {
			t.Errorf("findWireguard(_, %q): expected nil", tt.name)
		}
		if !tt.wantNil && (got == nil || id != tt.wantID) {
			t.Errorf("findWireguard(_, %q) = %v, %q; want id=%q", tt.name, got, id, tt.wantID)
		}
	}
}

func TestContainsWord(t *testing.T) {
	tests := []struct {
		s, word string
		want    bool
	}{
		{"inet 10.8.0.2/24 scope", "10.8.0.2", true},
		{"inet 10.8.0.20/24 scope", "10.8.0.2", false},
		{"inet 10.8.0.2/24", "10.8.0.2", true},
		{"10.8.0.2", "10.8.0.2", true},
	}
	for _, tt := range tests {
		got := containsWord(tt.s, tt.word)
		if got != tt.want {
			t.Errorf("containsWord(%q, %q) = %v, want %v", tt.s, tt.word, got, tt.want)
		}
	}
}
