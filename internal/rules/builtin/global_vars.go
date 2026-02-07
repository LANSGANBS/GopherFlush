package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
)

// GlobalVarsRule 全局变量检测规则
type GlobalVarsRule struct {
	enabled bool
}

// NewGlobalVarsRule 创建全局变量规则
func NewGlobalVarsRule() *GlobalVarsRule {
	return &GlobalVarsRule{
		enabled: true,
	}
}

func (r *GlobalVarsRule) Name() string {
	return "global-vars"
}

func (r *GlobalVarsRule) Description() string {
	return "检测全局变量滥用"
}

func (r *GlobalVarsRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	// TODO: 实现全局变量检测逻辑
	return nil
}

func (r *GlobalVarsRule) Enabled() bool {
	return r.enabled
}
