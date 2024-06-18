package v1

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// OnError 错误类型
type OnError string

const (
	ContinueErrorOutput   OnError = "continueErrorOutput"   // 继续运行workflow，并把错误消息保留到共享变量
	ContinueRegularOutput OnError = "continueRegularOutput" // 继续运行workflow，并把错误消息转换为正常输出结果存储到共享变量
	StopWorkflow          OnError = "stopWorkflow"          // 出现错误时立刻停止workflow，并将workflow设置为失败状态
)

// Workflow
// @Description: workflow对象，描述了DSL中的结构
type Workflow struct {
	Id          string            `yaml:"id"`          // 唯一标识
	Name        string            `yaml:"name"`        // workflow的名称
	Version     string            `yaml:"version"`     // workflow的版本信息
	Description string            `yaml:"description"` // workflow的描述信息
	Tags        map[string]string `yaml:"tags"`        // workflow的标签
	Variables   map[string]string `yaml:"variables"`   // workflow的公共变量
	Schedule    WorkflowSchedule  `yaml:"schedule"`    // 描述workflow何时启动实例
	Settings    WorkflowSettings  `yaml:"settings"`    // workflow的设置选项
	Nodes       []Node            `yaml:"nodes"`       // Node节点
	Edges       []Edge            `yaml:"edges"`       // 边，描述Node之间的连接关系
}

// WorkflowSchedule
// @Description: 描述workflow如何启动实例，三者可以同时存在，但必须存在其中一个
type WorkflowSchedule struct {
	Cron    Cron    `yaml:"cron,omitempty"`    // 基于cron启动
	Webhook Webhook `yaml:"webhook,omitempty"` // 基于webhook启动
	Event   Event   `yaml:"event,omitempty"`   // 基于kafka中event启动
}

type Cron struct {
	Expression string `yaml:"expression"` // cron表达式
	NextNode   string `yaml:"nextNode"`   // 触发后启动的节点Id
}

type Webhook struct {
	Method   string `yaml:"method"`   // webhook的方法
	URL      string `yaml:"url"`      // webhook的URL
	NextNode string `yaml:"nextNode"` // 激活的节点
}

type Event struct {
	Host      string   `yaml:"host"` //
	Topic     string   `yaml:"topic"`
	Group     string   `yaml:"group"`
	EventType []string `yaml:"eventType"`
	NextNode  string   `yaml:"nextNode"`
}

type Node struct {
	Id                 string                 `yaml:"id"`                 // 唯一标识
	Name               string                 `yaml:"name"`               // node名称
	Description        string                 `yaml:"description"`        // node的描述信息
	Position           []string               `yaml:"position"`           // node在画布上的坐标
	Activity           string                 `yaml:"activity"`           // node运行哪个activity
	ActivityVersion    string                 `yaml:"activityVersion"`    // node运行activity的版本
	ActivityParameters map[string]interface{} `yaml:"activityParameters"` // node运行activity的参数
	RetryOnFail        bool                   `yaml:"retryOnFail"`        // node运行activity时出现错误是否尝试重试
	MaxTries           int64                  `yaml:"maxTries"`           // 最大尝试次数
	WaitBetweenTries   int64                  `yaml:"waitBetweenTries"`   // 尝试的间隔等待时间 单位ms
	AlwaysOutputData   bool                   `yaml:"alwaysOutputData"`   // 是否总是输出数据
	OnError            OnError                `yaml:"onError"`            // node运行出错时如何处理workflow
}

type Edge struct {
	Id     string `yaml:"id"`     // 唯一标识
	Source string `yaml:"source"` // 源节点Id
	Target string `yaml:"target"` // 目标节点Id
	Label  string `yaml:"label"`  // 边上的显示文字
}

type WorkflowSettings struct{}

// ParseDSLFromYamlFile
//
//	@Description: 从Yaml文件中解析DSL
//	@Author cplinux98 2024-06-17 15:07:26
//	@param filename
//	@return *Workflow
//	@return error
func ParseDSLFromYamlFile(filename string) (Workflow, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return Workflow{}, err
	}

	var workflow Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return Workflow{}, err
	}

	return workflow, nil
}

// TopSort
//
//	@Description: 对节点生成拓扑图排序，
//	@Author cplinux98 2024-06-17 17:22:11
//	@param nodes
//	@param edges
//	@return []string
//	@return error
func TopSort(nodes []Node, edges []Edge) ([]string, error) {
	inDegree := make(map[string]int)
	graph := make(map[string][]string)
	var sorted []string

	// 初始化入度和图结构
	for _, node := range nodes {
		inDegree[node.Id] = 0
		graph[node.Id] = []string{}
	}

	// 构建图和入度计数
	for _, edge := range edges {
		graph[edge.Source] = append(graph[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}

	// 初始化队列，将所有入度为0的节点入队
	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// 处理队列中的节点
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		sorted = append(sorted, node)

		// 对邻接节点入度减1，并将入度变为0的节点入队
		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 检查是否存在环
	if len(sorted) != len(nodes) {
		return nil, fmt.Errorf("cycle detected")
	}

	return sorted, nil
}

// DetectParallelExecution
//
//	@Description: 串并行检测
//	@Author cplinux98 2024-06-17 16:14:45
//	@param sortedNodes
//	@param edges
//	@return [][]string
func DetectParallelExecution(sortedNodes []string, edges []Edge) [][]string {
	nodeLevels := make(map[string]int)
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// 初始化图和入度
	for _, node := range sortedNodes {
		graph[node] = []string{}
		inDegree[node] = 0
	}

	for _, edge := range edges {
		graph[edge.Source] = append(graph[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}

	// 计算每个节点的层级
	level := 0
	var queue []string
	for _, node := range sortedNodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
			nodeLevels[node] = level
		}
	}

	for len(queue) > 0 {
		level++
		nextQueue := []string{}
		for _, node := range queue {
			for _, neighbor := range graph[node] {
				inDegree[neighbor]--
				if inDegree[neighbor] == 0 {
					nextQueue = append(nextQueue, neighbor)
					nodeLevels[neighbor] = level
				}
			}
		}
		queue = nextQueue
	}

	// 按层级分组节点
	levels := make(map[int][]string)
	for node, lvl := range nodeLevels {
		levels[lvl] = append(levels[lvl], node)
	}

	// 将层级映射转换为切片
	var result [][]string
	for i := 0; i < len(levels); i++ {
		result = append(result, levels[i])
	}

	return result
}

// MarkNodesFromStart 标记从指定节点开始的所有子图节点
func MarkNodesFromStart(edges []Edge, start string) map[string]bool {
	if start == "" {
		return nil
	}
	marked := make(map[string]bool)
	var dfs func(string)
	dfs = func(nodeId string) {
		if marked[nodeId] {
			return
		}
		marked[nodeId] = true
		for _, edge := range edges {
			if edge.Source == nodeId {
				dfs(edge.Target)
			}
		}
	}
	dfs(start)
	return marked
}
