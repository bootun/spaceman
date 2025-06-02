package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bootun/cosmica/agent"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type agentCreator struct {
	createFunc agent.CreateAgentFunc
}

func NewAgentCreator(createFunc agent.CreateAgentFunc) *agentCreator {
	return &agentCreator{createFunc: createFunc}
}

func (ac *agentCreator) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "create_agent",
		Desc: `create an assistant to help you solve task. 
You can specify the tools that the assistant can use, define the problem it wants to solve, and the assistant will return the final result to you. 
Generally speaking, tasks assigned to assistants should not be too complex, otherwise assistants may not be able to handle the work well. 
If there are really complex tasks, you can try breaking them down into small tasks and assigning each task to an assistant to execute.

this is the assistant list:
- netizen: a netizen who can operate the browser and answer questions, if you want to use the browser, you can create netizen assistant.
`,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"name": {
				Desc:     "the name of the assistant",
				Type:     schema.String,
				Enum:     []string{"netizen"},
				Required: true,
			},
			"task": {
				Desc:     "tasks that the assistant needs to solve",
				Type:     schema.String,
				Required: true,
			},
		}),
	}, nil
}

func (ac *agentCreator) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	param, err := ac.parseParams(argumentsInJSON)
	if err != nil {
		return "", fmt.Errorf("parse params: %w", err)
	}

	agent, err := ac.createFunc(ctx, param.Task)
	if err != nil {
		return "", fmt.Errorf("create agent: %w", err)
	}
	// ts, err := tools.NewToolSet(
	// 	shell.NewShellExecutor(),
	// 	base.NewBell(),
	// 	file.NewFileReader(),
	// 	file.NewDirReader(),
	// )
	// if err != nil {
	// 	return "", fmt.Errorf("create tool set: %w", err)
	// }
	// if err = agent.AddTools(ts); err != nil {
	// 	return "", fmt.Errorf("bind tools: %w", err)
	// }
	res, err := agent.HandleQuestion(ctx, param.Task, []*schema.Message{
		schema.SystemMessage("You are an assistant to help solve the task, after completing all tasks, you need to summarize the content and results of the tasks and call the tool to end the conversation"),
	})
	if err != nil {
		return "", fmt.Errorf("handle question: %w", err)
	}
	log.Printf("res: %v", res)
	panic("not implemented") // TODO: Implement this function and return the appropriate dat
	return "", nil
}

type createAgentParams struct {
	Task string `json:"task"`
}

func (ac *agentCreator) parseParams(argumentsInJSON string) (*createAgentParams, error) {
	var param createAgentParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &param); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return &param, nil
}
