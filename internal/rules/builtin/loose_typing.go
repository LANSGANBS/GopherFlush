package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
)

// LooseTypingRule 宽泛类型检测规则
type LooseTypingRule struct {
	enabled bool
}

// NewLooseTypingRule 创建宽泛类型检测规则
func NewLooseTypingRule() *LooseTypingRule {
	return &LooseTypingRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *LooseTypingRule) Name() string {
	return "loose-typing"
}

// Description 返回规则描述
func (r *LooseTypingRule) Description() string {
	return "检测函数参数和返回值类型约束过于宽泛（如interface{}）"
}

// Enabled 返回规则是否启用
func (r *LooseTypingRule) Enabled() bool {
	return r.enabled
}

// Check 检查文件中的宽泛类型使用
func (r *LooseTypingRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 遍历所有函数
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// 检查函数参数
		if funcDecl.Type.Params != nil {
			for _, field := range funcDecl.Type.Params.List {
				if r.isLooseType(field.Type) {
					position := fset.Position(field.Pos())
					violation := &types.Violation{
						RuleName:   r.Name(),
						Severity:   types.SeverityMedium,
						FilePath:   filePath,
						Line:       position.Line,
						Column:     position.Column,
						Message:    "函数 '" + funcDecl.Name.Name + "' 的参数使用了过于宽泛的类型 interface{}",
						Suggestion: "使用具体的类型或定义明确的接口，以提供更好的类型安全性",
					}
					violations = append(violations, violation)
				}
			}
		}

		// 检查函数返回值
		if funcDecl.Type.Results != nil {
			for _, field := range funcDecl.Type.Results.List {
				if r.isLooseType(field.Type) {
					position := fset.Position(field.Pos())
					violation := &types.Violation{
						RuleName:   r.Name(),
						Severity:   types.SeverityMedium,
						FilePath:   filePath,
						Line:       position.Line,
						Column:     position.Column,
						Message:    "函数 '" + funcDecl.Name.Name + "' 的返回值使用了过于宽泛的类型 interface{}",
						Suggestion: "使用具体的类型或定义明确的接口，以提供更好的类型安全性",
					}
					violations = append(violations, violation)
				}
			}
		}

		return true
	})

	return violations
}

// isLooseType 检查类型是否为宽泛类型
func (r *LooseTypingRule) isLooseType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.InterfaceType:
		// 空接口 interface{}
		if t.Methods == nil || len(t.Methods.List) == 0 {
			return true
		}
	case *ast.Ident:
		// any 类型（Go 1.18+）
		if t.Name == "any" {
			return true
		}
	}
	return false
}
