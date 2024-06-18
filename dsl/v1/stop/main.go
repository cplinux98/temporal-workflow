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

	err = c.TerminateWorkflow(context.Background(), "greeting-workflow", "ab779394-f83d-4938-a899-47859eee3455", "stop")
	if err != nil {
		log.Fatalln("Unable to terminate workflow", err)
	}

	err = c.CancelWorkflow(context.Background(), "greeting-workflow", "ab779394-f83d-4938-a899-47859eee3455")
	if err != nil {
		log.Fatalln("Unable to cancel workflow", err)
	}
}
