package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

var (
	ErrFileNotExist   = errors.New("file not exist")
	ErrInvalidLineArg = errors.New("line param must be in format L{start}-L{end}")
)

// NewFileReader returns a new fileReader instance.
func NewFileReader() *fileReader {
	return &fileReader{}
}

type fileReader struct{}

func (fr *fileReader) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "file_reader",
		Desc: "read file content as string format (supports partial read by line numbers)",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"filename": {
				Desc:     "file name you want to read",
				Type:     schema.String,
				Required: true,
			},
			"line": {
				Desc:     "line range you want to read, format is L{start}-L{end}. For example, L1-L200 means you want to read lines 1-200 (inclusive) of this file. If omitted the whole file is returned.",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

func (fr *fileReader) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	params, err := fr.parseFileReaderParams(argumentsInJSON)
	if err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if strings.TrimSpace(params.Filename) == "" {
		return "", fmt.Errorf("文件名不能为空")
	}

	// If no line param provided – return full file content.
	if strings.TrimSpace(params.Line) == "" {
		file, err := os.ReadFile(params.Filename)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return "", ErrFileNotExist
			}
			return "", err
		}
		return string(file), nil
	}

	// Partial read by line range.
	start, end, err := parseLineRange(params.Line)
	if err != nil {
		return "", err
	}

	// Validate logical range.
	if start <= 0 || end < start {
		return "", fmt.Errorf("无效的行号范围: %d-%d", start, end)
	}

	f, err := os.Open(params.Filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrFileNotExist
		}
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var sb strings.Builder
	current := 0

	for scanner.Scan() {
		current++
		if current < start {
			continue
		}
		if current > end {
			break
		}
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

type fileReaderParams struct {
	Filename string `json:"filename"`
	Line     string `json:"line"`
}

func (fr *fileReader) parseFileReaderParams(argumentsInJSON string) (*fileReaderParams, error) {
	var params fileReaderParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return nil, err
	}
	return &params, nil
}

// parseLineRange converts a string like "L10-L20" to numerical start,end values.
func parseLineRange(arg string) (start, end int, err error) {
	arg = strings.TrimSpace(arg)
	if !strings.HasPrefix(arg, "L") {
		return 0, 0, ErrInvalidLineArg
	}
	parts := strings.SplitN(arg[1:], "-L", 2) // remove first 'L' then split at "-L"
	if len(parts) != 2 {
		return 0, 0, ErrInvalidLineArg
	}
	start, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, ErrInvalidLineArg
	}
	end, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, ErrInvalidLineArg
	}
	return
}
