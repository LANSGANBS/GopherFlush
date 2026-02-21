package analyzer

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopherflush/internal/rules"
	"gopherflush/pkg/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Analyzer struct {
	registry *rules.Registry
	parser   *Parser
}

func NewAnalyzer(registry *rules.Registry) *Analyzer {
	return &Analyzer{
		registry: registry,
		parser:   NewParser(),
	}
}

func (a *Analyzer) Analyze(path string) (*types.Report, error) {
	report := &types.Report{
		TotalFiles:      0,
		TotalViolations: 0,
		Violations:      []*types.Violation{},
	}

	enabledRules := a.registry.GetAll()
	if len(enabledRules) == 0 {
		return report, fmt.Errorf("没有启用的规则")
	}

	var filePaths []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		lang := DetectLanguage(filePath)
		if lang == LanguageUnknown {
			return nil
		}

		filePaths = append(filePaths, filePath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	report.TotalFiles = len(filePaths)

	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		workerSem  = make(chan struct{}, 8)
		violations = make([]*types.Violation, 0, len(filePaths)*10)
	)

	ruleSet := a.buildRuleSet(enabledRules)

	for _, filePath := range filePaths {
		wg.Add(1)
		go func(fp string) {
			defer wg.Done()
			workerSem <- struct{}{}
			defer func() { <-workerSem }()

			lang := DetectLanguage(fp)
			var fileViolations []*types.Violation

			if lang == LanguageGo {
				fileViolations = a.analyzeGoFileFast(fp, enabledRules)
			} else {
				fileViolations = a.analyzeTextFileFast(fp, lang, ruleSet)
			}

			if len(fileViolations) > 0 {
				mu.Lock()
				violations = append(violations, fileViolations...)
				mu.Unlock()
			}
		}(filePath)
	}

	wg.Wait()

	report.Violations = violations
	report.TotalViolations = len(violations)

	return report, nil
}

type ruleSet struct {
	fileSize     bool
	functionSize bool
	globalVars   bool
	duplicates   bool
}

func (a *Analyzer) buildRuleSet(enabledRules []rules.Rule) *ruleSet {
	rs := &ruleSet{}
	for _, rule := range enabledRules {
		switch rule.Name() {
		case "file-size":
			rs.fileSize = true
		case "function-size":
			rs.functionSize = true
		case "global-vars":
			rs.globalVars = true
		case "duplicates":
			rs.duplicates = true
		}
	}
	return rs
}

func (a *Analyzer) analyzeGoFileFast(filePath string, enabledRules []rules.Rule) []*types.Violation {
	file, err := a.parser.ParseFile(filePath)
	if err != nil {
		return nil
	}

	violations := make([]*types.Violation, 0, 10)
	fset := a.parser.GetFileSet()

	for _, rule := range enabledRules {
		vs := rule.Check(file, fset, filePath)
		violations = append(violations, vs...)
	}

	return violations
}

func (a *Analyzer) analyzeTextFileFast(filePath string, lang Language, rs *ruleSet) []*types.Violation {
	textAnalyzer := NewTextAnalyzer(lang)

	fileInfo, err := textAnalyzer.AnalyzeFile(filePath)
	if err != nil {
		return nil
	}

	violations := make([]*types.Violation, 0, 10)

	if rs.fileSize {
		violations = append(violations, a.checkFileSizeFast(fileInfo)...)
	}

	if rs.functionSize {
		violations = append(violations, a.checkFunctionSizeFast(fileInfo)...)
	}

	if rs.globalVars {
		violations = append(violations, a.checkGlobalVarsFast(fileInfo)...)
	}

	if rs.duplicates {
		violations = append(violations, a.checkDuplicatesFast(fileInfo)...)
	}

	return violations
}

func (a *Analyzer) checkFileSizeFast(fileInfo *FileInfo) []*types.Violation {
	lineCount := len(fileInfo.Lines)
	maxLines := 800

	if lineCount <= maxLines {
		return nil
	}

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

	return []*types.Violation{{
		RuleName:   "file-size",
		Severity:   severity,
		FilePath:   fileInfo.Path,
		Line:       1,
		Column:     1,
		Message:    fmt.Sprintf("文件过大：共 %d 行（阈值：%d 行）[%s]", lineCount, maxLines, fileInfo.Language),
		Suggestion: "考虑将文件拆分为多个较小的文件，提高代码可维护性",
	}}
}

func (a *Analyzer) checkFunctionSizeFast(fileInfo *FileInfo) []*types.Violation {
	maxLines := 200
	violations := make([]*types.Violation, 0, len(fileInfo.Functions)/4)

	for _, fn := range fileInfo.Functions {
		if fn.LineCount <= maxLines {
			continue
		}

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

		violations = append(violations, &types.Violation{
			RuleName:   "function-size",
			Severity:   severity,
			FilePath:   fileInfo.Path,
			Line:       fn.StartLine,
			Column:     1,
			Message:    fmt.Sprintf("函数 '%s' 过大：共 %d 行（阈值：%d 行）[%s]", fn.Name, fn.LineCount, maxLines, fileInfo.Language),
			Suggestion: "考虑将函数拆分为多个较小的函数，提高代码可读性和可维护性",
		})
	}

	return violations
}

func (a *Analyzer) checkGlobalVarsFast(fileInfo *FileInfo) []*types.Violation {
	violations := make([]*types.Violation, 0, len(fileInfo.Globals))

	for _, gv := range fileInfo.Globals {
		severity := types.SeverityMedium
		if gv.IsExported {
			severity = types.SeverityHigh
		}

		violations = append(violations, &types.Violation{
			RuleName:   "global-vars",
			Severity:   severity,
			FilePath:   fileInfo.Path,
			Line:       gv.Line,
			Column:     1,
			Message:    fmt.Sprintf("检测到全局变量 '%s' [%s]", gv.Name, fileInfo.Language),
			Suggestion: "考虑将全局变量改为局部变量、函数参数或使用依赖注入模式",
		})
	}

	return violations
}

func (a *Analyzer) checkDuplicatesFast(fileInfo *FileInfo) []*types.Violation {
	if len(fileInfo.Functions) < 2 {
		return nil
	}

	hashMap := make(map[string][]*FunctionInfo, len(fileInfo.Functions))
	emptyHash := calculateHashFast("")

	for _, fn := range fileInfo.Functions {
		content := extractFunctionContentFast(fileInfo.Lines, fn.StartLine, fn.EndLine)
		normalized := normalizeTextFast(content)
		hash := calculateHashFast(normalized)

		if hash == emptyHash {
			continue
		}

		hashMap[hash] = append(hashMap[hash], fn)
	}

	violations := make([]*types.Violation, 0, len(hashMap)/4)

	for _, funcs := range hashMap {
		if len(funcs) < 2 {
			continue
		}

		funcNames := make([]string, 0, len(funcs))
		for _, fn := range funcs {
			funcNames = append(funcNames, fn.Name)
		}

		severity := types.SeverityMedium
		if len(funcs) > 3 {
			severity = types.SeverityHigh
		}

		fn := funcs[0]
		violations = append(violations, &types.Violation{
			RuleName:   "duplicates",
			Severity:   severity,
			FilePath:   fileInfo.Path,
			Line:       fn.StartLine,
			Column:     1,
			Message:    fmt.Sprintf("检测到 %d 个重复函数: %s [%s]", len(funcs), strings.Join(funcNames, ", "), fileInfo.Language),
			Suggestion: "考虑提取公共逻辑到一个函数中，避免代码重复",
		})
	}

	return violations
}

func extractFunctionContentFast(lines []string, startLine, endLine int) string {
	if startLine < 1 || endLine > len(lines) || startLine > endLine {
		return ""
	}

	var buf bytes.Buffer
	buf.Grow((endLine - startLine + 1) * 80)

	for i := startLine - 1; i < endLine && i < len(lines); i++ {
		buf.WriteString(lines[i])
		buf.WriteByte('\n')
	}

	return buf.String()
}

func normalizeTextFast(content string) string {
	var buf bytes.Buffer
	buf.Grow(len(content))

	for _, ch := range content {
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func calculateHashFast(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

func (a *Analyzer) AnalyzeFile(filePath string) (*types.Report, error) {
	report := &types.Report{
		TotalFiles:      1,
		TotalViolations: 0,
		Violations:      []*types.Violation{},
	}

	enabledRules := a.registry.GetAll()
	if len(enabledRules) == 0 {
		return report, fmt.Errorf("没有启用的规则")
	}

	lang := DetectLanguage(filePath)
	if lang == LanguageUnknown {
		return report, nil
	}

	var violations []*types.Violation

	if lang == LanguageGo {
		violations = a.analyzeGoFileFast(filePath, enabledRules)
	} else {
		ruleSet := a.buildRuleSet(enabledRules)
		violations = a.analyzeTextFileFast(filePath, lang, ruleSet)
	}

	report.Violations = violations
	report.TotalViolations = len(violations)

	return report, nil
}

func (a *Analyzer) analyzeGoFile(filePath string, enabledRules []rules.Rule, report *types.Report) error {
	file, err := a.parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("警告: 解析文件 %s 失败: %v\n", filePath, err)
		return nil
	}

	fset := a.parser.GetFileSet()
	for _, rule := range enabledRules {
		violations := rule.Check(file, fset, filePath)
		report.Violations = append(report.Violations, violations...)
		report.TotalViolations += len(violations)
	}

	return nil
}

func (a *Analyzer) analyzeTextFile(filePath string, lang Language, report *types.Report) error {
	textAnalyzer := NewTextAnalyzer(lang)

	fileInfo, err := textAnalyzer.AnalyzeFile(filePath)
	if err != nil {
		fmt.Printf("警告: 分析文件 %s 失败: %v\n", filePath, err)
		return nil
	}

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

func (a *Analyzer) checkFileSize(fileInfo *FileInfo, report *types.Report) {
	violations := a.checkFileSizeFast(fileInfo)
	report.Violations = append(report.Violations, violations...)
	report.TotalViolations += len(violations)
}

func (a *Analyzer) checkFunctionSize(fileInfo *FileInfo, report *types.Report) {
	violations := a.checkFunctionSizeFast(fileInfo)
	report.Violations = append(report.Violations, violations...)
	report.TotalViolations += len(violations)
}

func (a *Analyzer) checkGlobalVars(fileInfo *FileInfo, report *types.Report) {
	violations := a.checkGlobalVarsFast(fileInfo)
	report.Violations = append(report.Violations, violations...)
	report.TotalViolations += len(violations)
}

func (a *Analyzer) checkDuplicates(fileInfo *FileInfo, report *types.Report) {
	violations := a.checkDuplicatesFast(fileInfo)
	report.Violations = append(report.Violations, violations...)
	report.TotalViolations += len(violations)
}

func (a *Analyzer) extractFunctionContent(lines []string, startLine, endLine int) string {
	return extractFunctionContentFast(lines, startLine, endLine)
}

func (a *Analyzer) normalizeTextContent(content string) string {
	return normalizeTextFast(content)
}

func (a *Analyzer) calculateTextHash(content string) string {
	return calculateHashFast(content)
}

func (a *Analyzer) joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

func (a *Analyzer) isRuleEnabled(ruleName string) bool {
	rules := a.registry.GetAll()
	for _, rule := range rules {
		if rule.Name() == ruleName {
			return true
		}
	}
	return false
}
