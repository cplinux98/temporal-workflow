package activities

import (
	"bytes"
	"context"
	"os/exec"
)

type PythonCommand struct {
}

type PythonCommandParams struct {
	Interpreter string `yaml:"interpreter"` // Python解释器路径
	Script      string `yaml:"script"`      // 执行的python脚本内容
}

type PythonCommandResult struct {
	Code   int64  `yaml:"code"`   // 退出状态码
	Stdout string `yaml:"stdout"` // 标准输出
	Stderr string `yaml:"stderr"` // 标准错误输出
}

func (p *PythonCommand) PythonCommandActivity(ctx context.Context, params PythonCommandParams) (*PythonCommandResult, error) {
	cmd := exec.CommandContext(ctx, params.Interpreter, "-c", params.Script)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 运行命令
	err := cmd.Run()

	// 获取退出码
	var exitCode int64 = 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = int64(exitError.ExitCode())
		} else {
			exitCode = int64(-1)
		}
	} else {
		exitCode = 0
	}

	// 构建输出结果
	result := &PythonCommandResult{
		Code:   exitCode,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	return result, nil
}
