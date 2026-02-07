package analyzer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopherflush/internal/rules"
	"gopherflush/pkg/types"
	"os"
	"path/filepath"
	"strings"
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

	// 只运行启用的规则对应的检测
	if a.isRuleEnabled("file-size") {
		a.checkFileSize(fileInfo, report)
	}

	if a.isRuleEnabled("function-size") {
		a.checkFunctionSize(fileInfo, report)
	}

	if a.isRuleEnabled("global-vars") {
		a.checkGlobalVars(fileInfo, report)
	}

	if a.isRuleEnabled("duplicates") {
		a.checkDuplicates(fileInfo, report)
	}

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

// checkDuplicates 检测重复函数
func (a *Analyzer) checkDuplicates(fileInfo *FileInfo, report *types.Report) {
	// 收集函数内容和哈希值
	type funcWithHash struct {
		info *FunctionInfo
		hash string
	}

	funcsWithHash := make([]funcWithHash, 0, len(fileInfo.Functions))

	for _, fn := range fileInfo.Functions {
		// 提取函数内容
		content := a.extractFunctionContent(fileInfo.Lines, fn.StartLine, fn.EndLine)

		// 标准化并计算哈希
		normalized := a.normalizeTextContent(content)
		hash := a.calculateTextHash(normalized)

		funcsWithHash = append(funcsWithHash, funcWithHash{
			info: fn,
			hash: hash,
		})
	}

	// 查找重复
	hashMap := make(map[string][]*FunctionInfo)
	for _, fh := range funcsWithHash {
		// 跳过空函数
		if fh.hash == a.calculateTextHash("") {
			continue
		}
		hashMap[fh.hash] = append(hashMap[fh.hash], fh.info)
	}

	// 生成违规记录
	for _, funcs := range hashMap {
		if len(funcs) < 2 {
			continue
		}

		// 收集函数名
		funcNames := make([]string, 0, len(funcs))
		for _, fn := range funcs {
			funcNames = append(funcNames, fn.Name)
		}

		// 确定严重程度
		severity := types.SeverityMedium
		if len(funcs) > 3 {
			severity = types.SeverityHigh
		}

		// 只为第一个函数创建违规记录
		fn := funcs[0]
		message := fmt.Sprintf("检测到 %d 个重复函数: %s [%s]", len(funcs),
			strings.Join(funcNames, ", "), fileInfo.Language)

		violation := &types.Violation{
			RuleName:   "duplicates",
			Severity:   severity,
			FilePath:   fileInfo.Path,
			Line:       fn.StartLine,
			Column:     1,
			Message:    message,
			Suggestion: "考虑提取公共逻辑到一个函数中，避免代码重复",
		}

		report.Violations = append(report.Violations, violation)
		report.TotalViolations++
	}
}

// extractFunctionContent 提取函数内容
func (a *Analyzer) extractFunctionContent(lines []string, startLine, endLine int) string {
	if startLine < 1 || endLine > len(lines) || startLine > endLine {
		return ""
	}

	// 注意：行号是从1开始的，但数组索引是从0开始的
	content := ""
	for i := startLine - 1; i < endLine && i < len(lines); i++ {
		content += lines[i] + "\n"
	}
	return content
}

// normalizeTextContent 标准化文本内容（去除空白、换行等）
func (a *Analyzer) normalizeTextContent(content string) string {
	// 去除所有空白字符
	normalized := ""
	for _, ch := range content {
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			normalized += string(ch)
		}
	}
	return normalized
}

// calculateTextHash 计算文本内容的哈希值
func (a *Analyzer) calculateTextHash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// joinStrings 连接字符串数组
func (a *Analyzer) joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// isRuleEnabled 检查规则是否启用
func (a *Analyzer) isRuleEnabled(ruleName string) bool {
	rules := a.registry.GetAll()
	for _, rule := range rules {
		if rule.Name() == ruleName {
			return true
		}
	}
	return false
}
