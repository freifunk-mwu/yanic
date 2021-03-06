package respondd

import (
	"io/ioutil"
	"strings"
	"sync"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/lib/duration"

	babelState "github.com/Vivena/babelweb2/state"
)

func trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

type Daemon struct {
	DataInterval duration.Duration `toml:"data_interval"`
	Listen       []struct {
		Address   string `toml:"address"`
		Interface string `toml:"interface"`
		Port      int    `toml:"port"`
	} `toml:"listen"`

	Batman    []string               `toml:"batman"`
	Babel     string                 `toml:"babel"`
	babelData *babelState.BabelState `toml:"-"`

	dataMX          sync.Mutex
	dataByInterface map[string]*data.ResponseData

	Answer        *AnswerConfig            `toml:"defaults"`
	AnswerByZones map[string]*AnswerConfig `toml:"zones"`
}

type AnswerConfig struct {
	NodeID            string         `toml:"node_id"`
	Hostname          string         `toml:"hostname"`
	SiteCode          string         `toml:"site_code"`
	DomainCode        string         `toml:"domain_code"`
	Location          *data.Location `json:"location,omitempty"`
	VPN               bool           `toml:"vpn"`
	InterfaceMAC      string         `toml:"interface_mac"`
	InterfacesTraffic []string       `toml:"interfaces_traffic"`
	InterfacesAddress []string       `toml:"interfaces_address"`
}

func (d *Daemon) getAnswer(iface string) (*AnswerConfig, string) {
	config := d.Answer
	if v, ok := d.AnswerByZones[iface]; iface != "" && ok {
		config = v
	}

	nodeID := config.NodeID
	if nodeID == "" {
		if v, err := ioutil.ReadFile("/etc/machine-id"); err == nil {
			nodeID = trim(string(v))[:12]
		}
	}
	return config, nodeID
}
