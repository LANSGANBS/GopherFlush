package analyzer

import (
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

	// 遍历目录下的所有 Go 文件
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .go 文件
		if !strings.HasSuffix(filePath, ".go") {
			return nil
		}

		// 解析文件
		file, err := a.parser.ParseFile(filePath)
		if err != nil {
			fmt.Printf("警告: 解析文件 %s 失败: %v\n", filePath, err)
			return nil // 继续处理其他文件
		}

		report.TotalFiles++

		// 应用所有规则
		fset := a.parser.GetFileSet()
		for _, rule := range enabledRules {
			violations := rule.Check(file, fset, filePath)
			report.Violations = append(report.Violations, violations...)
			report.TotalViolations += len(violations)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return report, nil
}
