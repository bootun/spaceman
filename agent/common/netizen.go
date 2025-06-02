package common

import (
	"context"
	"fmt"
	"log"

	"github.com/bootun/cosmica/agent"
	"github.com/bootun/cosmica/config"
	"github.com/bootun/cosmica/tools"
	"github.com/bootun/cosmica/tools/base"
	"github.com/bootun/cosmica/utils"
	"github.com/bootun/cosmica/utils/text"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/cloudwego/eino/schema"
)

type Netizen struct {
	model   *openai.ChatModel
	toolSet *tools.ToolSet
}

func NewNetizen(ctx context.Context, task string) (agent.Agent, error) {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	spacemanCfg := cfg.Agents.Spaceman
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: spacemanCfg.BaseURL,
		APIKey:  spacemanCfg.Token,
		Model:   spacemanCfg.ModelID,
	})
	if err != nil {
		return nil, fmt.Errorf("create chat model: %w", err)
	}
	bt, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{
		Headless: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create browser use tool: %w", err)
	}
	// 为AI配置工具集
	ts, err := tools.NewToolSet(
		base.NewBell(),
		bt,
	)
	if err != nil {
		return nil, fmt.Errorf("create tool set: %w", err)
	}
	if err = chatModel.BindTools(ts.Infos()); err != nil {
		return nil, fmt.Errorf("bind tools: %w", err)
	}
	return &Netizen{model: chatModel, toolSet: ts}, nil
}

// TODO(bootun): refactor me
func (sm *Netizen) HandleQuestion(ctx context.Context, question string, history []*schema.Message) (chatHistory []*schema.Message, err error) {
	if len(history) < 1 {
		chatHistory = []*schema.Message{
			schema.SystemMessage("你是netizen, 一个严格遵守用户指令，不会偷懒的人工智能。你擅长使用浏览器从网络上获取知识、进行操作"),
		}
	} else {
		chatHistory = history
	}
	chatHistory = append(chatHistory, schema.UserMessage(question))

	maxRetry := 10
	i := 0
	finished := false
	for {
		if i > maxRetry || finished {
			break
		}
		i++
		// 生成回答
		stream, err := sm.model.Stream(ctx, chatHistory)
		if err != nil {
			return chatHistory, fmt.Errorf("chat with stream: %w", err)
		}
		msg, err := utils.DealStream(stream, func(msg *schema.Message) {
			fmt.Print(msg.Content)
		})
		if err != nil {
			return chatHistory, fmt.Errorf("deal message: %w", err)
		}
		fmt.Println()
		chatHistory = append(chatHistory, schema.AssistantMessage(msg.Content, msg.ToolCalls))

		if len(msg.ToolCalls) > 0 {
			// 工具调用
			for _, toolCall := range msg.ToolCalls {
				toolName := toolCall.Function.Name
				toolParams := toolCall.Function.Arguments

				fmt.Println(text.Colorize(fmt.Sprintf("<tool call: %s, args: %v>", toolName, toolParams), text.Black, text.BgYellow))
				// 获取工具

				t, err := sm.toolSet.GetTool(toolName)
				if err != nil {
					log.Printf("获取%s工具调用出现错误: %v, 参数: %v", toolName, err, toolParams)
					chatHistory = append(chatHistory, schema.ToolMessage(fmt.Sprintf("调用工具出现了错误: %v", err), toolCall.ID))
					continue
				}
				// 调用工具
				res, err := t.InvokableRun(ctx, toolParams)
				if err != nil {
					log.Printf("调用%s工具时出现了错误: %v, 参数: %v", toolName, err, toolParams)
					chatHistory = append(chatHistory, schema.ToolMessage(fmt.Sprintf("调用工具出现了错误: %v", err), toolCall.ID))
					continue
				}
				chatHistory = append(chatHistory, schema.ToolMessage(res, toolCall.ID))
				if res == base.FinishFlag {
					finished = true
					break
				}
			}
		}
	}
	return
}

// // AddTool implements agent.Agent.
// func (sm *SpaceMan) AddTools(tools ...tool.InvokableTool) error {
// 	for _, t := range tools {
// 		if err := sm.toolSet.AddTool(t); err!= nil {
// 			return fmt.Errorf("add tool: %w", err)
// 		}
// 	}
// 	sm.model.BindTools(sm.toolSet.Infos())
// 	return nil
// }
