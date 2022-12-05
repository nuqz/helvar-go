package testing

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/nuqz/helvar-go/members"
	"github.com/nuqz/helvar-go/message"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
)

type Network struct {
	Clusters []members.Cluster
	Routers  []members.Router
	Groups   []members.Group
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

func toStrs[T constraints.Integer](in []T) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = strconv.Itoa(int(v))
	}

	return out
}

func joinIDs[T constraints.Integer](ids []T) string {
	return strings.Join(toStrs(ids), message.Delimiter.String())
}

func joinStrs(strs []string) string {
	return strings.Join(strs, message.Delimiter.String())
}

// Net is virtual HelvarNET network for testing purpose.
var Net = Network{
	Clusters: []members.Cluster{{ID: 1}, {ID: 2}, {ID: 3}},
	Routers:  []members.Router{{ID: 251}, {ID: 252}, {ID: 253}},
	Groups: []members.Group{
		{
			ID:        11,
			Name:      "Group 11",
			LastScene: 0,
			Devices: []members.Device{
				{Address: "1.251.1", Name: "Lamp 1 in Group 11", State: members.OK},
			},
		},
		{
			ID:        12,
			Name:      "Group 12",
			LastScene: 0,
			Devices: []members.Device{
				{Address: "1.252.1", Name: "Lamp 1 in Group 12", State: members.OK},
			},
		},
		{
			ID:        13,
			Name:      "Group 13",
			LastScene: 0,
			Devices: []members.Device{
				{Address: "1.253.1", Name: "Lamp 1 in Group 13", State: members.OK},
			},
		},
	},
}

type Router struct {
	Address string

	net Network

	listener  net.Listener
	listening bool

	log *logrus.Entry
}

func NewRouter(addr string, net Network) *Router {
	return &Router{
		Address: addr,
		net:     net,
		log:     logrus.New().WithField("package", "helvar-go/testing"),
	}
}

func (r *Router) IsListening() bool {
	return r.listening
}

func (r *Router) Listen() error {
	var err error
	r.listener, err = net.Listen("tcp", r.Address)
	if err != nil {
		return err
	}

	r.listening = true
	go func() {
		defer func() {
			r.listening = false
			if err := r.listener.Close(); err != nil {
				r.log.WithError(err).
					Error("failed to close listener properly")
			}
			r.log.Info("listener stopped")
		}()

		r.log.Infof("listening @%s", r.Address)
		for {
			conn, err := r.listener.Accept()
			if err != nil {
				r.log.WithError(err).Error("failed to accept a connection")
				continue
			}

			go func() {
				if err := r.handleClient(conn); err != nil {
					r.log.WithField("client", conn.RemoteAddr().String()).
						Error(err)
				}
			}()
		}
	}()

	return nil
}

func (r *Router) handleClient(conn net.Conn) error {
	log := r.log.WithField("client", conn.RemoteAddr().String())
	log.Info("handling client")

	defer func() {
		if err := conn.Close(); err != nil {
			log.WithError(err).
				WithField("client", conn.RemoteAddr().String()).
				Error("failed to close client connection properly")
		}
	}()

	buf := bufio.NewReader(conn)
	for {
		requestStr, err := buf.ReadString(message.Terminator.Byte())
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		log.Infof("received: %s", requestStr)

		msg, err := message.Parse(requestStr)
		if err != nil {
			log.WithError(err).
				Error("failed to parse incoming message: %s", requestStr)
			continue
		}

		cmdID := msg.GetCommandID()

		// TODO: Don't know if real router may respond with error message.
		if !message.NeedResponse(msg) {
			switch cmdID {
			case message.RecallSceneGroup:
				// TODO: add something meaningful
			case message.RecallSceneDevice:
				// TODO: add something meaningful
			case message.DirectLevelGroup:
				// TODO: add something meaningful
			case message.DirectLevelDevice:
				// TODO: add something meaningful
			default:
				// TODO: Unsupported command error
			}
			continue
		}

		reply := &message.Message{
			Type:       message.TReply,
			Parameters: msg.Parameters,
		}

		switch cmdID {
		case message.QueryClusters:
			reply.Answer = joinIDs(Net.GetClusterIDs())
		case message.QueryRouters:
			reply.Answer = joinIDs(Net.GetRouterIDs())
		case message.QueryGroups:
			reply.Answer = joinIDs(Net.GetGroupIDs())
		case message.QueryGroupDescription:
			reply.Answer = Net.GetGroupByID(msg.GetGroupID()).Name
		case message.QueryGroup:
			reply.Answer = joinStrs(Net.GetGroupDevices(msg.GetGroupID()))
		case message.QueryDeviceDescription:
			reply.Answer = Net.GetDeviceByAddress(msg.GetAddress()).Name
		case message.QueryTime:
			reply.Answer = strconv.Itoa(int(time.Now().Unix()))
		case message.NoCommand:
			fallthrough
		default:
			reply.Type = message.TError
			reply.Answer = strconv.Itoa(int(message.EInvalidMessageCommand))
		}

		out := reply.Bytes()
		if n, err := conn.Write(out); err != nil {
			log.WithError(err).
				Error("unable to write response %s for incoming message %s",
					reply, msg)
		} else if n != len(out) {
			log.Error(
				"only part of response %s was sent for incoming message %s",
				reply, msg)
		} else {
			log.Infof("sent: %s", reply)
		}
	}

	return nil
}
