package builtin

import (
	"fmt"
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
)

type CyclomaticComplexityRule struct {
	enabled       bool
	maxComplexity int
}

func NewCyclomaticComplexityRule() *CyclomaticComplexityRule {
	return &CyclomaticComplexityRule{
		enabled:       true,
		maxComplexity: 15,
	}
}

func (r *CyclomaticComplexityRule) Name() string {
	return "cyclomatic-complexity"
}

func (r *CyclomaticComplexityRule) Description() string {
	return "检测圈复杂度过高的函数"
}

func (r *CyclomaticComplexityRule) Enabled() bool {
	return r.enabled
}

func (r *CyclomaticComplexityRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		complexity := r.calculateComplexity(funcDecl.Body)
		if complexity > r.maxComplexity {
			pos := fset.Position(funcDecl.Pos())
			severity := r.determineSeverity(complexity)

			violations = append(violations, &types.Violation{
				RuleName:   r.Name(),
				Severity:   severity,
				FilePath:   filePath,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    fmt.Sprintf("函数 '%s' 圈复杂度过高: %d (阈值: %d)", funcDecl.Name.Name, complexity, r.maxComplexity),
				Suggestion: "考虑拆分函数或简化逻辑，降低圈复杂度",
			})
		}

		return true
	})

	return violations
}

func (r *CyclomaticComplexityRule) calculateComplexity(body *ast.BlockStmt) int {
	complexity := 1

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			complexity++
			complexity += r.countBinaryExprComplexity(node.Cond)
		case *ast.ForStmt:
			complexity++
			complexity += r.countBinaryExprComplexity(node.Cond)
		case *ast.RangeStmt:
			complexity++
		case *ast.CaseClause:
			if len(node.List) > 0 {
				complexity += len(node.List)
			} else {
				complexity++
			}
		case *ast.SwitchStmt:
			complexity++
		case *ast.TypeSwitchStmt:
			complexity++
		case *ast.SelectStmt:
			complexity++
		case *ast.LabeledStmt:
			complexity++
		}
		return true
	})

	return complexity
}

func (r *CyclomaticComplexityRule) countBinaryExprComplexity(expr ast.Expr) int {
	count := 0

	ast.Inspect(expr, func(n ast.Node) bool {
		if binaryExpr, ok := n.(*ast.BinaryExpr); ok {
			if binaryExpr.Op == token.LAND || binaryExpr.Op == token.LOR {
				count++
			}
		}
		return true
	})

	return count
}

func (r *CyclomaticComplexityRule) determineSeverity(complexity int) types.Severity {
	if complexity > 30 {
		return types.SeverityCritical
	} else if complexity > 25 {
		return types.SeverityHigh
	} else if complexity > 20 {
		return types.SeverityMedium
	}
	return types.SeverityLow
}
