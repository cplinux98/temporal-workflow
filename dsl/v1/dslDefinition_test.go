package v1

import (
	"reflect"
	"testing"
)

func TestParseDSLFromYamlFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    Workflow
		wantErr bool
	}{
		{name: "test_01", args: args{filename: "/Users/cp/SynologyDrive/MyProject/temporal-workflow/dsl/v1/helloworld_test.yaml"}, want: Workflow{
			Id:          "helloworld_test",
			Name:        "示例workflow",
			Version:     "1.0",
			Description: "用于测试",
			Tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			Variables: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			Schedule: WorkflowSchedule{
				Cron: Cron{
					Expression: "*/5 * * * *",
					NextNode:   "cron_start_node",
				},
				Webhook: Webhook{
					Method:   "get",
					URL:      "/helloworld",
					NextNode: "webhook_start_node",
				},
				Event: Event{
					Host:      "127.0.0.1:9092",
					Topic:     "test_topic",
					Group:     "test_group",
					EventType: []string{"test_event"},
					NextNode:  "event_start_node",
				},
			},
			Settings: WorkflowSettings{},
			Nodes: []Node{
				{
					Id:              "cron_start_node",
					Name:            "计划任务启动的node",
					Description:     "当计划任务到期时触发的node",
					Position:        []string{"100", "200"},
					Activity:        "ShellCommandActivity",
					ActivityVersion: "1",
					ActivityParameters: map[string]interface{}{
						"command": "echo 'hello cron'",
					},
					RetryOnFail:      false,
					MaxTries:         0,
					WaitBetweenTries: 0,
					AlwaysOutputData: true,
					OnError:          StopWorkflow,
				},
				{
					Id:              "webhook_start_node",
					Name:            "webhook启动的node",
					Description:     "webhook启动的node",
					Position:        []string{"300", "400"},
					Activity:        "ShellCommandActivity",
					ActivityVersion: "1",
					ActivityParameters: map[string]interface{}{
						"command": "echo 'hello webhook'",
					},
					RetryOnFail:      false,
					MaxTries:         0,
					WaitBetweenTries: 0,
					AlwaysOutputData: true,
					OnError:          StopWorkflow,
				},
				{
					Id:              "event_start_node",
					Name:            "event启动的node",
					Description:     "event启动的node",
					Position:        []string{"400", "500"},
					Activity:        "ShellCommandActivity",
					ActivityVersion: "1",
					ActivityParameters: map[string]interface{}{
						"command": "echo 'hello event'",
					},
					RetryOnFail:      false,
					MaxTries:         0,
					WaitBetweenTries: 0,
					AlwaysOutputData: true,
					OnError:          StopWorkflow,
				},
				{
					Id:              "simple_node",
					Name:            "普通node节点",
					Description:     "普通node节点",
					Position:        []string{"500", "600"},
					Activity:        "ShellCommandActivity",
					ActivityVersion: "1",
					ActivityParameters: map[string]interface{}{
						"command": "echo 'hello simple'",
					},
					RetryOnFail:      false,
					MaxTries:         0,
					WaitBetweenTries: 0,
					AlwaysOutputData: true,
					OnError:          StopWorkflow,
				},
			},
			Edges: []Edge{
				{
					Id:     "cron_start_node-simple_node",
					Source: "cron_start_node",
					Target: "simple_node",
				},
				{
					Id:     "webhook_start_node-simple_node",
					Source: "webhook_start_node",
					Target: "simple_node",
				},
				{
					Id:     "event_start_node-simple_node",
					Source: "event_start_node",
					Target: "simple_node",
				},
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDSLFromYamlFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDSLFromYamlFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDSLFromYamlFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopSort(t *testing.T) {
	type args struct {
		nodes []Node
		edges []Edge
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "no cycle", args: args{
			nodes: []Node{
				{Id: "01", Name: "01node"},
				{Id: "02", Name: "02node"},
				{Id: "03", Name: "03node"},
				{Id: "04", Name: "04node"},
			},
			edges: []Edge{
				{Id: "01_02", Source: "01", Target: "02"},
				{Id: "01_03", Source: "01", Target: "03"},
				{Id: "03_04", Source: "03", Target: "04"},
			}}, want: []string{"01", "02", "03", "04"}, wantErr: false},
		{name: "cycle_1", args: args{
			nodes: []Node{
				{Id: "01", Name: "01node"},
				{Id: "02", Name: "02node"},
				{Id: "03", Name: "03node"},
				{Id: "04", Name: "04node"},
			},
			edges: []Edge{
				{Id: "01_02", Source: "01", Target: "02"},
				{Id: "01_03", Source: "01", Target: "03"},
				{Id: "03_04", Source: "03", Target: "04"},
				{Id: "04_01", Source: "04", Target: "01"},
			}}, want: nil, wantErr: true},
		{name: "cycle_2", args: args{
			nodes: []Node{
				{Id: "01", Name: "01node"},
				{Id: "02", Name: "02node"},
				{Id: "03", Name: "03node"},
				{Id: "04", Name: "04node"},
			},
			edges: []Edge{
				{Id: "01_02", Source: "01", Target: "02"},
				{Id: "01_03", Source: "01", Target: "03"},
				{Id: "03_04", Source: "03", Target: "04"},
				{Id: "02_01", Source: "02", Target: "01"},
			}}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TopSort(tt.args.nodes, tt.args.edges)
			if (err != nil) != tt.wantErr {
				t.Errorf("TopSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopSort() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectParallelExecution(t *testing.T) {
	type args struct {
		sortedNodes []string
		edges       []Edge
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{name: "parallel", args: args{
			sortedNodes: []string{"01", "02", "03", "04"},
			edges: []Edge{
				{Id: "01_02", Source: "01", Target: "02"},
				{Id: "01_03", Source: "01", Target: "03"},
				{Id: "03_04", Source: "03", Target: "04"},
			},
		}, want: [][]string{{"01"}, {"02", "03"}, {"04"}}},
		{name: "sequence", args: args{
			sortedNodes: []string{"01", "02", "03", "04"},
			edges: []Edge{
				{Id: "01_02", Source: "01", Target: "02"},
				{Id: "02_03", Source: "02", Target: "03"},
				{Id: "03_04", Source: "03", Target: "04"},
			},
		}, want: [][]string{{"01"}, {"02"}, {"03"}, {"04"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectParallelExecution(tt.args.sortedNodes, tt.args.edges); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectParallelExecution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkNodesFromStart(t *testing.T) {
	type args struct {
		edges []Edge
		start string
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{name: "start_02", args: args{edges: []Edge{
			{Id: "01_03", Source: "01", Target: "03"},
			{Id: "02_03", Source: "02", Target: "03"},
			{Id: "03_04", Source: "03", Target: "04"},
			{Id: "04_05", Source: "04", Target: "05"},
		}, start: "02"}, want: map[string]bool{
			"02": true,
			"03": true,
			"04": true,
			"05": true,
		}},
		{name: "start_01", args: args{edges: []Edge{
			{Id: "01_03", Source: "01", Target: "03"},
			{Id: "02_03", Source: "02", Target: "03"},
			{Id: "03_04", Source: "03", Target: "04"},
			{Id: "04_05", Source: "04", Target: "05"},
		}, start: "01"}, want: map[string]bool{
			"01": true,
			"03": true,
			"04": true,
			"05": true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarkNodesFromStart(tt.args.edges, tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkNodesFromStart() = %v, want %v", got, tt.want)
			}
		})
	}
}
