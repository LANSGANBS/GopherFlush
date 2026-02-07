package config

// Config 配置结构
type Config struct {
	Rules   RulesConfig   `yaml:"rules"`
	Output  OutputConfig  `yaml:"output"`
	Exclude []string      `yaml:"exclude"`
}

// RulesConfig 规则配置
type RulesConfig struct {
	FileSize     FileSizeConfig     `yaml:"file_size"`
	FunctionSize FunctionSizeConfig `yaml:"function_size"`
	GlobalVars   GlobalVarsConfig   `yaml:"global_vars"`
	Duplicates   DuplicatesConfig   `yaml:"duplicates"`
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

// OutputConfig 输出配置
type OutputConfig struct {
	Format string `yaml:"format"` // console, json
	Path   string `yaml:"path"`   // JSON输出路径
}
