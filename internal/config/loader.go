package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Loader 配置加载器
type Loader struct{}

// NewLoader 创建配置加载器
func NewLoader() *Loader {
	return &Loader{}
}

// Load 从文件加载配置
func (l *Loader) Load(path string) (*Config, error) {
	// TODO: 实现配置加载逻辑
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadDefault 加载默认配置
func (l *Loader) LoadDefault() *Config {
	return &Config{
		Rules: RulesConfig{
			FileSize: FileSizeConfig{
				Enabled:  true,
				MaxLines: 800,
			},
			FunctionSize: FunctionSizeConfig{
				Enabled:  true,
				MaxLines: 200,
			},
			GlobalVars: GlobalVarsConfig{
				Enabled: true,
			},
			Duplicates: DuplicatesConfig{
				Enabled: true,
			},
			CommentedCode: CommentedCodeConfig{
				Enabled: true,
			},
			InconsistentComment: InconsistentCommentConfig{
				Enabled: true,
			},
			ResourceLeak: ResourceLeakConfig{
				Enabled: true,
			},
			LooseTyping: LooseTypingConfig{
				Enabled: true,
			},
			InaccurateConstant: InaccurateConstantConfig{
				Enabled: true,
			},
			MissingValidation: MissingValidationConfig{
				Enabled: true,
			},
			HardcodedSecrets: HardcodedSecretsConfig{
				Enabled: true,
			},
		},
		Output: OutputConfig{
			Format: "console",
		},
		Exclude: []string{
			"vendor/",
			"node_modules/",
			".git/",
		},
	}
}
