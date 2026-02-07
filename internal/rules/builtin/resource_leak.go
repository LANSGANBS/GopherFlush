package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"strings"
)

// ResourceLeakRule 资源泄漏检测规则
type ResourceLeakRule struct {
	enabled bool
}

// NewResourceLeakRule 创建资源泄漏检测规则
func NewResourceLeakRule() *ResourceLeakRule {
	return &ResourceLeakRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *ResourceLeakRule) Name() string {
	return "resource-leak"
}

// Description 返回规则描述
func (r *ResourceLeakRule) Description() string {
	return "检测未释放的资源（文件、连接、锁等）"
}

// Enabled 返回规则是否启用
func (r *ResourceLeakRule) Enabled() bool {
	return r.enabled
}

// Check 检查文件中的资源泄漏
func (r *ResourceLeakRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 遍历所有函数
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		// 检查函数体中的资源泄漏
		leaks := r.checkFunctionBody(funcDecl, fset, filePath)
		violations = append(violations, leaks...)

		return true
	})

	return violations
}

// ResourceAllocation 资源分配信息
type ResourceAllocation struct {
	VarName      string // 变量名
	ResourceType string // 资源类型（file, http, db, mutex等）
	Line         int    // 行号
	HasDefer     bool   // 是否有defer释放
}

// checkFunctionBody 检查函数体中的资源泄漏
func (r *ResourceLeakRule) checkFunctionBody(funcDecl *ast.FuncDecl, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 收集资源分配
	allocations := r.collectResourceAllocations(funcDecl.Body, fset)

	// 收集defer语句
	deferCalls := r.collectDeferCalls(funcDecl.Body)

	// 检查每个资源分配是否有对应的defer释放
	for _, alloc := range allocations {
		if !r.hasMatchingDefer(alloc, deferCalls) {
			severity := r.determineSeverity(alloc.ResourceType)

			violation := &types.Violation{
				RuleName:   r.Name(),
				Severity:   severity,
				FilePath:   filePath,
				Line:       alloc.Line,
				Column:     1,
				Message:    r.generateMessage(alloc),
				Suggestion: r.generateSuggestion(alloc),
			}
			violations = append(violations, violation)
		}
	}

	return violations
}

// collectResourceAllocations 收集资源分配
func (r *ResourceLeakRule) collectResourceAllocations(body *ast.BlockStmt, fset *token.FileSet) []*ResourceAllocation {
	allocations := []*ResourceAllocation{}

	ast.Inspect(body, func(n ast.Node) bool {
		// 查找赋值语句
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// 检查右侧是否是资源分配调用
		for i, rhs := range assignStmt.Rhs {
			if callExpr, ok := rhs.(*ast.CallExpr); ok {
				resourceType := r.identifyResourceType(callExpr)
				if resourceType != "" && i < len(assignStmt.Lhs) {
					// 获取变量名
					varName := r.extractVarName(assignStmt.Lhs[i])
					if varName != "" {
						allocations = append(allocations, &ResourceAllocation{
							VarName:      varName,
							ResourceType: resourceType,
							Line:         fset.Position(assignStmt.Pos()).Line,
							HasDefer:     false,
						})
					}
				}
			}
		}

		return true
	})

	return allocations
}

// identifyResourceType 识别资源类型
func (r *ResourceLeakRule) identifyResourceType(callExpr *ast.CallExpr) string {
	// 获取函数调用的名称
	funcName := r.getFunctionName(callExpr)

	// 文件操作
	if strings.Contains(funcName, "os.Open") || strings.Contains(funcName, "os.Create") {
		return "file"
	}

	// HTTP请求
	if strings.Contains(funcName, "http.Get") || strings.Contains(funcName, "http.Post") ||
	   strings.Contains(funcName, "http.Do") {
		return "http"
	}

	// 数据库连接
	if strings.Contains(funcName, "sql.Open") {
		return "database"
	}

	return ""
}

// getFunctionName 获取函数调用的名称
func (r *ResourceLeakRule) getFunctionName(callExpr *ast.CallExpr) string {
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		// 例如：os.Open, http.Get
		if ident, ok := fun.X.(*ast.Ident); ok {
			return ident.Name + "." + fun.Sel.Name
		}
	case *ast.Ident:
		// 例如：Open
		return fun.Name
	}
	return ""
}

// extractVarName 提取变量名
func (r *ResourceLeakRule) extractVarName(expr ast.Expr) string {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// collectDeferCalls 收集defer语句中释放的资源
func (r *ResourceLeakRule) collectDeferCalls(body *ast.BlockStmt) []string {
	deferVars := []string{}

	ast.Inspect(body, func(n ast.Node) bool {
		// 查找defer语句
		deferStmt, ok := n.(*ast.DeferStmt)
		if !ok {
			return true
		}

		// 提取被defer的变量名
		varName := r.extractDeferVarName(deferStmt.Call)
		if varName != "" {
			deferVars = append(deferVars, varName)
		}

		return true
	})

	return deferVars
}

// extractDeferVarName 从defer调用中提取变量名
func (r *ResourceLeakRule) extractDeferVarName(callExpr *ast.CallExpr) string {
	// 处理 defer file.Close() 这种形式
	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// hasMatchingDefer 检查资源分配是否有对应的defer释放
func (r *ResourceLeakRule) hasMatchingDefer(alloc *ResourceAllocation, deferVars []string) bool {
	for _, deferVar := range deferVars {
		if deferVar == alloc.VarName {
			return true
		}
	}
	return false
}

// determineSeverity 根据资源类型确定严重程度
func (r *ResourceLeakRule) determineSeverity(resourceType string) types.Severity {
	switch resourceType {
	case "database":
		return types.SeverityCritical // 数据库连接泄漏非常严重
	case "http":
		return types.SeverityHigh // HTTP连接泄漏严重
	case "file":
		return types.SeverityHigh // 文件描述符泄漏严重
	default:
		return types.SeverityMedium
	}
}

// generateMessage 生成违规消息
func (r *ResourceLeakRule) generateMessage(alloc *ResourceAllocation) string {
	resourceName := ""
	switch alloc.ResourceType {
	case "file":
		resourceName = "文件"
	case "http":
		resourceName = "HTTP响应"
	case "database":
		resourceName = "数据库连接"
	default:
		resourceName = "资源"
	}

	return "变量 '" + alloc.VarName + "' 分配了" + resourceName + "但没有使用defer进行释放"
}

// generateSuggestion 生成修复建议
func (r *ResourceLeakRule) generateSuggestion(alloc *ResourceAllocation) string {
	closeMethod := ""
	switch alloc.ResourceType {
	case "file":
		closeMethod = alloc.VarName + ".Close()"
	case "http":
		closeMethod = alloc.VarName + ".Body.Close()"
	case "database":
		closeMethod = alloc.VarName + ".Close()"
	default:
		closeMethod = alloc.VarName + ".Close()"
	}

	return "在资源分配后立即添加 defer " + closeMethod + " 以确保资源被正确释放"
}
