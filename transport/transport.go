// Package transport is an interface for synchronous connection based communication
package transport

import (
	"time"
)

// Transport is an interface which is used for communication between
// services. It uses connection based socket send/recv semantics and
// has various implementations; http, grpc, quic.
type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Client, error) //Client内部包含一个到addr的持久连接（conn)，通过conn来处理http请求，内部维护一个缓存，缓存请求的响应
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}

type Message struct {
	Header map[string]string
	Body   []byte
}

type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
	Local() string
	Remote() string
}

type Client interface {
	Socket
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error //默认是内部http的listen函数
}

type Option func(*Options)

type DialOption func(*DialOptions)

type ListenOption func(*ListenOptions)

var (
	DefaultTransport Transport = newHTTPTransport()

	DefaultDialTimeout = time.Second * 5
)

func NewTransport(opts ...Option) Transport {
	return newHTTPTransport(opts...)
}
