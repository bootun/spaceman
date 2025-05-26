package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ToolSet struct {
	toolsMap map[string]tool.InvokableTool
	infos    []*schema.ToolInfo
}

func NewToolSet(tools ...tool.InvokableTool) (*ToolSet, error) {
	ts := &ToolSet{}
	ts.toolsMap = make(map[string]tool.InvokableTool, len(tools))
	ts.infos = make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		info, err := t.Info(context.Background())
		if err != nil {
			return nil, fmt.Errorf("get tool info: %w", err)
		}
		toolName := info.Name
		ts.toolsMap[toolName] = t
		ts.infos = append(ts.infos, info)
	}
	return ts, nil
}

func (ts *ToolSet) Infos() []*schema.ToolInfo {
	return ts.infos
}

var (
	ErrToolNotFound = errors.New("tool not found")
)

func (ts *ToolSet) GetTool(name string) (tool.InvokableTool, error) {
	t, ok := ts.toolsMap[name]
	if !ok {
		return nil, ErrToolNotFound
	}
	return t, nil
}

func (ts *ToolSet) AddTool(tool tool.InvokableTool) error {
	info, err := tool.Info(context.Background())
	if err != nil {
		return fmt.Errorf("get tool info: %w", err)
	}
	name := info.Name
	if _, ok := ts.toolsMap[name]; ok {
		return fmt.Errorf("tool %s already exists", name)
	}
	ts.toolsMap[name] = tool
	ts.infos = append(ts.infos, info)
	return nil
}
