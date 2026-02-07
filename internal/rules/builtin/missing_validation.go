package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"strings"
)

// MissingValidationRule 缺少验证检测规则
type MissingValidationRule struct {
	enabled bool
}

// NewMissingValidationRule 创建缺少验证检测规则
func NewMissingValidationRule() *MissingValidationRule {
	return &MissingValidationRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *MissingValidationRule) Name() string {
	return "missing-validation"
}

// Description 返回规则描述
func (r *MissingValidationRule) Description() string {
	return "检测缺少异常处理、边界条件检查或输入验证"
}

// Enabled 返回规则是否启用
func (r *MissingValidationRule) Enabled() bool {
	return r.enabled
}

// Check 检查文件中的缺少验证问题
func (r *MissingValidationRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 遍历所有函数
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		// 检查函数体中的危险操作
		r.checkDangerousOperations(funcDecl, fset, filePath, &violations)

		return true
	})

	return violations
}

// checkDangerousOperations 检查函数中的危险操作
func (r *MissingValidationRule) checkDangerousOperations(funcDecl *ast.FuncDecl, fset *token.FileSet, filePath string, violations *[]*types.Violation) {
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// 检查数据库操作
		if callExpr, ok := n.(*ast.CallExpr); ok {
			r.checkDatabaseOperation(callExpr, fset, filePath, violations)
			r.checkDivisionByZero(callExpr, fset, filePath, violations)
		}

		// 检查字符串拼接构建SQL
		if binExpr, ok := n.(*ast.BinaryExpr); ok {
			r.checkSQLStringConcatenation(binExpr, funcDecl, fset, filePath, violations)
		}

		return true
	})
}

// checkDatabaseOperation 检查数据库操作是否缺少输入验证
func (r *MissingValidationRule) checkDatabaseOperation(callExpr *ast.CallExpr, fset *token.FileSet, filePath string, violations *[]*types.Violation) {
	funcName := r.getFunctionName(callExpr)

	// 检查是否是数据库操作
	if !strings.Contains(funcName, "Exec") && !strings.Contains(funcName, "Query") {
		return
	}

	// 检查参数中是否有字符串拼接或变量
	for _, arg := range callExpr.Args {
		if r.isUnsafeArgument(arg) {
			position := fset.Position(callExpr.Pos())
			*violations = append(*violations, &types.Violation{
				RuleName:   r.Name(),
				Severity:   types.SeverityHigh,
				FilePath:   filePath,
				Line:       position.Line,
				Column:     position.Column,
				Message:    "数据库操作可能缺少输入验证，存在SQL注入风险",
				Suggestion: "使用参数化查询（PreparedStatement）或对输入进行严格验证",
			})
			return
		}
	}
}

// checkDivisionByZero 检查除法操作是否缺少零检查
func (r *MissingValidationRule) checkDivisionByZero(callExpr *ast.CallExpr, fset *token.FileSet, filePath string, violations *[]*types.Violation) {
	// 这个方法暂时留空，因为除法操作通常是 BinaryExpr 而不是 CallExpr
	// 实际的除法检查应该在遍历 BinaryExpr 时进行
}

// checkSQLStringConcatenation 检查字符串拼接构建SQL
func (r *MissingValidationRule) checkSQLStringConcatenation(binExpr *ast.BinaryExpr, funcDecl *ast.FuncDecl, fset *token.FileSet, filePath string, violations *[]*types.Violation) {
	// 检查是否是字符串拼接操作
	if binExpr.Op != token.ADD {
		return
	}

	// 检查是否包含SQL关键字
	if r.containsSQLKeywords(binExpr) {
		position := fset.Position(binExpr.Pos())
		*violations = append(*violations, &types.Violation{
			RuleName:   r.Name(),
			Severity:   types.SeverityHigh,
			FilePath:   filePath,
			Line:       position.Line,
			Column:     position.Column,
			Message:    "使用字符串拼接构建SQL语句，存在SQL注入风险",
			Suggestion: "使用参数化查询（PreparedStatement）代替字符串拼接",
		})
	}
}

// getFunctionName 获取函数调用的名称
func (r *MissingValidationRule) getFunctionName(callExpr *ast.CallExpr) string {
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		return fun.Sel.Name
	case *ast.Ident:
		return fun.Name
	}
	return ""
}

// isUnsafeArgument 检查参数是否不安全（包含变量或拼接）
func (r *MissingValidationRule) isUnsafeArgument(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.Ident:
		// 变量
		return true
	case *ast.BinaryExpr:
		// 字符串拼接
		return true
	case *ast.CallExpr:
		// 函数调用结果
		return true
	case *ast.BasicLit:
		// 字面量是安全的
		return false
	}
	return false
}

// containsSQLKeywords 检查表达式是否包含SQL关键字
func (r *MissingValidationRule) containsSQLKeywords(expr ast.Expr) bool {
	sqlKeywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "FROM", "WHERE"}

	// 递归检查表达式中的字符串字面量
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			value := strings.ToUpper(lit.Value)
			for _, keyword := range sqlKeywords {
				if strings.Contains(value, keyword) {
					found = true
					return false
				}
			}
		}
		return true
	})

	return found
}




