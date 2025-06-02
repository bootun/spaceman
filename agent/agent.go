// 仅用作依赖倒置
package agent

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

type Agent interface {
	// 添加工具
	// AddTools(tools ...tool.InvokableTool) error
	HandleQuestion(ctx context.Context, question string, history []*schema.Message) (chatHistory []*schema.Message, err error)
}

const (
	AgentSpaceman = "spaceman"
)

// CreateAgentFunc 定义创建 agent 的函数类型
type CreateAgentFunc func(ctx context.Context, task string) (Agent, error)

// func NewAgent(ctx context.Context, task string) (Agent, error) {
// 	cfg, err := config.LoadConfig("config.yml")
// 	if err != nil {
// 		return nil, fmt.Errorf("load config: %w", err)
// 	}
// 	spacemanCfg := cfg.Agents.Spaceman
// 	spaceman, err := common.NewSpaceMan(ctx, openai.ChatModelConfig{
// 		BaseURL: spacemanCfg.BaseURL,
// 		APIKey:  spacemanCfg.Token,
// 		Model:   spacemanCfg.ModelID,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("create spaceman agent: %w", err)
// 	}
// 	return spaceman, nil
// }
