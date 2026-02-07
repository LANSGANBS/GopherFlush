package builtin

import (
	"fmt"
	"gopherflush/internal/cli"
)

// HelpCommand 帮助命令
type HelpCommand struct {
	registry *cli.Registry
}

// NewHelpCommand 创建帮助命令
func NewHelpCommand(registry *cli.Registry) *HelpCommand {
	return &HelpCommand{
		registry: registry,
	}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "显示所有可用命令"
}

func (c *HelpCommand) Usage() string {
	return "/help"
}

func (c *HelpCommand) Execute(args []string) error {
	fmt.Println("\n可用命令:")
	fmt.Println("========================================")

	commands := c.registry.List()
	for _, cmd := range commands {
		fmt.Printf("  /%s\n", cmd.Name())
		fmt.Printf("    %s\n", cmd.Description())
		fmt.Printf("    用法: %s\n\n", cmd.Usage())
	}

	return nil
}
