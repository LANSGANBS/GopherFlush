package analyzer

import (
	"fmt"
	"gopherflush/internal/rules"
	"gopherflush/pkg/types"
	"os"
	"path/filepath"
)

// Analyzer 代码分析器
type Analyzer struct {
	registry *rules.Registry
	parser   *Parser
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer(registry *rules.Registry) *Analyzer {
	return &Analyzer{
		registry: registry,
		parser:   NewParser(),
	}
}

// Analyze 分析指定路径的代码
func (a *Analyzer) Analyze(path string) (*types.Report, error) {
	report := &types.Report{
		TotalFiles:      0,
		TotalViolations: 0,
		Violations:      []*types.Violation{},
	}

	// 获取所有启用的规则
	enabledRules := a.registry.GetAll()
	if len(enabledRules) == 0 {
		return report, fmt.Errorf("没有启用的规则")
	}

	// 遍历目录下的所有支持的文件
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检测文件语言
		lang := DetectLanguage(filePath)
		if lang == LanguageUnknown {
			return nil // 跳过不支持的文件
		}

		report.TotalFiles++

		// 根据语言类型选择不同的分析方式
		if lang == LanguageGo {
			// Go 文件：使用 AST 解析
			return a.analyzeGoFile(filePath, enabledRules, report)
		} else {
			// 其他语言：使用文本分析
			return a.analyzeTextFile(filePath, lang, report)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return report, nil
}

// analyzeGoFile 分析Go文件（使用AST）
func (a *Analyzer) analyzeGoFile(filePath string, enabledRules []rules.Rule, report *types.Report) error {
	// 解析文件
	file, err := a.parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("警告: 解析文件 %s 失败: %v\n", filePath, err)
		return nil // 继续处理其他文件
	}

	// 应用所有规则
	fset := a.parser.GetFileSet()
	for _, rule := range enabledRules {
		violations := rule.Check(file, fset, filePath)
		report.Violations = append(report.Violations, violations...)
		report.TotalViolations += len(violations)
	}

	return nil
}

// analyzeTextFile 分析非Go文件（使用文本分析）
func (a *Analyzer) analyzeTextFile(filePath string, lang Language, report *types.Report) error {
	// 创建文本分析器
	textAnalyzer := NewTextAnalyzer(lang)

	// 分析文件
	fileInfo, err := textAnalyzer.AnalyzeFile(filePath)
	if err != nil {
		fmt.Printf("警告: 分析文件 %s 失败: %v\n", filePath, err)
		return nil
	}

	// 检测文件大小
	a.checkFileSize(fileInfo, report)

	// 检测函数大小
	a.checkFunctionSize(fileInfo, report)

	// 检测全局变量
	a.checkGlobalVars(fileInfo, report)

	return nil
}

// checkFileSize 检测文件大小
func (a *Analyzer) checkFileSize(fileInfo *FileInfo, report *types.Report) {
	lineCount := len(fileInfo.Lines)
	maxLines := 800

	if lineCount <= maxLines {
		return
	}

	// 计算严重程度
	var severity types.Severity
	if lineCount > 2000 {
		severity = types.SeverityCritical
	} else if lineCount > 1500 {
		severity = types.SeverityHigh
	} else if lineCount > 1000 {
		severity = types.SeverityMedium
	} else {
		severity = types.SeverityLow
	}

	violation := &types.Violation{
		RuleName:   "file-size",
		Severity:   severity,
		FilePath:   fileInfo.Path,
		Line:       1,
		Column:     1,
		Message:    fmt.Sprintf("文件过大：共 %d 行（阈值：%d 行）[%s]", lineCount, maxLines, fileInfo.Language),
		Suggestion: "考虑将文件拆分为多个较小的文件，提高代码可维护性",
	}

	report.Violations = append(report.Violations, violation)
	report.TotalViolations++
}

// checkFunctionSize 检测函数大小
func (a *Analyzer) checkFunctionSize(fileInfo *FileInfo, report *types.Report) {
	maxLines := 200

	for _, fn := range fileInfo.Functions {
		if fn.LineCount <= maxLines {
			continue
		}

		// 计算严重程度
		var severity types.Severity
		if fn.LineCount > 800 {
			severity = types.SeverityCritical
		} else if fn.LineCount > 500 {
			severity = types.SeverityHigh
		} else if fn.LineCount > 300 {
			severity = types.SeverityMedium
		} else {
			severity = types.SeverityLow
		}

		violation := &types.Violation{
			RuleName:   "function-size",
			Severity:   severity,
			FilePath:   fileInfo.Path,
			Line:       fn.StartLine,
			Column:     1,
			Message:    fmt.Sprintf("函数 '%s' 过大：共 %d 行（阈值：%d 行）[%s]", fn.Name, fn.LineCount, maxLines, fileInfo.Language),
			Suggestion: "考虑将函数拆分为多个较小的函数，提高代码可读性和可维护性",
		}

		report.Violations = append(report.Violations, violation)
		report.TotalViolations++
	}
}

// checkGlobalVars 检测全局变量
func (a *Analyzer) checkGlobalVars(fileInfo *FileInfo, report *types.Report) {
	for _, gv := range fileInfo.Globals {
		// 根据是否可导出设置严重程度
		severity := types.SeverityMedium
		if gv.IsExported {
			severity = types.SeverityHigh
		}

		violation := &types.Violation{
			RuleName:   "global-vars",
			Severity:   severity,
			FilePath:   fileInfo.Path,
			Line:       gv.Line,
			Column:     1,
			Message:    fmt.Sprintf("检测到全局变量 '%s' [%s]", gv.Name, fileInfo.Language),
			Suggestion: "考虑将全局变量改为局部变量、函数参数或使用依赖注入模式",
		}

		report.Violations = append(report.Violations, violation)
		report.TotalViolations++
	}
}
