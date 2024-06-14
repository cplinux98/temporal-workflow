package v1

import (
	"context"
	"fmt"
)

type Workflow struct {
	Id          string            // 唯一标识
	Name        string            // workflow的名称
	Version     string            // workflow的版本信息
	Description string            // workflow的描述信息
	Tags        map[string]string // workflow的标签
	Variables   map[string]string // workflow的公共变量
	Nodes       []Node            // Node节点
	Settings    WorkflowSettings  // workflow的设置选项
}

type Node struct {
	Id               string                 // 唯一标识
	Name             string                 // node名称
	Description      string                 // node的描述信息
	Position         []string               // node在画布上的坐标
	Type             string                 // node运行的类型
	TypeVersion      string                 // node运行的类型的版本
	TypeParameters   map[string]interface{} // node运行的类型的参数
	RetryOnFail      bool                   // node运行时出现错误是否尝试重试
	MaxTries         int64                  // 最大尝试次数
	WaitBetweenTries int64                  // 尝试的间隔等待时间 单位ms
	AlwaysOutputData bool                   // 是否总是输出数据
	ExecuteOnce      bool                   // 是否在整个workflow生命周期中只运行一次
	OnError          OnError                // node运行出错时如何处理workflow
}

type WorkflowSettings struct{}

// OnError 错误类型
type OnError string

const (
	ContinueErrorOutput   OnError = "continueErrorOutput"   // 继续运行workflow，并把错误消息输出到下一个节点
	ContinueRegularOutput OnError = "continueRegularOutput" // 继续运行workflow，并把错误消息当做成功消息传递到下一个节点
	StopWorkflow          OnError = "stopWorkflow"          // 出现错误时立刻停止workflow，并将workflow设置为失败状态
)

// NodeType 所有的nodeType都要实现这些接口
type NodeType interface {
	Execute(ctx context.Context, parameters map[string]interface{}) (map[string]interface{}, error) // 执行具体逻辑
	Type() string                                                                                   // 获取当前节点类型
}

// CronTrigger 一个触发器类型的节点，通过cron去触发workflow
type CronTrigger struct {
}

func (n *CronTrigger) Execute(ctx context.Context, parameters map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("TriggerNode.Execute")
	// 当计划到期后，就可以执行了
	return nil, nil
}
func (n *CronTrigger) Type() string {
	return "trigger"
}

// WebhookTrigger 一个触发器类型的节点，通过webhook去触发workflow
type WebhookTrigger struct {
}

func (n *WebhookTrigger) Execute(ctx context.Context, parameters map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("WebhookTrigger.Execute")
	// 当接接收到回掉后，就可以执行了
	return nil, nil
}
func (n *WebhookTrigger) Type() string {
	return "trigger"
}

// ShellAction 一个action类型的节点，可以执行shell命令
type ShellAction struct {
}

func (n *ShellAction) Execute(ctx context.Context, parameters map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("ShellAction.Execute")
	return nil, nil
}
func (n *ShellAction) Type() string {
	return "action"
}

// PythonAction 一个action类型的节点，可以执行python脚本
type PythonAction struct {
}

func (n *PythonAction) Execute(ctx context.Context, parameters map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("PythonAction.Execute")
	return nil, nil
}
func (n *PythonAction) Type() string {
	return "action"
}
