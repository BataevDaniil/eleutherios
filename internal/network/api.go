package network

import (
	"context"
	"sort"
	"strconv"

	"github.com/BataevDaniil/eleutherios/internal/keenetic"
)

type Bridge struct {
	Name        string
	LinuxName   string
	Address     string
	Description string
	Index       int
}

func Bridges(ctx context.Context) []Bridge {
	ifaces, err := keenetic.ShowInterfaces(ctx)
	if err != nil {
		return nil
	}
	out := make([]Bridge, 0, len(ifaces))
	for _, v := range ifaces {
		if v.Type != "Bridge" {
			continue
		}
		out = append(out, Bridge{
			Name:        bridgeName(v),
			LinuxName:   bridgeLinuxName(v.Index),
			Address:     v.Address,
			Description: v.Description,
			Index:       v.Index,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Index < out[j].Index })
	return out
}

func bridgeName(v keenetic.Interface) string {
	if v.InterfaceName != "" {
		return v.InterfaceName
	}
	return v.ID
}

func bridgeLinuxName(index int) string {
	return "br" + strconv.Itoa(index)
}
