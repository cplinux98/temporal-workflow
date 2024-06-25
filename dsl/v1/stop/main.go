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

	err = c.TerminateWorkflow(context.Background(), "dsl_7910b371-5352-4ef0-b771-df289ddc0150", "cc8f5d56-5390-4837-9bae-84b1e8dc029e", "stop")
	if err != nil {
		log.Fatalln("Unable to terminate workflow", err)
	}

	//err = c.CancelWorkflow(context.Background(), "dsl_7910b371-5352-4ef0-b771-df289ddc0150", "cc8f5d56-5390-4837-9bae-84b1e8dc029e"")
	//if err != nil {
	//	log.Fatalln("Unable to cancel workflow", err)
	//}
}
