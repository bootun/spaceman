package common

import (
	"context"
	"fmt"
	"log"

	"github.com/bootun/cosmica/tools"
	"github.com/bootun/cosmica/tools/base"
	"github.com/bootun/cosmica/tools/file"
	"github.com/bootun/cosmica/tools/shell"
	"github.com/bootun/cosmica/utils"
	"github.com/bootun/cosmica/utils/text"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/cloudwego/eino/schema"
)

type SpaceMan struct {
	model   *openai.ChatModel
	toolSet *tools.ToolSet
}

func NewSpaceMan(ctx context.Context, config openai.ChatModelConfig) (*SpaceMan, error) {
	chatModel, err := openai.NewChatModel(ctx, &config)
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
		shell.NewShellExecutor(),
		base.NewBell(),
		file.NewFileReader(),
		file.NewDirReader(),
		bt,
	)
	if err != nil {
		return nil, fmt.Errorf("create tool set: %w", err)
	}
	if err = chatModel.BindTools(ts.Infos()); err != nil {
		return nil, fmt.Errorf("bind tools: %w", err)
	}
	return &SpaceMan{model: chatModel, toolSet: ts}, nil
}

// TODO(bootun): refactor me
func (sm *SpaceMan) HandleQuestion(ctx context.Context, question string, history []*schema.Message) (chatHistory []*schema.Message, err error) {
	if len(history) < 1 {
		chatHistory = []*schema.Message{
			schema.SystemMessage("你是spaceman, 一个严格遵守用户指令，不会偷懒的人工智能，负责规划并解决用户提出的问题。在进行所有行动之前，你需要预先规划为了完成这件事，接下来要做的事情，并告诉用户，然后才行动、调用工具等。"),
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
