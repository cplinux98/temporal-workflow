package main

import (
	"context"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort:  "192.168.10.99:7233",
		Namespace: "default",
	})
	if err != nil {
		log.Fatalln("Unable to connect to temporal server", err)
	}
	defer c.Close()

	err = c.CancelWorkflow(context.Background(), "dsl_69a89bc6-657e-4450-9e3d-110b182e164c", "f9443a84-2819-4a48-8029-e95b236ba113")
	if err != nil {
		log.Fatalln("Unable to cancel workflow", err)
	}
}
