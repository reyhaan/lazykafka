package lazykafka

import (
	"crypto/tls"
	"regexp"
	"time"
)

const defaultClientID = "lazykafka"

var validID = regexp.MustCompile(`\A[A-Za-z0-9._-]+\z`)

// Config used to describe default configs for various kafka stuffs
type Config struct {
	Net struct {
		MaxOpenRequests int
		DialTimeout     time.Duration
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration

		TLS struct {
			Enable bool
			Config *tls.Config
		}

		KeepAlive time.Duration
	}

	Metadata struct {
	}

	Producer struct {
	}

	Consumer struct {
	}

	ClientID          string
	ChannelBufferSize int
}

// NewConfig blah
func NewConfig() *Config {
	c := &Config{}

	c.Net.MaxOpenRequests = 5
	c.Net.DialTimeout = 30 * time.Second
	c.Net.ReadTimeout = 30 * time.Second
	c.Net.WriteTimeout = 30 * time.Second

	c.ClientID = defaultClientID
	c.ChannelBufferSize = 256

	return c
}
