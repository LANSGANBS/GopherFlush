package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"strings"
)

// HardcodedSecretsRule 硬编码配置检测规则
type HardcodedSecretsRule struct {
	enabled bool
}

// NewHardcodedSecretsRule 创建硬编码配置检测规则
func NewHardcodedSecretsRule() *HardcodedSecretsRule {
	return &HardcodedSecretsRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *HardcodedSecretsRule) Name() string {
	return "hardcoded-secrets"
}

// Description 返回规则描述
func (r *HardcodedSecretsRule) Description() string {
	return "检测硬编码的配置、URL、API token等敏感信息"
}

// Enabled 返回规则是否启用
func (r *HardcodedSecretsRule) Enabled() bool {
	return r.enabled
}

// Check 检查文件中的硬编码配置
func (r *HardcodedSecretsRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 遍历所有声明
	ast.Inspect(file, func(n ast.Node) bool {
		// 检查变量和常量声明
		if genDecl, ok := n.(*ast.GenDecl); ok {
			r.checkDeclarations(genDecl, fset, filePath, &violations)
		}

		// 检查赋值语句
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			r.checkAssignments(assignStmt, fset, filePath, &violations)
		}

		return true
	})

	return violations
}

// checkDeclarations 检查变量和常量声明
func (r *HardcodedSecretsRule) checkDeclarations(genDecl *ast.GenDecl, fset *token.FileSet, filePath string, violations *[]*types.Violation) {
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for i, value := range valueSpec.Values {
			if i >= len(valueSpec.Names) {
				continue
			}

			varName := valueSpec.Names[i].Name
			if violation := r.checkValue(value, varName, fset, filePath); violation != nil {
				*violations = append(*violations, violation)
			}
		}
	}
}

// checkAssignments 检查赋值语句
func (r *HardcodedSecretsRule) checkAssignments(assignStmt *ast.AssignStmt, fset *token.FileSet, filePath string, violations *[]*types.Violation) {
	for i, rhs := range assignStmt.Rhs {
		if i >= len(assignStmt.Lhs) {
			continue
		}

		varName := ""
		if ident, ok := assignStmt.Lhs[i].(*ast.Ident); ok {
			varName = ident.Name
		}

		if violation := r.checkValue(rhs, varName, fset, filePath); violation != nil {
			*violations = append(*violations, violation)
		}
	}
}

// checkValue 检查值是否包含硬编码的敏感信息
func (r *HardcodedSecretsRule) checkValue(expr ast.Expr, varName string, fset *token.FileSet, filePath string) *types.Violation {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}

	value := lit.Value
	position := fset.Position(lit.Pos())

	// 检查URL
	if r.isURL(value) {
		return &types.Violation{
			RuleName:   r.Name(),
			Severity:   types.SeverityMedium,
			FilePath:   filePath,
			Line:       position.Line,
			Column:     position.Column,
			Message:    "检测到硬编码的URL: " + varName,
			Suggestion: "将URL配置移到配置文件或环境变量中",
		}
	}

	// 检查API token/key
	if r.isAPIToken(varName, value) {
		return &types.Violation{
			RuleName:   r.Name(),
			Severity:   types.SeverityHigh,
			FilePath:   filePath,
			Line:       position.Line,
			Column:     position.Column,
			Message:    "检测到硬编码的API token或密钥: " + varName,
			Suggestion: "将敏感信息移到环境变量或密钥管理系统中",
		}
	}

	// 检查数据库连接字符串
	if r.isDatabaseConnection(value) {
		return &types.Violation{
			RuleName:   r.Name(),
			Severity:   types.SeverityHigh,
			FilePath:   filePath,
			Line:       position.Line,
			Column:     position.Column,
			Message:    "检测到硬编码的数据库连接字符串",
			Suggestion: "将数据库连接信息移到配置文件或环境变量中",
		}
	}

	return nil
}

// isURL 检查字符串是否是URL
func (r *HardcodedSecretsRule) isURL(value string) bool {
	value = strings.ToLower(value)
	return strings.Contains(value, "http://") || strings.Contains(value, "https://")
}

// isAPIToken 检查变量名和值是否表示API token
func (r *HardcodedSecretsRule) isAPIToken(varName, value string) bool {
	varNameLower := strings.ToLower(varName)

	// 检查变量名是否包含敏感关键字
	sensitiveKeywords := []string{"token", "key", "secret", "password", "passwd", "pwd", "apikey", "api_key"}
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(varNameLower, keyword) {
			// 检查值是否看起来像真实的token（长度大于10且不是占位符）
			cleanValue := strings.Trim(value, "\"'`")
			if len(cleanValue) > 10 && !strings.Contains(strings.ToLower(cleanValue), "your") && !strings.Contains(strings.ToLower(cleanValue), "example") {
				return true
			}
		}
	}

	return false
}

// isDatabaseConnection 检查字符串是否是数据库连接字符串
func (r *HardcodedSecretsRule) isDatabaseConnection(value string) bool {
	valueLower := strings.ToLower(value)

	// 常见的数据库连接字符串模式
	patterns := []string{
		"mysql://",
		"postgresql://",
		"postgres://",
		"mongodb://",
		"redis://",
		"user:password@",
		"jdbc:",
	}

	for _, pattern := range patterns {
		if strings.Contains(valueLower, pattern) {
			return true
		}
	}

	return false
}



