package builtin

import (
	"fmt"
	"gopherflush/internal/analyzer"
	"gopherflush/internal/config"
	"gopherflush/internal/reporter"
	"gopherflush/internal/rules"
	rulesBuiltin "gopherflush/internal/rules/builtin"
	"strings"
)

// RunCommand 运行检测命令
type RunCommand struct {
	cfg *config.Config
}

// NewRunCommand 创建运行命令
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
	// 解析参数
	path := "."
	var ruleNames string

	for i, arg := range args {
		if strings.HasPrefix(arg, "--rules=") {
			ruleNames = strings.TrimPrefix(arg, "--rules=")
		} else if i == 0 && !strings.HasPrefix(arg, "--") {
			path = arg
		}
	}

	// 初始化规则注册表
	registry := c.initializeRules(ruleNames)

	// 初始化分析器
	a := analyzer.NewAnalyzer(registry)

	// 执行代码检测
	fmt.Printf("正在检测: %s\n", path)
	report, err := a.Analyze(path)
	if err != nil {
		return fmt.Errorf("检测失败: %w", err)
	}

	// 生成控制台报告
	consoleReporter := reporter.NewConsoleReporter()
	if err := consoleReporter.Generate(report); err != nil {
		return fmt.Errorf("生成控制台报告失败: %w", err)
	}

	// 生成详细报告文件
	outputPath := "gopherflush-report.json"
	jsonReporter := reporter.NewJSONReporter(outputPath)
	if err := jsonReporter.Generate(report); err != nil {
		fmt.Printf("\n警告: 生成详细报告失败: %v\n", err)
	} else {
		fmt.Printf("\n详细报告已保存到: %s\n", outputPath)
	}

	return nil
}

func (c *RunCommand) initializeRules(ruleNames string) *rules.Registry {
	registry := rules.NewRegistry()

	// 解析要运行的规则列表
	var selectedRules map[string]bool
	if ruleNames != "" {
		selectedRules = make(map[string]bool)
		for _, name := range strings.Split(ruleNames, ",") {
			selectedRules[strings.TrimSpace(name)] = true
		}
	}

	// 注册规则
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
