package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	v1 "temporal-workflow/dsl/v1"
	"temporal-workflow/dsl/v1/activities"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort:  "192.168.10.99:7233",
		Namespace: "default",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "dsl", worker.Options{})

	w.RegisterWorkflow(v1.SimpleDSLWorkflow)
	w.RegisterActivity(&activities.PythonCommand{})
	w.RegisterActivity(&activities.ShellCommand{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
