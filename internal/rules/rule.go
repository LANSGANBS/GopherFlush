package rules

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
)

// Rule 规则接口
type Rule interface {
	// Name 返回规则名称
	Name() string

	// Description 返回规则描述
	Description() string

	// Check 检查代码并返回违规记录
	Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation

	// Enabled 规则是否启用
	Enabled() bool
}
