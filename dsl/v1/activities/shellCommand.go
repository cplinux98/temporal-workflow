package activities

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type ShellCommandV1 struct {
}

type ShellCommandParams struct {
	Command string `yaml:"command"` // 执行命令
}

type ShellCommandResult struct {
	Code   int64  `yaml:"code"`   // 退出状态码
	Stdout string `yaml:"stdout"` // 标准输出
	Stderr string `yaml:"stderr"` // 标准错误输出
}

func (s *ShellCommandV1) ShellCommandV1(ctx context.Context, params *ShellCommandParams) (*ShellCommandResult, error) {
	cmd := exec.Command("/bin/sh", "-c", params.Command)

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
	result := &ShellCommandResult{
		Code:   exitCode,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if exitCode != 0 {
		return result, fmt.Errorf(stderr.String())
	}

	return result, nil
}
