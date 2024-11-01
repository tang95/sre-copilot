package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai/jsonschema"
	"os/exec"
)

type Bash struct{}

func (b *Bash) Name() string {
	return "bash"
}

func (b *Bash) Description() string {
	return "execute bash command"
}

func (b *Bash) Parameters() any {
	return jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"script": {
				Type:        jsonschema.String,
				Description: "script to execute",
			},
		},
		Required: []string{"script"},
	}
}
func (b *Bash) Call(ctx context.Context, input string) (string, error) {
	// 解析输入为 JSON 对象
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", fmt.Errorf("failed to unmarshal input: %v", err)
	}

	// 获取脚本内容
	script, ok := params["script"].(string)
	if !ok {
		return "", fmt.Errorf("invalid input: script is required")
	}

	// 执行脚本
	out, err := exec.CommandContext(ctx, "bash", "-c", script).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %v", err)
	}

	return string(out), nil
}
