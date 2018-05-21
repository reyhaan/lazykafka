package lazykafka

import (
	"bufio"
	"net"
)

func withRecover(fn func()) {
	defer func() {
		handler := PanicHandler
		if handler != nil {
			if err := recover(); err != nil {
				handler(err)
			}
		}
	}()

	fn()
}

func safeAsyncClose(b *Broker) {
	tmp := b // local var prevents clobbering in goroutine
	go withRecover(func() {
		if connected, _ := tmp.Connected(); connected {
			if err := tmp.Close(); err != nil {
				Logger.Println("Error closing broker", tmp.ID(), ":", err)
			}
		}
	})
}

// Encoder is a simple interface for any type that can be encoded as an array of bytes
// in order to be sent as the key or value of a Kafka message. Length() is provided as an
// optimization, and must return the same as len() on the result of Encode().
type Encoder interface {
	Encode() ([]byte, error)
	Length() int
}

// make strings and byte slices encodable for convenience so they can be used as keys
// and/or values in kafka messages

// StringEncoder implements the Encoder interface for Go strings so that they can be used
// as the Key or Value in a ProducerMessage.
type StringEncoder string

// Encode blah
func (s StringEncoder) Encode() ([]byte, error) {
	return []byte(s), nil
}

// Length blah
func (s StringEncoder) Length() int {
	return len(s)
}

// ByteEncoder implements the Encoder interface for Go byte slices so that they can be used
// as the Key or Value in a ProducerMessage.
type ByteEncoder []byte

// Encode blah
func (b ByteEncoder) Encode() ([]byte, error) {
	return b, nil
}

// Length blah
func (b ByteEncoder) Length() int {
	return len(b)
}

// bufConn wraps a net.Conn with a buffer for reads to reduce the number of
// reads that trigger syscalls.
type bufConn struct {
	net.Conn
	buf *bufio.Reader
}

func newBufConn(conn net.Conn) *bufConn {
	return &bufConn{
		Conn: conn,
		buf:  bufio.NewReader(conn),
	}
}

func (bc *bufConn) Read(b []byte) (n int, err error) {
	return bc.buf.Read(b)
}
