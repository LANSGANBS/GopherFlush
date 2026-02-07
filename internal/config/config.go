package config

// Config 配置结构
type Config struct {
	Rules   RulesConfig   `yaml:"rules"`
	Output  OutputConfig  `yaml:"output"`
	Exclude []string      `yaml:"exclude"`
}

// RulesConfig 规则配置
type RulesConfig struct {
	FileSize            FileSizeConfig            `yaml:"file_size"`
	FunctionSize        FunctionSizeConfig        `yaml:"function_size"`
	GlobalVars          GlobalVarsConfig          `yaml:"global_vars"`
	Duplicates          DuplicatesConfig          `yaml:"duplicates"`
	CommentedCode       CommentedCodeConfig       `yaml:"commented_code"`
	InconsistentComment InconsistentCommentConfig `yaml:"inconsistent_comment"`
	ResourceLeak        ResourceLeakConfig        `yaml:"resource_leak"`
	LooseTyping         LooseTypingConfig         `yaml:"loose_typing"`
	InaccurateConstant  InaccurateConstantConfig  `yaml:"inaccurate_constant"`
	MissingValidation   MissingValidationConfig   `yaml:"missing_validation"`
	HardcodedSecrets    HardcodedSecretsConfig    `yaml:"hardcoded_secrets"`
}

// FileSizeConfig 文件大小规则配置
type FileSizeConfig struct {
	Enabled  bool `yaml:"enabled"`
	MaxLines int  `yaml:"max_lines"`
}

// FunctionSizeConfig 函数大小规则配置
type FunctionSizeConfig struct {
	Enabled  bool `yaml:"enabled"`
	MaxLines int  `yaml:"max_lines"`
}

// GlobalVarsConfig 全局变量规则配置
type GlobalVarsConfig struct {
	Enabled bool `yaml:"enabled"`
}

// DuplicatesConfig 重复代码规则配置
type DuplicatesConfig struct {
	Enabled bool `yaml:"enabled"`
}

// CommentedCodeConfig 注释代码规则配置
type CommentedCodeConfig struct {
	Enabled bool `yaml:"enabled"`
}

// InconsistentCommentConfig 不一致注释规则配置
type InconsistentCommentConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ResourceLeakConfig 资源泄漏规则配置
type ResourceLeakConfig struct {
	Enabled bool `yaml:"enabled"`
}

// LooseTypingConfig 宽泛类型规则配置
type LooseTypingConfig struct {
	Enabled bool `yaml:"enabled"`
}

// InaccurateConstantConfig 不准确常量规则配置
type InaccurateConstantConfig struct {
	Enabled bool `yaml:"enabled"`
}

// MissingValidationConfig 缺少验证规则配置
type MissingValidationConfig struct {
	Enabled bool `yaml:"enabled"`
}

// HardcodedSecretsConfig 硬编码配置规则配置
type HardcodedSecretsConfig struct {
	Enabled bool `yaml:"enabled"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	Format string `yaml:"format"` // console, json
	Path   string `yaml:"path"`   // JSON输出路径
}
