package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// TODO(bootun): 命令行新开个线程，这样可以和AI互动?
func NewShellExecutor() *shellExecutor {
	return &shellExecutor{
		OS: runtime.GOOS,
	}
}

type shellExecutor struct {
	OS string
}

func (s *shellExecutor) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "shell_executor",
		Desc: "command line shell, the user current operating system is " + s.OS,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"command": {
				Desc:     "want to execute command",
				Type:     schema.String,
				Required: true,
			},
		}),
	}, nil
}

func (s *shellExecutor) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	params, err := parseShellParams(argumentsInJSON)
	if err != nil {
		return "", fmt.Errorf("解析命令参数失败: %w", err)
	}

	if strings.TrimSpace(params.Command) == "" {
		return "", fmt.Errorf("命令不能为空")
	}

	var cmd *exec.Cmd

	switch s.OS {
	case "windows":
		cmd = exec.CommandContext(ctx, "cmd", "/C", params.Command)
	case "darwin", "linux":
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", params.Command)
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", s.OS)
	}

	// 捕获标准输出和标准错误
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		log.Printf("执行命令失败: %v, 错误: %v", params.Command, errMsg)
		return "", fmt.Errorf("执行命令失败: %v", errMsg)
	}

	output := stdout.String()
	// log.Printf("命令: %v, 输出: %v", params.Command, output)
	if output == "" {
		return "the command did not return a result", nil
	}
	return output, nil
}

type shellParams struct {
	Command string `json:"command"`
}

func parseShellParams(argumentsInJSON string) (*shellParams, error) {
	var params shellParams
	err := json.Unmarshal([]byte(argumentsInJSON), &params)
	if err != nil {
		return nil, err
	}
	return &params, nil
}
