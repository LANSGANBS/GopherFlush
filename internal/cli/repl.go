package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// REPL 交互式命令行
type REPL struct {
	registry *Registry
	scanner  *bufio.Scanner
}

// NewREPL 创建新的 REPL
func NewREPL(registry *Registry) *REPL {
	return &REPL{
		registry: registry,
		scanner:  bufio.NewScanner(os.Stdin),
	}
}

// Start 启动交互式命令行
func (r *REPL) Start() {
	fmt.Println("GopherFlush - AI代码质量检测工具")
	fmt.Println("========================================")
	fmt.Println("输入 /help 查看可用命令")
	fmt.Println("输入 /exit 退出程序")
	fmt.Println()

	for {
		fmt.Print("gopherflush> ")

		if !r.scanner.Scan() {
			break
		}

		line := strings.TrimSpace(r.scanner.Text())
		if line == "" {
			continue
		}

		if err := r.processLine(line); err != nil {
			fmt.Printf("错误: %v\n", err)
		}
	}
}

// processLine 处理输入行
func (r *REPL) processLine(line string) error {
	// 检查是否是命令（以 / 开头）
	if !strings.HasPrefix(line, "/") {
		return fmt.Errorf("命令必须以 / 开头，例如: /help")
	}

	// 移除 / 前缀
	line = strings.TrimPrefix(line, "/")

	// 解析命令和参数
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return fmt.Errorf("请输入命令")
	}

	cmdName := parts[0]
	args := parts[1:]

	// 执行命令
	return r.registry.Execute(cmdName, args)
}
