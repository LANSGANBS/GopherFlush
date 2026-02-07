package main

import (
	"flag"
	"fmt"
	"gopherflush/internal/analyzer"
	"gopherflush/internal/cli"
	cliBuiltin "gopherflush/internal/cli/builtin"
	"gopherflush/internal/config"
	"gopherflush/internal/reporter"
	"gopherflush/internal/rules"
	"gopherflush/internal/rules/builtin"
	"os"
	"strings"
)

func main() {
	// 定义命令行参数
	var (
		interactive = flag.Bool("interactive", false, "启动交互模式")
		path        = flag.String("path", ".", "要检测的代码路径")
		configFile  = flag.String("config", "", "配置文件路径")
		ruleNames   = flag.String("rules", "", "要运行的规则（逗号分隔，留空表示运行所有规则）")
		format      = flag.String("format", "console", "输出格式（console/json）")
		outputPath  = flag.String("output", "gopherflush-report.json", "详细报告输出路径")
	)

	flag.Parse()

	// 加载配置
	cfg := loadConfig(*configFile)

	// 如果没有提供任何参数或明确指定交互模式，进入交互模式
	if *interactive || (flag.NFlag() == 0 && len(os.Args) == 1) {
		startInteractiveMode(cfg, flag.CommandLine)
		return
	}

	// 否则执行一次性检测
	runOnceMode(cfg, *path, *ruleNames, *format, *outputPath)
}

// startInteractiveMode 启动交互模式
func startInteractiveMode(cfg *config.Config, flagSet *flag.FlagSet) {
	// 初始化命令注册中心
	cmdRegistry := cli.NewRegistry()

	// 注册内置命令
	cmdRegistry.Register(cliBuiltin.NewHelpCommand(cmdRegistry))
	cmdRegistry.Register(cliBuiltin.NewShowCommand(flagSet))
	cmdRegistry.Register(cliBuiltin.NewRunCommand(cfg))
	cmdRegistry.Register(cliBuiltin.NewExitCommand())

	// 启动 REPL
	repl := cli.NewREPL(cmdRegistry)
	repl.Start()
}

// runOnceMode 执行一次性检测
func runOnceMode(cfg *config.Config, path, ruleNames, format, outputPath string) {
	fmt.Println("GopherFlush - AI代码质量检测工具")
	fmt.Println("========================================")

	// 初始化规则注册表
	registry := initializeRules(cfg, ruleNames)

	// 初始化分析器
	a := analyzer.NewAnalyzer(registry)

	// 执行代码检测
	fmt.Printf("正在检测: %s\n", path)
	report, err := a.Analyze(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 生成控制台报告
	consoleReporter := reporter.NewConsoleReporter()
	if err := consoleReporter.Generate(report); err != nil {
		fmt.Fprintf(os.Stderr, "生成报告失败: %v\n", err)
		os.Exit(1)
	}

	// 生成详细报告文件
	if format == "json" || outputPath != "" {
		jsonReporter := reporter.NewJSONReporter(outputPath)
		if err := jsonReporter.Generate(report); err != nil {
			fmt.Fprintf(os.Stderr, "生成详细报告失败: %v\n", err)
		} else {
			fmt.Printf("\n详细报告已保存到: %s\n", outputPath)
		}
	}

	// 如果有违规记录，返回非零退出码
	if report.TotalViolations > 0 {
		os.Exit(1)
	}
}

// loadConfig 加载配置
func loadConfig(configFile string) *config.Config {
	loader := config.NewLoader()

	// 如果指定了配置文件，尝试加载
	if configFile != "" {
		cfg, err := loader.Load(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 加载配置文件失败: %v，使用默认配置\n", err)
			return loader.LoadDefault()
		}
		return cfg
	}

	// 否则使用默认配置
	return loader.LoadDefault()
}

// initializeRules 初始化规则注册表
func initializeRules(cfg *config.Config, ruleNames string) *rules.Registry {
	registry := rules.NewRegistry()

	// 解析要运行的规则列表
	var selectedRules map[string]bool
	if ruleNames != "" {
		selectedRules = make(map[string]bool)
		for _, name := range strings.Split(ruleNames, ",") {
			selectedRules[strings.TrimSpace(name)] = true
		}
	}

	// 注册文件大小规则
	if shouldEnableRule("file-size", selectedRules, cfg.Rules.FileSize.Enabled) {
		registry.Register(builtin.NewFileSizeRule(cfg.Rules.FileSize.MaxLines))
	}

	// 注册函数大小规则
	if shouldEnableRule("function-size", selectedRules, cfg.Rules.FunctionSize.Enabled) {
		registry.Register(builtin.NewFunctionSizeRule(cfg.Rules.FunctionSize.MaxLines))
	}

	// 注册全局变量规则
	if shouldEnableRule("global-vars", selectedRules, cfg.Rules.GlobalVars.Enabled) {
		registry.Register(builtin.NewGlobalVarsRule())
	}

	// 注册重复代码规则
	if shouldEnableRule("duplicates", selectedRules, cfg.Rules.Duplicates.Enabled) {
		registry.Register(builtin.NewDuplicatesRule())
	}

	return registry
}

// shouldEnableRule 判断是否应该启用规则
func shouldEnableRule(ruleName string, selectedRules map[string]bool, configEnabled bool) bool {
	// 如果指定了特定规则，只启用这些规则
	if selectedRules != nil {
		return selectedRules[ruleName]
	}
	// 否则根据配置决定
	return configEnabled
}
