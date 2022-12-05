package helvargo

import (
	"bufio"
	"net"
	"time"

	"github.com/nuqz/chanfan"
	"github.com/nuqz/helvar-go/message"
	"github.com/pkg/errors"
)

const KeepAliveDuration = 120 * time.Second

type Transceiver struct {
	*chanfan.Transceiver[*message.Message, *message.Message]

	conn net.Conn
	r    *bufio.Reader
}

var terminatorByte = message.Terminator.Byte()

func NewTransceiver(
	conn net.Conn,
	in <-chan *chanfan.IO[*message.Message, *message.Message],
) *Transceiver {
	t := chanfan.NewTransceiver(in)
	t.KeepAliveDuration = KeepAliveDuration
	t.Terminate = func() error {
		if err := conn.Close(); err != nil {
			return errors.Wrap(err,
				"failed to close transceiver connection properly")
		}
		return nil
	}

	return &Transceiver{
		Transceiver: t,

		conn: conn,
		r:    bufio.NewReader(conn),
	}
}

func (t *Transceiver) transceive(
	msg *message.Message,
) (*message.Message, error) {
	out := msg.Bytes()
	if n, err := t.conn.Write(out); err != nil {
		return nil, errors.Wrapf(err,
			"failed to sent message: %s", msg)
	} else if n != len(out) {
		return nil, errors.Wrapf(err,
			"message was sent partially: %s", msg)
	}

	if message.NeedResponse(msg) {
		resp, err := t.r.ReadString(terminatorByte)
		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to receive response for: %s", msg)
		}

		return message.ParsePartial(resp)
	}

	return nil, nil
}

func (t *Transceiver) Go() <-chan error {
	t.KeepAlive = func() error {
		_, err := t.transceive(message.NewQueryTime())
		return err
	}

	return t.Transceiver.Go(t.transceive)
}
