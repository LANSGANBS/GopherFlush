package builtin

import (
	"fmt"
	"os"
)

// ExitCommand 退出命令
type ExitCommand struct{}

// NewExitCommand 创建退出命令
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
	fmt.Println("再见！")
	os.Exit(0)
	return nil
}
