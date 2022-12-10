package testing

import (
	"os"

	"github.com/nuqz/helvar-go/members"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Network struct {
	Clusters []members.Cluster
	Routers  []members.Router
	Groups   []members.Group
}

func NetFromYAML(bs []byte) (Network, error) {
	out := Network{}

	if err := yaml.Unmarshal(bs, &out); err != nil {
		return out, err
	}

	return out, nil
}

func NetFromYAMLFile(path string) (Network, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return Network{}, err
	}

	return NetFromYAML(bs)
}

func MustNetFromYAMLFile(path string) Network {
	net, err := NetFromYAMLFile(path)
	if err != nil {
		logrus.Fatal(err)
	}

	return net
}

func (n Network) GetClusterIDs() []uint8 {
	out := make([]uint8, len(n.Clusters))
	for i, c := range n.Clusters {
		out[i] = c.ID
	}

	return out
}

func (n Network) GetRouterIDs() []uint8 {
	out := make([]uint8, len(n.Routers))
	for i, r := range n.Routers {
		out[i] = r.ID
	}

	return out
}

func (n Network) GetGroupIDs() []uint16 {
	out := make([]uint16, len(n.Groups))
	for i, g := range n.Groups {
		out[i] = g.ID
	}

	return out
}

func (n Network) GetGroupByID(id uint16) members.Group {
	for _, g := range n.Groups {
		if g.ID == id {
			return g
		}
	}

	return members.Group{}
}

func (n Network) GetGroupDevices(id uint16) []string {
	g := n.GetGroupByID(id)
	out := make([]string, len(g.Devices))
	for i, d := range g.Devices {
		out[i] = "@" + d.Address
	}

	return out
}

func (n Network) GetDeviceByAddress(addr string) members.Device {
	for _, g := range n.Groups {
		for _, d := range g.Devices {
			if d.Address == addr {
				return d
			}
		}
	}

	return members.Device{}
}

// TODO: ...
// DefaultNet is virtual HelvarNET network for testing this particular package.
// var DefaultNet = MustNetFromYAMLFile(path.Join(
// 	os.Getenv("GOPATH"), "src", "github.com", "nuqz", "helvar-go",
// 	"testing", "test_net.yml"))
