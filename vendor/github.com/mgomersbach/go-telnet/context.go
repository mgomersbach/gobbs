package telnet

import (
	"net"
)

type Context interface {
	Logger() Logger
	Conn() net.Conn

	InjectLogger(Logger) Context
	InjectConn(net.Conn) Context
}

type internalContext struct {
	logger Logger
	conn   net.Conn
}

func NewContext() Context {
	ctx := internalContext{}

	return &ctx
}

func (ctx *internalContext) Logger() Logger {
	return ctx.logger
}

func (ctx *internalContext) InjectLogger(logger Logger) Context {
	ctx.logger = logger

	return ctx
}

func (ctx *internalContext) Conn() net.Conn {
	return ctx.conn
}

func (ctx *internalContext) InjectConn(conn net.Conn) Context {
	ctx.conn = conn

	return ctx
}
