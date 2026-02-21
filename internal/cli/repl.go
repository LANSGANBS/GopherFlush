package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type REPL struct {
	registry *Registry
	scanner  *bufio.Scanner
}

func NewREPL(registry *Registry) *REPL {
	return &REPL{
		registry: registry,
		scanner:  bufio.NewScanner(os.Stdin),
	}
}

func (r *REPL) Start() {
	r.printWelcome()

	for {
		fmt.Print(r.styledPrompt())

		if !r.scanner.Scan() {
			break
		}

		line := strings.TrimSpace(r.scanner.Text())
		if line == "" {
			continue
		}

		if err := r.processLine(line); err != nil {
			fmt.Printf("%s错误:%s %v\n", ColorRed, ColorReset, err)
		}
	}
}

func (r *REPL) printWelcome() {
	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", ColorCyan, strings.Repeat("─", width), ColorReset)

	title := "GopherFlush"
	subtitle := "AI 代码质量检测工具"
	version := "v1.0.0"

	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", ColorCyan, strings.Repeat(" ", padding), ColorBold+ColorWhite, title, ColorReset)

	padding = (width - len(subtitle)) / 2
	fmt.Printf("%s%s%s%s%s\n", ColorCyan, strings.Repeat(" ", padding), ColorDim, subtitle, ColorReset)

	padding = (width - len(version)) / 2
	fmt.Printf("%s%s%s%s%s\n", ColorCyan, strings.Repeat(" ", padding), ColorDim, version, ColorReset)

	fmt.Printf("%s%s%s\n", ColorCyan, strings.Repeat("─", width), ColorReset)

	fmt.Println()
	fmt.Printf("  %s•%s 输入 %s/help%s 查看可用命令\n", ColorCyan, ColorReset, ColorCyan, ColorReset)
	fmt.Printf("  %s•%s 输入 %s/exit%s 退出程序\n", ColorCyan, ColorReset, ColorCyan, ColorReset)
	fmt.Println()

	fmt.Printf("%s%s%s\n", ColorDim, strings.Repeat("─", width), ColorReset)
	fmt.Println()
}

func (r *REPL) styledPrompt() string {
	return fmt.Sprintf("%s❯%s ", ColorCyan, ColorReset)
}

func (r *REPL) processLine(line string) error {
	if !strings.HasPrefix(line, "/") {
		return fmt.Errorf("命令必须以 / 开头，例如: /help")
	}

	line = strings.TrimPrefix(line, "/")

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return fmt.Errorf("请输入命令")
	}

	cmdName := parts[0]
	args := parts[1:]

	return r.registry.Execute(cmdName, args)
}
