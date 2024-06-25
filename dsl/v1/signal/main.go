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

	// 给工作流发送信号
	signalInput := make(map[string]string)
	signalInput["hello"] = "world"

	err = c.SignalWorkflow(context.Background(), "dsl_9a8da51f-e120-4964-9487-e4b0da8452f0", "9a46eea8-5235-45b9-9808-39d1b87bd1e3", "helloworld", signalInput)

	if err != nil {
		log.Fatalln("Unable to signal workflow", err)
	}
}
