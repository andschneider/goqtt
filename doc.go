/*
Package goqtt is a MQTT 3.1.1 client library.

It does not implement the full client specification, see README for more information.

Use an instantiated Client to interact with a MQTT broker:

	client := goqtt.NewClient("example.broker")
	err := client.Connect()
	...
	err := client.Subscribe()
	...
	err := client.Publish("hello world")
	...
	mes, err := client.ReadLoop() // call in a 'for' loop for indefinite reading
	...

It is nice to let the broker know when you are disconnecting:

	err := client.Connect()
	if err != nil {
		// handle error
	}
	defer client.Disconnect()

To configure a Client, the only required parameter is the address of the broker.
There are default values set for the other options. See the Constants section for
their values.

	client := goqtt.NewClient("example.broker")

The other configuration options are: Port, Topic, ClientId, and KeepAlive.
To set these, use their functions and pass to the NewClient:

	p := goqtt.Port("1883")
	t := goqtt.Topic("hello/world")
	cid := goqtt.ClientId("goqtt-clientId")
	ka := goqtt.KeepAlive(15)  // seconds

	client := goqtt.NewClient("example.broker", p, t, cid, ka)

Note that you can pass any of these options, you don't need to pass all of them.

*/
package goqtt
