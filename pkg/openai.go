package pkg

import (
	"errors"
	"github.com/sashabaranov/go-openai"
)

func CreateOpenAiClient(cfg *ModelConfig) (*openai.Client, error) {
	var clientConfig openai.ClientConfig
	if cfg.Type == "openai" {
		clientConfig = openai.DefaultConfig(cfg.ApiKey)
		if cfg.BaseURL != "" {
			clientConfig.BaseURL = cfg.BaseURL
		}
	} else if cfg.Type == "azure_openai" {
		clientConfig = openai.DefaultAzureConfig(cfg.ApiKey, cfg.BaseURL)
		clientConfig.AzureModelMapperFunc = func(model string) string {
			azureModelMapping := map[string]string{
				cfg.Model: cfg.Model,
			}
			return azureModelMapping[model]
		}
	} else {
		return nil, errors.New("invalid type")
	}
	return openai.NewClientWithConfig(clientConfig), nil
}
