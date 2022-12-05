package helvargo

import (
	"fmt"
	"image/color"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/nuqz/chanfan"
	"github.com/nuqz/helvar-go/members"
	"github.com/nuqz/helvar-go/message"
	"github.com/pkg/errors"
)

// Client is a HelvarNET client.
type Client struct {
	host      string
	hostParts []string
	port      int
	address   string

	toSend chan<- *chanfan.IO[*message.Message, *message.Message]
}

// NewClient returns new client, which will communicate to specified router.
func NewClient(host string, port int) *Client {
	if host == "localhost" {
		host = "127.0.0.1"
	}

	address := fmt.Sprintf("%s:%d", host, port)
	return &Client{
		host:      host,
		hostParts: strings.Split(host, "."),
		port:      port,
		address:   address,
	}
}

// IsSameSubnet returns true when ... TODO
func (c *Client) IsSameSubnet(addr string) bool {
	addrParts := strings.Split(addr, ".")
	return addrParts[0] == c.hostParts[2] && addrParts[1] == c.hostParts[3]
}

// Connect ... TODO ...
func (c *Client) Connect(nTransceivers, bufSize int) ([]<-chan error, error) {
	in := make(chan *chanfan.IO[*message.Message, *message.Message], bufSize)
	errs := make([]<-chan error, bufSize)
	c.toSend = in
	for i := 0; i < nTransceivers; i++ {
		conn, err := net.Dial("tcp", c.address)
		if err != nil {
			return nil, errors.Wrapf(err,
				"couldn't establish connection #%d to %s", i+1, c.address)
		}

		errs[i] = NewTransceiver(conn, in).Go()
	}

	return errs, nil
}

func (c *Client) Disconnect() { close(c.toSend) }

func (c *Client) Transceive(msg *message.Message) (*message.Message, error) {
	ret := make(chan *chanfan.Result[*message.Message])
	c.toSend <- chanfan.NewIO(msg, ret)
	resp := <-ret

	if resp.Error != nil {
		return nil, errors.Wrap(resp.Error, "failed to transceive message")
	}

	if resp.Value != nil && resp.Value.Type == message.TError {
		// TODO: add error description by it's code
		return nil, errors.Errorf("received a message with an error: %s",
			resp.Value)
	}

	return resp.Value, nil
}

func (c *Client) queryIDs(msg *message.Message) ([]int, error) {
	msg, err := c.Transceive(msg)
	if err != nil {
		return nil, err
	}

	out, err := msg.AnswerIDs()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetClusters() ([]members.Cluster, error) {
	clusterIDs, err := c.queryIDs(message.NewQueryClusters())
	if err != nil {
		return nil, err
	}

	clusters := make([]members.Cluster, len(clusterIDs))
	for i, id := range clusterIDs {
		clusters[i] = members.Cluster{ID: uint8(id)}
	}

	return clusters, nil
}

func (c *Client) GetRouters() ([]members.Router, error) {
	routerIDs, err := c.queryIDs(message.NewQueryRouters("0"))
	if err != nil {
		return nil, err
	}

	routers := make([]members.Router, len(routerIDs))
	for i, id := range routerIDs {
		routers[i] = members.Router{ID: uint8(id)}
	}

	return routers, nil
}

func (c *Client) GetGroups() ([]members.Group, error) {
	groupIDs, err := c.queryIDs(message.NewQueryGroups())
	if err != nil {
		return nil, err
	}

	groups := make([]members.Group, len(groupIDs))
	for i, id := range groupIDs {
		groups[i] = members.Group{ID: uint16(id)}
	}

	return groups, nil
}

// UpdateName sends query group description command to a given router r and
// then updates group name with response it received. Returns an error if
// something goes wrong.
func (c *Client) GetGroupName(g members.Group) (string, error) {
	msg, err := c.Transceive(message.NewQueryGroupDescription(g.ID))
	if err != nil {
		return "", err
	}

	return msg.Answer, nil
}

// UpdateDevices sends query group command to a given router r and then
// updates group devices with name it received. Returns a slice of devices as
// a result or an error if something goes wrong.
func (c *Client) GetDevices(g members.Group) ([]members.Device, error) {
	msg, err := c.Transceive(message.NewQueryGroup(g.ID))
	if err != nil {
		return nil, err
	}

	devices := []members.Device{}
	for _, rawAddr := range msg.AnswerStrings() {
		addr := strings.TrimLeft(rawAddr, "@")

		// TODO: ...
		// Only lookup devices with same cluster and router ID as current
		// router
		// if c.IsSameSubnet(addr) {
		devices = append(devices, members.Device{Address: addr})
		// }
	}

	return devices, nil
}

// UpdateName sends query device description to a given router r and then
// updates device name with response it received. Returns an error if
// something goes wrong.
func (c *Client) GetDeviceName(d members.Device) (string, error) {
	msg, err := c.Transceive(message.NewQueryDeviceDescription(d.Address))
	if err != nil {
		return "", err
	}

	return msg.Answer, nil
}

func (c *Client) GetDeviceState(d members.Device) (members.DeviceState, error) {
	// TODO
	return 0, errors.New("not implemented")
}

func (c *Client) GetTime() (time.Time, error) {
	reply, err := c.Transceive(message.NewQueryTime())
	if err != nil {
		return time.Time{}, err
	}

	ts, err := strconv.ParseInt(reply.Answer, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrapf(err,
			"failed to query network time - timestamp %s is invalid integer",
			reply.Answer)
	}

	return time.Unix(ts, 0), nil
}

func (c *Client) RecallSceneGroup(
	gid uint16,
	block, scene uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewRecallSceneGroup(gid, block, scene, params...))
	return err
}

func (c *Client) RecallSceneDevice(
	addr string,
	block, scene uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewRecallSceneDevice(addr, block, scene, params...))
	return err
}

func (c *Client) DirectLevelGroup(
	gid uint16,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewDirectLevelGroup(gid, level, params...))
	return err
}

func (c *Client) DirectLevelDevice(
	addr string,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewDirectLevelDevice(addr, level, params...))
	return err
}

func (c *Client) ColorTemperatureGroup(
	gid, tempK uint16,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewColorTemperatureGroup(gid, tempK, level, params...))
	return err
}

func (c *Client) ColorTemperatureDevice(
	addr string,
	tempK uint16,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewColorTemperatureDevice(addr, tempK, level, params...))
	return err
}

func (c *Client) ColorGroup(
	gid uint16,
	color color.Color,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewColorGroup(gid, color, level, params...))
	return err
}

func (c *Client) ColorDevice(
	addr string,
	color color.Color,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewColorDevice(addr, color, level, params...))
	return err
}

func (c *Client) RGBGroup(
	gid uint16,
	r, g, b byte,
	level uint8,
	params ...message.Parameter,

) error {
	_, err := c.Transceive(
		message.NewRGBGroup(gid, r, g, b, level, params...))
	return err
}

func (c *Client) RGBDevice(
	addr string,
	r, g, b byte,
	level uint8,
	params ...message.Parameter,
) error {
	_, err := c.Transceive(
		message.NewRGBDevice(addr, r, g, b, level, params...))
	return err
}
