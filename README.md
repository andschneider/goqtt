# goqtt

![Build](https://github.com/andschneider/goqtt/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/andschneider/goqtt)](https://goreportcard.com/report/github.com/andschneider/goqtt)
[![License: MIT](https://img.shields.io/github/license/andschneider/goqtt)](https://img.shields.io/github/license/andschneider/goqtt)

---

## Install

First, install [Go](https://golang.org/doc/install).

Next use go modules to download this package:

`go get github.com/andschneider/goqtt`

## Usage

See the `examples` folder. There are examples for:

- simple connect and disconnect to a broker
- publishing a message a topic
- subscribing to a topic
- a minimal broker

If you have Go installed you can run them with `go run ./examples/<example>`.

Each example has a Dockerfile if you don't have Go installed (and have Docker installed). To build an image, run `docker build -t <image-name> .` from an example directory. For example, from the `examples/connect` directory run `docker build -t conngoqtt .`. 

There is also a Docker Compose file which will spin up a broker, a subscriber, and a publisher. To run this, from the `examples` directory use:

```bash
docker-compose up --build
```

If you would like to publish more messages you can use the following from a separate terminal:

```bash
 docker run --net examples_goqtt examples_publish -server broker -port 1885 -message "hi"
```