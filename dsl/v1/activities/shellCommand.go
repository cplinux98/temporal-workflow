package activities

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type ShellCommand struct {
}

type ShellCommandParams struct {
	Command string `yaml:"command"` // 执行命令
}

type ShellCommandResult struct {
	Code   int64  `yaml:"code"`   // 退出状态码
	Stdout string `yaml:"stdout"` // 标准输出
	Stderr string `yaml:"stderr"` // 标准错误输出
}

// ShellCommandV1
//
//	@Description: 运行shell命令 v1版本
//	@Author cplinux98 2024-06-19 13:45:12
//	@receiver s
//	@param ctx
//	@param params
//	@return *ShellCommandResult
//	@return error
func (s *ShellCommand) ShellCommandV1(ctx context.Context, params *ShellCommandParams) (*ShellCommandResult, error) {
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
