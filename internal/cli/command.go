package cli

// Command 命令接口
type Command interface {
	// Name 返回命令名称（不包含/前缀）
	Name() string

	// Description 返回命令描述
	Description() string

	// Usage 返回命令用法说明
	Usage() string

	// Execute 执行命令
	Execute(args []string) error
}
