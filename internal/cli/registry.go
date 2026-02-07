package cli

import (
	"fmt"
	"sort"
)

// Registry 命令注册中心
type Registry struct {
	commands map[string]Command
}

// NewRegistry 创建新的命令注册中心
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register 注册命令
func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Get 根据名称获取命令
func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// List 获取所有命令列表（按名称排序）
func (r *Registry) List() []Command {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	sort.Strings(names)

	commands := make([]Command, 0, len(r.commands))
	for _, name := range names {
		commands = append(commands, r.commands[name])
	}
	return commands
}

// Execute 执行命令
func (r *Registry) Execute(name string, args []string) error {
	cmd, ok := r.Get(name)
	if !ok {
		return fmt.Errorf("未知命令: %s，使用 /help 查看可用命令", name)
	}
	return cmd.Execute(args)
}
