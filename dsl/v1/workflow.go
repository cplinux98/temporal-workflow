package v1

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
)

// ListenWorkflow 主workflow definition，用于监听事件，然后通过workflow定义启动子workflow来运行workflow实例
func ListenWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	bindings := make(map[string]string)
	//workflowcheck:ignore Only iterates for building another map
	for k, v := range dslWorkflow.Variables {
		bindings[k] = v
	}

	// 监听事件
	channel := workflow.GetSignalChannel(ctx, "event_signal")
	var eventPayload string
	channel.Receive(ctx, &eventPayload)

	// 根据事件创建子workflow来解析运行DSL
	opt := workflow.ChildWorkflowOptions{}
	ctx = workflow.WithChildOptions(ctx, opt)
	err := workflow.ExecuteChildWorkflow(ctx, ExecuteWorkflow, eventPayload)
	if err != nil {
		return nil, nil
	}

	//// 在这里使用childrenWorkflow
	//// 使用workflow的settings配置workflow的options
	//
	//ao := workflow.ActivityOptions{
	//	StartToCloseTimeout: 10 * time.Second,
	//}
	//ctx = workflow.WithActivityOptions(ctx, ao)
	//logger := workflow.GetLogger(ctx)
	//
	////err := dslWorkflow.Root.execute(ctx, bindings)
	////if err != nil {
	////	logger.Error("DSL Workflow failed.", "Error", err)
	////	return nil, err
	////}
	//
	//logger.Info("DSL Workflow completed.")
	return nil, nil
}

// 子 Workflow 定义
func ExecuteWorkflow(ctx workflow.Context, payload string) error {
	// 解析 DSL 并执行相应的任务
	// 这里只是一个简单示例，实际实现需要根据 DSL 内容解析和执行
	workflow.ExecuteActivity(ctx, TaskActivity, payload).Get(ctx, nil)
	return nil
}

// 活动定义
func TaskActivity(ctx context.Context, payload string) error {
	// 执行具体任务
	// 这里只是一个简单示例
	fmt.Println("Executing task with payload:", payload)
	return nil
}
