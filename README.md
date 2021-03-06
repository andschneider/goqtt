<p align="left"><img src="logo.png" width="300" height="300"></p>

![Build](https://github.com/andschneider/goqtt/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/andschneider/goqtt)](https://goreportcard.com/report/github.com/andschneider/goqtt)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/andschneider/goqtt)](https://pkg.go.dev/mod/github.com/andschneider/goqtt)
[![License: MIT](https://img.shields.io/github/license/andschneider/goqtt)](https://img.shields.io/github/license/andschneider/goqtt)

# goqtt

goqtt is a simple MQTT client library. It is intended for lightweight MQTT applications where security and data guarantees are not important (see [limitations](https://github.com/andschneider/goqtt#limitations) below).

---

## Install

Using go `1.13` or higher, install using go modules:

`go get github.com/andschneider/goqtt`

## Usage

### subscription

Below is an example showing how to subscribe to a specified topic. Any Publish message received on the topic will be returned, allowing for further processing if desired.

```go
package main

import (
	"log"

	"github.com/andschneider/goqtt"
)

func main() {
	// Create Client
	topic := goqtt.Topic("goqtt")
	client := goqtt.NewClient("mqtt.eclipse.org", topic)

	// Connect to the specified MQTT broker
	err := client.Connect()
	if err != nil {
		log.Fatalf("could not connect to broker: %v\n", err)
	}
	defer client.Disconnect()

	// Subscribe to the topic
	err = client.Subscribe()
	if err != nil {
		log.Fatalf("could not subscribe to topic: %v\n", err)
	}

	// Read messages indefinitely
	for {
		log.Println("waiting for message")
		m, err := client.ReadLoop()
		if err != nil {
			log.Printf("error when reading in message: %v\n", err)
		}
		if m != nil {
			log.Printf("received message: '%s' from topic: '%s'", string(m.Message), m.Topic)
		}
	}
}
```

### publish

Below is an example showing how to Publish a message to a specified topic.

```go
package main

import (
	"log"

	"github.com/andschneider/goqtt"
)

func main() {
	// Create Client
	topic := goqtt.Topic("goqtt")
	client := goqtt.NewClient("mqtt.eclipse.org", topic)

	// Connect to the specified MQTT broker
	err := client.Connect()
	if err != nil {
		log.Fatalf("could not connect to broker: %v\n", err)
	}
	defer client.Disconnect()

	// Publish a message to the topic
	err = client.Publish("hello world")
	if err != nil {
		log.Printf("could not send message: %v\n", err)
	}
}
```

### more 

There are more examples in the `examples` folder. They have more client options exposed and produce verbose logging:

- connect and disconnect to a broker
- publishing a message a topic
- subscribing to a topic

#### Go

If you have Go installed you can run them with `go run ./examples/<example>`.

#### CLI

If you don't have Go installed, the examples have been compiled to binaries for Linux (x86 and Arm) and Darwin (macOS). These can found in the [releases](https://github.com/andschneider/goqtt/releases) page. Download the correct .tar.gz for your platform and once uncompressed, the binaries will be in an `examples` folder. These are CLIs and can be ran directly from your terminal. 

#### Docker

Each example has a Dockerfile if you don't have Go installed (and have Docker installed). To build an image, run `docker build -t <image-name> .` from an example directory. 

For example, from the `examples/connect` directory run `docker build -t conngoqtt .`. Then to use the container, run `docker run conngoqtt`.

## Limitations

Currently, goqtt does not implement the full [MQTT 3.1.1](http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html) client specification. 

The two main omissions from the spec are packets with **QoS > 0** and **username/password** message payloads in CONNECT packets. Because of these, **goqtt should be used with care in any environment that requires sensitive or important data**.

### QoS

For detailed information please refer to the spec, but with QoS set at 0 each packet is sent *at most once*. To quote the spec:

> The message is delivered according to the capabilities of the underlying network. No response is sent by the receiver and no retry is performed by the sender. The message arrives at the receiver either once or not at all.
  
If a packet is encountered by goqtt with a higher QoS it will likely crash.

### Authentication/Security

Without any support for username and password information to be sent to a broker there is no way for a client to authenticate. Please carefully consider if this is an acceptable security risk for your use case.

The MQTT spec does not require the use of TLS/SSL, however there are brokers that support TLS. **goqtt does not support TLS connections.**  Please carefully consider if this is an acceptable security risk for your use case.

### Other

A more minor omission is the lack of support for multiple topic subscriptions or wildcard subscriptions. A topic must be fully formed, e.g. `goqtt/hello`.

## Logging

*WIP*

goqtt uses [zerolog](https://github.com/rs/zerolog) internally for creating structured logs with different levels (such as DEBUG, INFO, and ERROR). 

However, all log statements are **disabled** by default. If you'd like to enable the logging, set the global log level to the desired level. This is done with the `zerolog.SetGlobalLevel()` function and a parameter such as `zerolog.InfoLevel`. See the examples and their documentation for more levels.

## Inspiration

The [eclipse paho golang](https://github.com/eclipse/paho.mqtt.golang) MQTT library has been a large inspiration and resource for this project. Much thanks to them.
