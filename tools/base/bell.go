package base

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

const (
	FinishFlag = "[finish]"
)

func NewBell() *bell {
	return &bell{}
}

type bell struct{}

func (s *bell) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "bell",
		Desc: `当且仅当出现以下任何一种情况时必须调用:
1.答案已完整给出，对话可结束。
2.已向用户提出问题或澄清请求，需要等待用户回复才能继续。

典型示例:
- Human:你好→ AI:有什么我可以帮你？→ AI:调用 bell（等待用户）。
- 技术解答完成 → 调用 bell（对话结束）。`,
		ParamsOneOf: schema.NewParamsOneOfByParams(nil),
	}, nil
}

func (s *bell) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	log.Printf("任务已完成")
	return FinishFlag, nil
}

// TODO(bootun): 终止原因
type bellParams struct {
	Reason string `json:"reason"`
}
