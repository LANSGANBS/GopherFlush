package builtin

import (
	"fmt"
	"gopherflush/internal/cli"
	"os"
	"strings"
)

type ExitCommand struct{}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

func (c *ExitCommand) Name() string {
	return "exit"
}

func (c *ExitCommand) Description() string {
	return "退出程序"
}

func (c *ExitCommand) Usage() string {
	return "/exit"
}

func (c *ExitCommand) Execute(args []string) error {
	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	title := "感谢使用 GopherFlush!"
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", cli.ColorCyan, strings.Repeat(" ", padding), cli.ColorBold+cli.ColorGreen, title, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	fmt.Printf("  %s祝您编码愉快，代码质量越来越高！%s\n", cli.ColorDim, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	fmt.Println()
	os.Exit(0)
	return nil
}
