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
	"path/filepath"
	"strings"
)

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
	colorDim    = "\033[90m"
)

func main() {
	var (
		interactive = flag.Bool("interactive", false, "启动交互模式")
		path        = flag.String("path", ".", "要检测的代码路径")
		configFile  = flag.String("config", "", "配置文件路径")
		ruleNames   = flag.String("rules", "", "要运行的规则（逗号分隔，留空表示运行所有规则）")
		format      = flag.String("format", "console", "输出格式（console/json）")
		outputPath  = flag.String("output", "gopherflush-report.json", "详细报告输出路径")
	)

	flag.Parse()

	cfg := loadConfig(*configFile)

	if *interactive || (flag.NFlag() == 0 && len(os.Args) == 1) {
		startInteractiveMode(cfg, flag.CommandLine)
		return
	}

	runOnceMode(cfg, *path, *ruleNames, *format, *outputPath)
}

func startInteractiveMode(cfg *config.Config, flagSet *flag.FlagSet) {
	cmdRegistry := cli.NewRegistry()

	cmdRegistry.Register(cliBuiltin.NewHelpCommand(cmdRegistry))
	cmdRegistry.Register(cliBuiltin.NewShowCommand(flagSet))
	cmdRegistry.Register(cliBuiltin.NewRunCommand(cfg))
	cmdRegistry.Register(cliBuiltin.NewExitCommand())

	repl := cli.NewREPL(cmdRegistry)
	repl.Start()
}

func runOnceMode(cfg *config.Config, path, ruleNames, format, outputPath string) {
	fmt.Println()

	printHeader("GopherFlush - 代码检测")

	absPath, _ := filepath.Abs(path)
	fmt.Printf("  %s检测路径:%s %s%s%s\n", colorDim, colorReset, colorCyan, absPath, colorReset)
	fmt.Printf("  %s状态:%s     %s正在分析...%s\n", colorDim, colorReset, colorYellow, colorReset)
	fmt.Println()

	registry := initializeRules(cfg, ruleNames)

	a := analyzer.NewAnalyzer(registry)

	report, err := a.Analyze(path)
	if err != nil {
		fmt.Printf("  %s✗ 检测失败:%s %s\n", colorRed, colorReset, err.Error())
		os.Exit(1)
	}

	consoleReporter := reporter.NewConsoleReporter()
	if err := consoleReporter.Generate(report); err != nil {
		fmt.Fprintf(os.Stderr, "生成报告失败: %v\n", err)
		os.Exit(1)
	}

	if format == "json" || outputPath != "" {
		jsonReporter := reporter.NewJSONReporter(outputPath)
		if err := jsonReporter.Generate(report); err != nil {
			fmt.Fprintf(os.Stderr, "生成详细报告失败: %v\n", err)
		} else {
			absOutputPath, err := filepath.Abs(outputPath)
			if err != nil {
				absOutputPath = outputPath
			}
			fmt.Printf("  %s✓ 详细报告已保存:%s %s\n", colorGreen, colorReset, absOutputPath)
		}
	}

	fmt.Println()

	if report.TotalViolations > 0 {
		os.Exit(1)
	}
}

func printHeader(title string) {
	width := 60
	fmt.Printf("%s%s%s\n", colorCyan, strings.Repeat("─", width), colorReset)
	padding := (width - len(title)) / 2
	fmt.Printf("%s%s%s%s%s%s\n", colorCyan, strings.Repeat(" ", padding), colorBold+colorWhite, title, colorReset, colorReset)
	fmt.Printf("%s%s%s\n", colorCyan, strings.Repeat("─", width), colorReset)
}

func loadConfig(configFile string) *config.Config {
	loader := config.NewLoader()

	if configFile != "" {
		cfg, err := loader.Load(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 加载配置文件失败: %v，使用默认配置\n", err)
			return loader.LoadDefault()
		}
		return cfg
	}

	return loader.LoadDefault()
}

func initializeRules(cfg *config.Config, ruleNames string) *rules.Registry {
	registry := rules.NewRegistry()

	var selectedRules map[string]bool
	if ruleNames != "" {
		selectedRules = make(map[string]bool)
		for _, name := range strings.Split(ruleNames, ",") {
			selectedRules[strings.TrimSpace(name)] = true
		}
	}

	if shouldEnableRule("file-size", selectedRules, cfg.Rules.FileSize.Enabled) {
		registry.Register(builtin.NewFileSizeRule(cfg.Rules.FileSize.MaxLines))
	}

	if shouldEnableRule("function-size", selectedRules, cfg.Rules.FunctionSize.Enabled) {
		registry.Register(builtin.NewFunctionSizeRule(cfg.Rules.FunctionSize.MaxLines))
	}

	if shouldEnableRule("global-vars", selectedRules, cfg.Rules.GlobalVars.Enabled) {
		registry.Register(builtin.NewGlobalVarsRule())
	}

	if shouldEnableRule("duplicates", selectedRules, cfg.Rules.Duplicates.Enabled) {
		registry.Register(builtin.NewDuplicatesRule())
	}

	if shouldEnableRule("commented-code", selectedRules, cfg.Rules.CommentedCode.Enabled) {
		registry.Register(builtin.NewCommentedCodeRule())
	}

	if shouldEnableRule("inconsistent-comment", selectedRules, cfg.Rules.InconsistentComment.Enabled) {
		registry.Register(builtin.NewInconsistentCommentRule())
	}

	if shouldEnableRule("resource-leak", selectedRules, cfg.Rules.ResourceLeak.Enabled) {
		registry.Register(builtin.NewResourceLeakRule())
	}

	if shouldEnableRule("loose-typing", selectedRules, cfg.Rules.LooseTyping.Enabled) {
		registry.Register(builtin.NewLooseTypingRule())
	}

	if shouldEnableRule("inaccurate-constant", selectedRules, cfg.Rules.InaccurateConstant.Enabled) {
		registry.Register(builtin.NewInaccurateConstantRule())
	}

	if shouldEnableRule("missing-validation", selectedRules, cfg.Rules.MissingValidation.Enabled) {
		registry.Register(builtin.NewMissingValidationRule())
	}

	if shouldEnableRule("hardcoded-secrets", selectedRules, cfg.Rules.HardcodedSecrets.Enabled) {
		registry.Register(builtin.NewHardcodedSecretsRule())
	}

	return registry
}

func shouldEnableRule(ruleName string, selectedRules map[string]bool, configEnabled bool) bool {
	if selectedRules != nil {
		return selectedRules[ruleName]
	}
	return configEnabled
}
