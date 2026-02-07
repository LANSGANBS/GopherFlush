package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"math"
	"strconv"
	"strings"
)

// InaccurateConstantRule 不准确常量检测规则
type InaccurateConstantRule struct {
	enabled bool
}

// NewInaccurateConstantRule 创建不准确常量检测规则
func NewInaccurateConstantRule() *InaccurateConstantRule {
	return &InaccurateConstantRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *InaccurateConstantRule) Name() string {
	return "inaccurate-constant"
}

// Description 返回规则描述
func (r *InaccurateConstantRule) Description() string {
	return "检测AI生成的不准确常量定义"
}

// Enabled 返回规则是否启用
func (r *InaccurateConstantRule) Enabled() bool {
	return r.enabled
}

// 已知的准确常量值
var knownConstants = map[string]float64{
	"PI":  3.14159265358979323846,
	"E":   2.71828182845904523536,
	"PHI": 1.61803398874989484820, // 黄金比例
}

// Check 检查文件中的不准确常量
func (r *InaccurateConstantRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 遍历所有声明
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		// 遍历常量声明
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// 检查每个常量
			for i, name := range valueSpec.Names {
				if i >= len(valueSpec.Values) {
					continue
				}

				// 检查常量名称和值
				violation := r.checkConstant(name.Name, valueSpec.Values[i], fset, filePath)
				if violation != nil {
					violations = append(violations, violation)
				}
			}
		}
	}

	return violations
}

// checkConstant 检查单个常量是否准确
func (r *InaccurateConstantRule) checkConstant(name string, value ast.Expr, fset *token.FileSet, filePath string) *types.Violation {
	// 检查常量名称是否匹配已知常量
	expectedValue, exists := knownConstants[strings.ToUpper(name)]
	if !exists {
		return nil
	}

	// 提取常量值
	actualValue := r.extractValue(value)
	if actualValue == 0 {
		return nil
	}

	// 比较值的准确性（允许小的误差）
	tolerance := 1e-12
	diff := math.Abs(expectedValue - actualValue)

	if diff > tolerance {
		position := fset.Position(value.Pos())
		return &types.Violation{
			RuleName:   r.Name(),
			Severity:   types.SeverityLow,
			FilePath:   filePath,
			Line:       position.Line,
			Column:     position.Column,
			Message:    "常量 '" + name + "' 的值不准确，当前值: " + strconv.FormatFloat(actualValue, 'f', -1, 64) + "，正确值应为: " + strconv.FormatFloat(expectedValue, 'f', -1, 64),
			Suggestion: "使用更准确的常量值，或使用标准库中的常量（如 math.Pi）",
		}
	}

	return nil
}

// extractValue 从 AST 表达式中提取数值
func (r *InaccurateConstantRule) extractValue(expr ast.Expr) float64 {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.FLOAT || v.Kind == token.INT {
			val, err := strconv.ParseFloat(v.Value, 64)
			if err == nil {
				return val
			}
		}
	}
	return 0
}

