package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
)

// DuplicatesRule 重复代码检测规则
type DuplicatesRule struct {
	enabled bool
}

// NewDuplicatesRule 创建重复代码规则
func NewDuplicatesRule() *DuplicatesRule {
	return &DuplicatesRule{
		enabled: true,
	}
}

func (r *DuplicatesRule) Name() string {
	return "duplicates"
}

func (r *DuplicatesRule) Description() string {
	return "检测重复的变量或函数声明"
}

func (r *DuplicatesRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	// TODO: 实现重复代码检测逻辑
	return nil
}

func (r *DuplicatesRule) Enabled() bool {
	return r.enabled
}
