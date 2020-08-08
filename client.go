package goqtt

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// Default values for the Client configuration. See the DefaultClientId function for
// the Client's default ClientId value.
const (
	DefaultPort      = "1883"
	DefaultKeepAlive = 60 // seconds
	DefaultTopic     = "goqtt"
)

// clientConfig holds the information needed to create a new Client.
// As it's unexported these can't be set directly, instead use a functional
// option, like Port.
type clientConfig struct {
	clientId  string
	keepAlive int
	// server is the address (IP or domain) of the broker.
	server string
	port   string
	topic  string
}

// Client is the main interaction point for sending and receiving Packets. An
// instantiated Client needs to call the Connect() method before sending/receiving any packets.
type Client struct {
	Config *clientConfig
	conn   net.Conn
}

// NewClient creates a Client struct which can be used to interact with a MQTT broker.
// It sets default values for most of the configuration, so only a server address is
// required to instantiate it. Other configuration options are available and can be
// passed in as needed.
func NewClient(addr string, opts ...option) *Client {
	var c = &Client{}
	// default configuration
	d := &clientConfig{
		clientId:  DefaultClientId(),
		keepAlive: DefaultKeepAlive,
		server:    addr,
		port:      DefaultPort,
		topic:     DefaultTopic,
	}
	c.Config = d
	// apply any other options that have been passed in
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// DefaultClientId creates a default value for the Client's clientId.
// It uses the process id and the current time to try to prevent client collisions with the broker.
func DefaultClientId() string {
	cid := fmt.Sprintf("%s-%d-%s", "goqtt", os.Getpid(), strconv.Itoa(time.Now().Second()))
	return cid
}

// option is a function used to modify/set a configuration field in a Client.
type option func(*Client)

// ClientId is the identifier you use to tell the MQTT broker who you are.
// A broker will require unique ClientId's for all connected clients.
func ClientId(cid string) option {
	return func(c *Client) {
		c.Config.clientId = cid
	}
}

// KeepAlive is the time in seconds that the broker should wait before
// disconnecting you if no packets are sent.
func KeepAlive(seconds int) option {
	return func(c *Client) {
		c.Config.keepAlive = seconds
	}
}

// Port is the port number of the server. Usually 1883 (insecure) or 8883 (TLS).
func Port(port string) option {
	return func(c *Client) {
		c.Config.port = port
	}
}

// Topic is the fully qualified topic to subscribe or publish to.
// Only single topics and no wildcards are accepted at this time.
func Topic(topic string) option {
	return func(c *Client) {
		c.Config.topic = topic
	}
}
