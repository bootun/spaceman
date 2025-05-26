package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

var (
	ErrDirNotExist = errors.New("directory not exist")
)

// NewDirReader returns a new dirReader instance.
func NewDirReader() *dirReader {
	return &dirReader{}
}

type dirReader struct{}

func (dr *dirReader) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "dir_reader",
		Desc: "list all files in a directory",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"dirname": {
				Desc:     "directory path you want to read",
				Type:     schema.String,
				Required: true,
			},
			// "recursive": {
			// 	Desc:     "whether to list files recursively in subdirectories",
			// 	Type:     schema.Boolean,
			// 	Required: false,
			// },
		}),
	}, nil
}

func (dr *dirReader) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	params, err := dr.parseDirReaderParams(argumentsInJSON)
	if err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if strings.TrimSpace(params.Dirname) == "" {
		return "", fmt.Errorf("目录名不能为空")
	}

	// Check if directory exists
	info, err := os.Stat(params.Dirname)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrDirNotExist
		}
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("路径 %s 不是一个目录", params.Dirname)
	}

	var fileList []string
	// if params.Recursive {
	// 	// Recursive listing
	// 	err = filepath.Walk(params.Dirname, func(path string, info fs.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		// Skip the root directory itself
	// 		if path != params.Dirname {
	// 			relPath, err := filepath.Rel(params.Dirname, path)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			if info.IsDir() {
	// 				fileList = append(fileList, relPath+"/")
	// 			} else {
	// 				fileList = append(fileList, relPath)
	// 			}
	// 		}
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		return "", err
	// 	}
	// } else

	// Non-recursive listing
	entries, err := os.ReadDir(params.Dirname)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fileList = append(fileList, entry.Name()+"/")
		} else {
			fileList = append(fileList, entry.Name())
		}
	}

	// Convert the file list to JSON
	result, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

type dirReaderParams struct {
	Dirname string `json:"dirname"`
	// Recursive bool   `json:"recursive"`
}

func (dr *dirReader) parseDirReaderParams(argumentsInJSON string) (*dirReaderParams, error) {
	var params dirReaderParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return nil, err
	}
	return &params, nil
}
