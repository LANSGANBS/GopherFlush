package builtin

import (
	"flag"
	"fmt"
	"strings"
)

// ShowCommand 显示信息命令
type ShowCommand struct {
	flagSet *flag.FlagSet
}

// NewShowCommand 创建显示命令
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
	fmt.Println("\n可用的命令行参数:")
	fmt.Println("========================================")

	c.flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  -%s\n", f.Name)
		fmt.Printf("    说明: %s\n", f.Usage)
		if f.DefValue != "" {
			fmt.Printf("    默认值: %s\n", f.DefValue)
		}
		fmt.Println()
	})

	return nil
}

func (c *ShowCommand) showRules() error {
	fmt.Println("\n可用的检测规则:")
	fmt.Println("========================================")
	fmt.Println("  file-size")
	fmt.Println("    检测文件是否超过指定行数（默认800行）")
	fmt.Println()
	fmt.Println("  function-size")
	fmt.Println("    检测函数是否超过指定行数（默认200行）")
	fmt.Println()
	fmt.Println("  global-vars")
	fmt.Println("    检测全局变量滥用")
	fmt.Println()
	fmt.Println("  duplicates")
	fmt.Println("    检测重复的变量或函数声明")
	fmt.Println()

	return nil
}

func (c *ShowCommand) showConfig() error {
	fmt.Println("\n配置文件格式:")
	fmt.Println("========================================")
	fmt.Println("配置文件使用 YAML 格式，示例:")
	fmt.Println()
	fmt.Println("rules:")
	fmt.Println("  file_size:")
	fmt.Println("    enabled: true")
	fmt.Println("    max_lines: 800")
	fmt.Println("  function_size:")
	fmt.Println("    enabled: true")
	fmt.Println("    max_lines: 200")
	fmt.Println()

	return nil
}
