package builtin

import (
	"fmt"
	"gopherflush/internal/cli"
	"strings"
)

type HelpCommand struct {
	registry *cli.Registry
}

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
	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	title := "可用命令"
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", cli.ColorCyan, strings.Repeat(" ", padding), cli.ColorBold+cli.ColorWhite, title, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	commands := c.registry.List()
	for _, cmd := range commands {
		fmt.Printf("  %s/%-10s%s %s\n", cli.ColorCyan, cmd.Name(), cli.ColorReset, cmd.Description())
	}

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	fmt.Println()
	fmt.Printf("  %s提示:%s 使用 %s/show <命令名>%s 查看命令详细用法\n", cli.ColorDim, cli.ColorReset, cli.ColorCyan, cli.ColorReset)
	fmt.Println()

	return nil
}
