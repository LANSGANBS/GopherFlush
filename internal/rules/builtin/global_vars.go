package builtin

import (
	"fmt"
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
	violations := []*types.Violation{}

	// 统计全局变量数量，用于确定严重程度
	globalVarCount := 0

	// 遍历文件中的所有顶层声明
	for _, decl := range file.Decls {
		// 只处理通用声明（GenDecl）
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// 只处理变量声明（VAR）
		if genDecl.Tok != token.VAR {
			continue
		}

		// 遍历声明中的所有变量
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// 为每个变量名创建违规记录
			for _, name := range valueSpec.Names {
				globalVarCount++

				// 获取位置信息
				pos := fset.Position(name.Pos())

				// 判断是否是可导出的全局变量（大写开头）
				isExported := name.IsExported()
				varName := name.Name

				violation := &types.Violation{
					RuleName:   r.Name(),
					Severity:   r.calculateSeverity(isExported),
					FilePath:   filePath,
					Line:       pos.Line,
					Column:     pos.Column,
					Message:    fmt.Sprintf("检测到全局变量 '%s'", varName),
					Suggestion: "考虑将全局变量改为局部变量、函数参数或使用依赖注入模式",
				}

				violations = append(violations, violation)
			}
		}
	}

	return violations
}

func (r *GlobalVarsRule) Enabled() bool {
	return r.enabled
}

// calculateSeverity 根据全局变量特征计算严重程度
func (r *GlobalVarsRule) calculateSeverity(isExported bool) types.Severity {
	// 可导出的全局变量更严重，因为它们可以被包外访问
	if isExported {
		return types.SeverityHigh // 严重
	}
	// 不可导出的全局变量相对较轻
	return types.SeverityMedium // 中等
}
