package pkg

import "time"

type ModelConfig struct {
	Type    string `mapstructure:"type"`
	Model   string `mapstructure:"model"`
	ApiKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

type PromptConfig struct {
	Main string `mapstructure:"main"`
}

type Config struct {
	Domain string `mapstructure:"domain"`
	Http   struct {
		Addr    string        `mapstructure:"addr"`
		Timeout time.Duration `mapstructure:"timeout"`
		Debug   bool          `mapstructure:"debug"`
	} `mapstructure:"http"`
	LogLevel string `mapstructure:"log_level"`
	Database struct {
		Driver string `mapstructure:"driver"`
		Source string `mapstructure:"source"`
	} `mapstructure:"database"`
	Robot struct {
		Type         string `mapstructure:"type"`
		ClientId     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
	} `mapstructure:"robot"`
	Model  ModelConfig  `mapstructure:"model"`
	Prompt PromptConfig `mapstructure:"prompt"`
}
