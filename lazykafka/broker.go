package lazykafka

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Broker we'll use to connect to the kafka cluster
type Broker struct {
	id        int32
	addr      string
	conf      *Config
	conn      net.Conn
	connErr   error
	lock      sync.Mutex
	opened    int32
	responses chan responsePromise
	done      chan bool
}

type responsePromise struct {
	requestTime   time.Time
	correlationID int32
	packets       chan []byte
	errors        chan error
}

// NewBroker blah
func NewBroker(addr string) *Broker {
	return &Broker{id: -1, addr: addr}
}

// Open function is responsible to open a connection with broker of given address
func (b *Broker) Open(conf *Config) error {
	if !atomic.CompareAndSwapInt32(&b.opened, 0, 1) {
		return ErrAlreadyConnected
	}

	if conf == nil {
		conf = NewConfig()
	}

	b.lock.Lock()

	go withRecover(func() {
		defer b.lock.Unlock()

		dialer := net.Dialer{
			Timeout:   conf.Net.DialTimeout,
			KeepAlive: conf.Net.KeepAlive,
		}

		if conf.Net.TLS.Enable {
			b.conn, b.connErr = tls.DialWithDialer(&dialer, "tcp", b.addr, conf.Net.TLS.Config)
		} else {
			b.conn, b.connErr = dialer.Dial("tcp", b.addr)
		}
		if b.connErr != nil {
			Logger.Printf("Failed to connect to broker %s: %s\n", b.addr, b.connErr)
			b.conn = nil
			atomic.StoreInt32(&b.opened, 0)
			return
		}
		b.conn = newBufConn(b.conn)

		b.conf = conf

		b.done = make(chan bool)
		b.responses = make(chan responsePromise, b.conf.Net.MaxOpenRequests-1)

		if b.id >= 0 {
			Logger.Printf("Connected to broker at %s (registered as #%d)\n", b.addr, b.id)
		} else {
			Logger.Printf("Connected to broker at %s (unregistered)\n", b.addr)
		}
		go withRecover(b.responseReceiver)
	})

	return nil
}

// Connected blah
func (b *Broker) Connected() (bool, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.conn != nil, b.connErr
}

// Close blah
func (b *Broker) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.conn == nil {
		return ErrNotConnected
	}

	close(b.responses)
	<-b.done

	err := b.conn.Close()

	b.conn = nil
	b.connErr = nil
	b.done = nil
	b.responses = nil

	if err == nil {
		Logger.Printf("Closed connection to broker %s\n", b.addr)
	} else {
		Logger.Printf("Error while closing connection to broker %s: %s\n", b.addr, err)
	}

	atomic.StoreInt32(&b.opened, 0)

	return err
}

// ID returns the broker ID retrieved from Kafka's metadata, or -1 if that is not known.
func (b *Broker) ID() int32 {
	return b.id
}

func (b *Broker) responseReceiver() {
	var dead error
	header := make([]byte, 8)
	for response := range b.responses {
		if dead != nil {
			response.errors <- dead
			continue
		}

		err := b.conn.SetReadDeadline(time.Now().Add(b.conf.Net.ReadTimeout))
		if err != nil {
			dead = err
			response.errors <- err
			continue
		}

		bytesReadHeader, err := io.ReadFull(b.conn, header)
		Logger.Printf("%d\n", bytesReadHeader)
		if err != nil {
			dead = err
			response.errors <- err
			continue
		}

		decodedHeader := responseHeader{}
		err = decode(header, &decodedHeader)
		if err != nil {
			dead = err
			response.errors <- err
			continue
		}
		if decodedHeader.correlationID != response.correlationID {
			// TODO if decoded ID < cur ID, discard until we catch up
			// TODO if decoded ID > cur ID, save it so when cur ID catches up we have a response
			dead = PacketDecodingError{fmt.Sprintf("correlation ID didn't match, wanted %d, got %d", response.correlationID, decodedHeader.correlationID)}
			response.errors <- dead
			continue
		}

		buf := make([]byte, decodedHeader.length-4)
		bytesReadBody, err := io.ReadFull(b.conn, buf)
		Logger.Printf("%d\n", bytesReadBody)
		if err != nil {
			dead = err
			response.errors <- err
			continue
		}

		response.packets <- buf
	}
	close(b.done)
}
