package network

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
)

type ifaceRecord struct {
	Address       string `json:"address"`
	ID            string `json:"id"`
	Index         int    `json:"index"`
	Description   string `json:"description"`
	InterfaceName string `json:"interface-name"`
	Type          string `json:"type"`
}

type Bridge struct {
	Name        string
	LinuxName   string
	Address     string
	Description string
	Index       int
}

func Bridges(ctx context.Context) []Bridge {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:79/rci/show/interface", nil)
	if err != nil {
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var raw map[string]ifaceRecord
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil
	}
	out := make([]Bridge, 0, len(raw))
	for _, v := range raw {
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

func bridgeName(v ifaceRecord) string {
	if v.InterfaceName != "" {
		return v.InterfaceName
	}
	return v.ID
}

func bridgeLinuxName(index int) string {
	return "br" + strconv.Itoa(index)
}
