package builtin

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"os"
)

// FileSizeRule 文件大小检测规则
type FileSizeRule struct {
	maxLines int
	enabled  bool
}

// NewFileSizeRule 创建文件大小规则
func NewFileSizeRule(maxLines int) *FileSizeRule {
	return &FileSizeRule{
		maxLines: maxLines,
		enabled:  true,
	}
}

func (r *FileSizeRule) Name() string {
	return "file-size"
}

func (r *FileSizeRule) Description() string {
	return "检测文件是否超过指定行数"
}

func (r *FileSizeRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	// 计算文件行数
	lineCount, err := countFileLines(filePath)
	if err != nil {
		return nil
	}

	// 如果文件行数未超过阈值，返回空
	if lineCount <= r.maxLines {
		return nil
	}

	// 根据超出程度确定严重级别
	severity := r.calculateSeverity(lineCount)

	// 创建违规记录
	violation := &types.Violation{
		RuleName:   r.Name(),
		Severity:   severity,
		FilePath:   filePath,
		Line:       1,
		Column:     1,
		Message:    fmt.Sprintf("文件过大：共 %d 行（阈值：%d 行）", lineCount, r.maxLines),
		Suggestion: "考虑将文件拆分为多个较小的文件，提高代码可维护性",
	}

	return []*types.Violation{violation}
}

// countFileLines 计算文件行数
func countFileLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}

// calculateSeverity 根据文件行数计算严重程度
func (r *FileSizeRule) calculateSeverity(lineCount int) types.Severity {
	if lineCount > 2000 {
		return types.SeverityCritical // 极其严重
	} else if lineCount > 1500 {
		return types.SeverityHigh // 严重
	} else if lineCount > 1000 {
		return types.SeverityMedium // 中等
	}
	return types.SeverityLow // 低
}

func (r *FileSizeRule) Enabled() bool {
	return r.enabled
}
