package testing

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/nuqz/helvar-go/message"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
)

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
			reply.Answer = joinIDs(r.net.GetClusterIDs())
		case message.QueryRouters:
			reply.Answer = joinIDs(r.net.GetRouterIDs())
		case message.QueryGroups:
			reply.Answer = joinIDs(r.net.GetGroupIDs())
		case message.QueryGroupDescription:
			reply.Answer = r.net.GetGroupByID(msg.GetGroupID()).Name
		case message.QueryGroup:
			reply.Answer = joinStrs(r.net.GetGroupDevices(msg.GetGroupID()))
		case message.QueryDeviceDescription:
			reply.Answer = r.net.GetDeviceByAddress(msg.GetAddress()).Name
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
