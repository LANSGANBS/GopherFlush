package builtin

import (
	"fmt"
	"gopherflush/internal/analyzer"
	"gopherflush/internal/cli"
	"gopherflush/internal/config"
	"gopherflush/internal/reporter"
	"gopherflush/internal/rules"
	rulesBuiltin "gopherflush/internal/rules/builtin"
	"path/filepath"
	"strings"
)

type RunCommand struct {
	cfg *config.Config
}

func NewRunCommand(cfg *config.Config) *RunCommand {
	return &RunCommand{
		cfg: cfg,
	}
}

func (c *RunCommand) Name() string {
	return "run"
}

func (c *RunCommand) Description() string {
	return "运行代码检测"
}

func (c *RunCommand) Usage() string {
	return "/run [path] [--rules=rule1,rule2]"
}

func (c *RunCommand) Execute(args []string) error {
	path := "."
	var ruleNames string

	for i, arg := range args {
		if strings.HasPrefix(arg, "--rules=") {
			ruleNames = strings.TrimPrefix(arg, "--rules=")
		} else if i == 0 && !strings.HasPrefix(arg, "--") {
			path = arg
		}
	}

	fmt.Println()

	width := 60
	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	title := "代码检测"
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s\n", cli.ColorCyan, strings.Repeat(" ", padding), cli.ColorBold+cli.ColorWhite, title, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	absPath, _ := filepath.Abs(path)
	fmt.Printf("  %s检测路径:%s %s%s%s\n", cli.ColorDim, cli.ColorReset, cli.ColorCyan, absPath, cli.ColorReset)
	fmt.Printf("  %s状态:%s     %s正在分析...%s\n", cli.ColorDim, cli.ColorReset, cli.ColorYellow, cli.ColorReset)

	fmt.Printf("%s%s%s\n", cli.ColorCyan, strings.Repeat("─", width), cli.ColorReset)

	registry := c.initializeRules(ruleNames)

	a := analyzer.NewAnalyzer(registry)

	report, err := a.Analyze(path)
	if err != nil {
		fmt.Println()
		fmt.Printf("  %s✗ 检测失败:%s %s\n", cli.ColorRed, cli.ColorReset, err.Error())
		return fmt.Errorf("检测失败: %w", err)
	}

	consoleReporter := reporter.NewConsoleReporter()
	if err := consoleReporter.Generate(report); err != nil {
		return fmt.Errorf("生成控制台报告失败: %w", err)
	}

	outputPath := "gopherflush-report.json"
	jsonReporter := reporter.NewJSONReporter(outputPath)
	if err := jsonReporter.Generate(report); err != nil {
		fmt.Printf("  %s⚠ 详细报告生成失败:%s %s\n", cli.ColorYellow, cli.ColorReset, err.Error())
	} else {
		absOutputPath, err := filepath.Abs(outputPath)
		if err != nil {
			absOutputPath = outputPath
		}
		fmt.Printf("  %s✓ 详细报告已保存:%s %s\n", cli.ColorGreen, cli.ColorReset, absOutputPath)
	}

	fmt.Println()

	return nil
}

func (c *RunCommand) initializeRules(ruleNames string) *rules.Registry {
	registry := rules.NewRegistry()

	var selectedRules map[string]bool
	if ruleNames != "" {
		selectedRules = make(map[string]bool)
		for _, name := range strings.Split(ruleNames, ",") {
			selectedRules[strings.TrimSpace(name)] = true
		}
	}

	if c.shouldEnableRule("file-size", selectedRules, c.cfg.Rules.FileSize.Enabled) {
		registry.Register(rulesBuiltin.NewFileSizeRule(c.cfg.Rules.FileSize.MaxLines))
	}

	if c.shouldEnableRule("function-size", selectedRules, c.cfg.Rules.FunctionSize.Enabled) {
		registry.Register(rulesBuiltin.NewFunctionSizeRule(c.cfg.Rules.FunctionSize.MaxLines))
	}

	if c.shouldEnableRule("global-vars", selectedRules, c.cfg.Rules.GlobalVars.Enabled) {
		registry.Register(rulesBuiltin.NewGlobalVarsRule())
	}

	if c.shouldEnableRule("duplicates", selectedRules, c.cfg.Rules.Duplicates.Enabled) {
		registry.Register(rulesBuiltin.NewDuplicatesRule())
	}

	return registry
}

func (c *RunCommand) shouldEnableRule(ruleName string, selectedRules map[string]bool, configEnabled bool) bool {
	if selectedRules != nil {
		return selectedRules[ruleName]
	}
	return configEnabled
}
