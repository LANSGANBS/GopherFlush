package builtin

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"os"
)

// FunctionSizeRule 函数大小检测规则
type FunctionSizeRule struct {
	maxLines int
	enabled  bool
}

// NewFunctionSizeRule 创建函数大小规则
func NewFunctionSizeRule(maxLines int) *FunctionSizeRule {
	return &FunctionSizeRule{
		maxLines: maxLines,
		enabled:  true,
	}
}

func (r *FunctionSizeRule) Name() string {
	return "function-size"
}

func (r *FunctionSizeRule) Description() string {
	return "检测函数是否超过指定行数"
}

func (r *FunctionSizeRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 读取文件内容以计算行数
	content, err := os.ReadFile(filePath)
	if err != nil {
		return violations
	}

	// 遍历所有函数声明
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// 计算函数行数
		lineCount := r.calculateFunctionLines(funcDecl, content)

		// 如果函数行数超过阈值，创建违规记录
		if lineCount > r.maxLines {
			severity := r.calculateSeverity(lineCount)
			funcName := funcDecl.Name.Name

			// 使用 FileSet 获取正确的行号
			pos := fset.Position(funcDecl.Pos())

			violation := &types.Violation{
				RuleName:   r.Name(),
				Severity:   severity,
				FilePath:   filePath,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    fmt.Sprintf("函数 '%s' 过大：共 %d 行（阈值：%d 行）", funcName, lineCount, r.maxLines),
				Suggestion: "考虑将函数拆分为多个较小的函数，提高代码可读性和可维护性",
			}
			violations = append(violations, violation)
		}

		return true
	})

	return violations
}

func (r *FunctionSizeRule) Enabled() bool {
	return r.enabled
}

// calculateFunctionLines 计算函数的行数
func (r *FunctionSizeRule) calculateFunctionLines(funcDecl *ast.FuncDecl, content []byte) int {
	// 获取函数的起始和结束位置
	start := int(funcDecl.Pos()) - 1
	end := int(funcDecl.End()) - 1

	// 确保位置在有效范围内
	if start < 0 || end > len(content) || start >= end {
		return 0
	}

	// 计算函数体内的行数
	funcContent := content[start:end]
	scanner := bufio.NewScanner(bufio.NewReader(bufio.NewReader(nil)))
	scanner.Split(bufio.ScanLines)

	lineCount := 1 // 至少有一行
	for _, b := range funcContent {
		if b == '\n' {
			lineCount++
		}
	}

	return lineCount
}

// calculateSeverity 根据函数行数计算严重程度
func (r *FunctionSizeRule) calculateSeverity(lineCount int) types.Severity {
	if lineCount > 800 {
		return types.SeverityCritical // 极其严重
	} else if lineCount > 500 {
		return types.SeverityHigh // 严重
	} else if lineCount > 300 {
		return types.SeverityMedium // 中等
	}
	return types.SeverityLow // 低
}
