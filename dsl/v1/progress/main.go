package main

import (
	"context"
	"fmt"
	"go.temporal.io/api/enums/v1"
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

	// 获取工作流的进度
	//encodeValue, err := c.QueryWorkflow(context.Background(), "dsl_095fcebc-6623-4bc3-ac1e-87ddcd262d4d", "974ce869-42e5-4768-bb9a-5fb4a7449757", "SimpleDSLWorkflow")
	//if err != nil {
	//	log.Fatalln("Unable to query workflow", err)
	//}
	//fmt.Println(encodeValue)
	//workflowRun := c.GetWorkflow(context.Background(), "dsl_095fcebc-6623-4bc3-ac1e-87ddcd262d4d", "974ce869-42e5-4768-bb9a-5fb4a7449757")
	//fmt.Println(workflowRun.GetWithOptions())
	//client.WorkflowRunGetOptions{}
	iter := c.GetWorkflowHistory(context.Background(), "dsl_80934566-a0f9-4543-bd86-8b31dd7b4ecb", "b4b777be-291d-44b3-ba15-779439d0aa31", false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)

	//events := []*shared.HistoryEvent{}
	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			log.Fatalf("Failed to get next history event: %v", err)
		}

		// 解析事件以确定当前节点或状态
		fmt.Printf("Event: %v\n", event)
	}

	//encodeValue, err := c.QueryWorkflow(context.Background(), "dsl_7910b371-5352-4ef0-b771-df289ddc0150", "cc8f5d56-5390-4837-9bae-84b1e8dc029e", "current_state")
	//if err != nil {
	//	log.Fatalln("Unable to query workflow", err)
	//}
	//
	//var state v1.WorkflowState
	//if err := encodeValue.Get(&state); err != nil {
	//	log.Fatalln("Unable to get workflow state", err)
	//}
	//
	//fmt.Printf("State: %+v\n", state)
}
