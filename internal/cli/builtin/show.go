package builtin

import (
	"flag"
	"fmt"
	"gopherflush/internal/cli"
	"strings"
)

type ShowCommand struct {
	flagSet *flag.FlagSet
}

func NewShowCommand(flagSet *flag.FlagSet) *ShowCommand {
	return &ShowCommand{
		flagSet: flagSet,
	}
}

func (c *ShowCommand) Name() string {
	return "show"
}

func (c *ShowCommand) Description() string {
	return "显示各种信息（flags/rules/config）"
}

func (c *ShowCommand) Usage() string {
	return "/show <flags|rules|config>"
}

func (c *ShowCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要显示的内容: flags, rules, config")
	}

	switch strings.ToLower(args[0]) {
	case "flags":
		return c.showFlags()
	case "rules":
		return c.showRules()
	case "config":
		return c.showConfig()
	default:
		return fmt.Errorf("未知选项: %s，可用选项: flags, rules, config", args[0])
	}
}

func (c *ShowCommand) showFlags() error {
	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	title := "命令行参数"
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", cli.ColorCyan, strings.Repeat(" ", padding), cli.ColorBold+cli.ColorWhite, title, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	c.flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  %s-%-12s%s %s\n", cli.ColorCyan, f.Name, cli.ColorReset, f.Usage)
		if f.DefValue != "" {
			fmt.Printf("              %s默认值:%s %s\n", cli.ColorDim, cli.ColorReset, cli.ColorGreen+f.DefValue+cli.ColorReset)
		}
	})

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)
	fmt.Println()

	return nil
}

func (c *ShowCommand) showRules() error {
	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	title := "检测规则"
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", cli.ColorCyan, strings.Repeat(" ", padding), cli.ColorBold+cli.ColorWhite, title, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	rules := []struct {
		name        string
		description string
	}{
		{"file-size", "检测文件是否超过指定行数（默认800行）"},
		{"function-size", "检测函数是否超过指定行数（默认200行）"},
		{"global-vars", "检测全局变量滥用"},
		{"duplicates", "检测重复的变量或函数声明"},
		{"commented-code", "检测注释掉的代码"},
		{"inconsistent-comment", "检测不一致的注释风格"},
		{"resource-leak", "检测资源泄漏风险"},
		{"loose-typing", "检测宽泛的类型定义"},
		{"inaccurate-constant", "检测不准确的常量值"},
		{"missing-validation", "检测缺少输入验证"},
		{"hardcoded-secrets", "检测硬编码的敏感信息"},
	}

	for _, rule := range rules {
		fmt.Printf("  %s%-16s%s %s\n", cli.ColorGreen, rule.name, cli.ColorReset, rule.description)
	}

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)
	fmt.Println()

	return nil
}

func (c *ShowCommand) showConfig() error {
	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	title := "配置文件格式"
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", cli.ColorCyan, strings.Repeat(" ", padding), cli.ColorBold+cli.ColorWhite, title, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	fmt.Printf("  %s配置文件使用 YAML 格式，示例:%s\n", cli.ColorDim, cli.ColorReset)
	fmt.Println()
	fmt.Println("  rules:")
	fmt.Println("    file_size:")
	fmt.Println("      enabled: true")
	fmt.Println("      max_lines: 800")
	fmt.Println("    function_size:")
	fmt.Println("      enabled: true")
	fmt.Println("      max_lines: 200")

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)
	fmt.Println()

	return nil
}
