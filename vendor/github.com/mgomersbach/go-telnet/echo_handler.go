package telnet

import (
	"io"
	"net"
)

// EchoHandler is a simple TELNET server which "echos" back to the client any (non-command)
// data back to the TELNET client, it received from the TELNET client.
var EchoHandler Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(ctx Context, conn *net.Conn, w Writer, r Reader) {
	io.Copy(w, r)
}
