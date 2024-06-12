package test

import (
	"github.com/serverlessworkflow/sdk-go/v2/parser"
	"testing"
)

func TestParseWorkflow(t *testing.T) {
	// 测试文件路径
	testFilePath := "test/1.helloworld_example.json"

	// 调用parse
	workflow, err := parser.FromFile(testFilePath)
	if err != nil {
		return
	}
	t.Logf("%+v", workflow)
}
