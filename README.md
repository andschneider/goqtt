# goqtt

![Build](https://github.com/andschneider/goqtt/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/andschneider/goqtt)](https://goreportcard.com/report/github.com/andschneider/goqtt)
[![License: MIT](https://img.shields.io/github/license/andschneider/goqtt)](https://img.shields.io/github/license/andschneider/goqtt)

---

## Install

Using go `1.13` or higher, install using go modules:

`go get github.com/andschneider/goqtt`

## Usage

### subscription

Below is a simple example showing a subscription to a specified topic. Any Publish message received on the topic will be returned, allowing for further processing if desired. *Note: there is no error handling in the example, which is not advised.*

```go
package main

import (
	"log"

	"github.com/andschneider/goqtt"
)

func main() {
	// Create Client
	clientId := "simple"
	keepAlive := 30 // seconds
	broker := "mqtt.eclipse.org:1883"
	topic := "goqtt"

	client := goqtt.NewClient(goqtt.NewClientConfig(clientId, keepAlive, broker, topic))

	// Attempt a connection to the specified MQTT broker
	client.Connect()
	defer client.Disconnect()

	// Attempt to subscribe to the topic
	client.Subscribe()

	// Read messages indefinitely
	for {
		log.Println("waiting for message")
		m, _ := client.ReadLoop()
		if m != nil {
			log.Printf("received message: '%s' from topic: '%s'", string(m.Message), m.Topic)
		}
	}
}
```

There are more example in the `examples` folder. You can find CLIs to do the following:

- simple connect and disconnect to a broker
- publishing a message a topic
- subscribing to a topic
- a minimal broker

If you have Go installed you can run them with `go run ./examples/<example>`.

If you don't have Go installed, the examples have been compiled to binaries for Linux (x86 and Arm) and Darwin (macOS). These can found in the [releases](https://github.com/andschneider/goqtt/releases) page. Download the correct .tar.gz for your platform and once uncompressed, the example binaries will be in an `examples` folder. These are CLIs and can be ran directly from your terminal. 

Each example has a Dockerfile if you don't have Go installed (and have Docker installed). To build an image, run `docker build -t <image-name> .` from an example directory. For example, from the `examples/connect` directory run `docker build -t conngoqtt .`. 

There is also a Docker Compose file which will spin up a broker, a subscriber, and a publisher. To run this, from the `examples` directory use:

```bash
docker-compose up --build
```

If you would like to publish more messages you can use the following from a separate terminal:

```bash
 docker run --net examples_goqtt examples_publish -server broker -port 1885 -message "hi"
```

## Logging

*WIP*

`goqtt` uses [zerolog](https://github.com/rs/zerolog) internally for creating structured logs with different levels (such as DEBUG, INFO, and ERROR). 

However, all log statements are **disabled** by default. If you'd like to enable the logging, set the global log level to the desired level. This is done with the `zerolog.SetGlobalLevel()` function and a parameter such as `zerolog.InfoLevel`. See the examples and their documentation for more levels.
