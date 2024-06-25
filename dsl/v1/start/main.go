package main

import (
	"context"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"log"
	v1 "temporal-workflow/dsl/v1"
	"temporal-workflow/dsl/v1/dslDefinition"
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

	// 读取DSL中workflow
	dslWorkflow, err := dslDefinition.ParseDSLFromYamlFile("/Users/cp/SynologyDrive/MyProject/temporal-workflow/dsl/v1/test_01.yaml")
	if err != nil {
		log.Fatalln("Unable to parse workflow definition", err)
	}

	// 启动workflow
	options := client.StartWorkflowOptions{
		ID:        "dsl_" + uuid.New(),
		TaskQueue: "dsl",
	}

	// 启动workflow
	workflowRun, err := c.ExecuteWorkflow(context.Background(), options, v1.SimpleDSLWorkflow, dslWorkflow, "cron_start_node")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", workflowRun.GetID(), workflowRun.GetRunID())
}
