package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"strings"
)

type ResourceLeakRule struct {
	enabled bool
}

func NewResourceLeakRule() *ResourceLeakRule {
	return &ResourceLeakRule{
		enabled: true,
	}
}

func (r *ResourceLeakRule) Name() string {
	return "resource-leak"
}

func (r *ResourceLeakRule) Description() string {
	return "检测未释放的资源（文件、连接、锁等）"
}

func (r *ResourceLeakRule) Enabled() bool {
	return r.enabled
}

func (r *ResourceLeakRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		leaks := r.checkFunctionBody(funcDecl, fset, filePath)
		violations = append(violations, leaks...)

		return true
	})

	return violations
}

type ResourceAllocation struct {
	VarName      string
	ResourceType string
	Line         int
	IsHTTPBody   bool
}

type deferInfo struct {
	varName    string
	methodName string
	isClosure  bool
	closureVars []string
}

func (r *ResourceLeakRule) checkFunctionBody(funcDecl *ast.FuncDecl, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	allocations := r.collectResourceAllocations(funcDecl.Body, fset)

	deferCalls := r.collectDeferCalls(funcDecl.Body)

	closedVars := r.collectExplicitCloses(funcDecl.Body)

	for _, alloc := range allocations {
		if r.isResourceHandled(alloc, deferCalls, closedVars) {
			continue
		}

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

	return violations
}

func (r *ResourceLeakRule) collectResourceAllocations(body *ast.BlockStmt, fset *token.FileSet) []*ResourceAllocation {
	allocations := []*ResourceAllocation{}

	ast.Inspect(body, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, rhs := range assignStmt.Rhs {
			if callExpr, ok := rhs.(*ast.CallExpr); ok {
				resourceType, isHTTPBody := r.identifyResourceType(callExpr)
				if resourceType != "" && i < len(assignStmt.Lhs) {
					varName := r.extractVarName(assignStmt.Lhs[i])
					if varName != "" {
						allocations = append(allocations, &ResourceAllocation{
							VarName:      varName,
							ResourceType: resourceType,
							Line:         fset.Position(assignStmt.Pos()).Line,
							IsHTTPBody:   isHTTPBody,
						})
					}
				}
			}
		}

		return true
	})

	return allocations
}

func (r *ResourceLeakRule) identifyResourceType(callExpr *ast.CallExpr) (string, bool) {
	funcName := r.getFunctionName(callExpr)

	if strings.Contains(funcName, "os.Open") || strings.Contains(funcName, "os.Create") ||
		strings.Contains(funcName, "os.OpenFile") {
		return "file", false
	}

	if strings.Contains(funcName, "http.Get") || strings.Contains(funcName, "http.Post") ||
		strings.Contains(funcName, "http.Do") || strings.Contains(funcName, "http.Head") {
		return "http", true
	}

	if strings.Contains(funcName, "sql.Open") || strings.Contains(funcName, "sql.OpenDB") {
		return "database", false
	}

	if strings.Contains(funcName, "net.Dial") || strings.Contains(funcName, "net.DialTimeout") {
		return "network", false
	}

	if strings.Contains(funcName, "os.Exec") || strings.Contains(funcName, "exec.Command") {
		return "process", false
	}

	return "", false
}

func (r *ResourceLeakRule) getFunctionName(callExpr *ast.CallExpr) string {
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := fun.X.(*ast.Ident); ok {
			return ident.Name + "." + fun.Sel.Name
		}
		if selector, ok := fun.X.(*ast.SelectorExpr); ok {
			if ident, ok := selector.X.(*ast.Ident); ok {
				return ident.Name + "." + selector.Sel.Name + "." + fun.Sel.Name
			}
		}
	case *ast.Ident:
		return fun.Name
	}
	return ""
}

func (r *ResourceLeakRule) extractVarName(expr ast.Expr) string {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

func (r *ResourceLeakRule) collectDeferCalls(body *ast.BlockStmt) []deferInfo {
	deferInfos := []deferInfo{}

	ast.Inspect(body, func(n ast.Node) bool {
		deferStmt, ok := n.(*ast.DeferStmt)
		if !ok {
			return true
		}

		info := r.analyzeDeferCall(deferStmt.Call)
		if info.varName != "" || info.isClosure {
			deferInfos = append(deferInfos, info)
		}

		return true
	})

	return deferInfos
}

func (r *ResourceLeakRule) analyzeDeferCall(callExpr *ast.CallExpr) deferInfo {
	info := deferInfo{}

	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			info.varName = ident.Name
			info.methodName = selectorExpr.Sel.Name
			return info
		}
	}

	if funcLit, ok := callExpr.Fun.(*ast.FuncLit); ok {
		info.isClosure = true
		info.closureVars = r.extractClosedVarsInClosure(funcLit)
		return info
	}

	return info
}

func (r *ResourceLeakRule) extractClosedVarsInClosure(funcLit *ast.FuncLit) []string {
	vars := []string{}

	ast.Inspect(funcLit.Body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				if selectorExpr.Sel.Name == "Close" || selectorExpr.Sel.Name == "CloseBody" {
					vars = append(vars, ident.Name)
				}
			}
		}

		return true
	})

	return vars
}

func (r *ResourceLeakRule) collectExplicitCloses(body *ast.BlockStmt) []string {
	closedVars := []string{}

	ast.Inspect(body, func(n ast.Node) bool {
		exprStmt, ok := n.(*ast.ExprStmt)
		if !ok {
			return true
		}

		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			return true
		}

		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if selectorExpr.Sel.Name == "Close" {
				if ident, ok := selectorExpr.X.(*ast.Ident); ok {
					closedVars = append(closedVars, ident.Name)
				}
			}
		}

		return true
	})

	return closedVars
}

func (r *ResourceLeakRule) isResourceHandled(alloc *ResourceAllocation, deferCalls []deferInfo, closedVars []string) bool {
	for _, dv := range closedVars {
		if dv == alloc.VarName {
			return true
		}
	}

	for _, info := range deferCalls {
		if info.varName == alloc.VarName {
			return true
		}

		if info.isClosure {
			for _, cv := range info.closureVars {
				if cv == alloc.VarName {
					return true
				}
				if alloc.IsHTTPBody && cv == alloc.VarName+".Body" {
					return true
				}
			}
		}

		if alloc.IsHTTPBody && info.varName == alloc.VarName && info.methodName == "Close" {
			return true
		}
	}

	for _, info := range deferCalls {
		if alloc.IsHTTPBody && info.varName == alloc.VarName && info.methodName == "Body" {
			return true
		}
	}

	return false
}

func (r *ResourceLeakRule) determineSeverity(resourceType string) types.Severity {
	switch resourceType {
	case "database":
		return types.SeverityCritical
	case "http":
		return types.SeverityHigh
	case "file":
		return types.SeverityHigh
	case "network":
		return types.SeverityHigh
	case "process":
		return types.SeverityMedium
	default:
		return types.SeverityMedium
	}
}

func (r *ResourceLeakRule) generateMessage(alloc *ResourceAllocation) string {
	resourceName := ""
	switch alloc.ResourceType {
	case "file":
		resourceName = "文件句柄"
	case "http":
		resourceName = "HTTP响应"
	case "database":
		resourceName = "数据库连接"
	case "network":
		resourceName = "网络连接"
	case "process":
		resourceName = "进程"
	default:
		resourceName = "资源"
	}

	if alloc.IsHTTPBody {
		return "变量 '" + alloc.VarName + "' 包含HTTP响应，但未确保 Body 被正确关闭"
	}

	return "变量 '" + alloc.VarName + "' 分配了" + resourceName + "，但未确保资源被正确释放"
}

func (r *ResourceLeakRule) generateSuggestion(alloc *ResourceAllocation) string {
	if alloc.IsHTTPBody {
		return "在获取响应后立即添加 defer " + alloc.VarName + ".Body.Close() 以确保响应体被正确关闭"
	}

	closeMethod := ""
	switch alloc.ResourceType {
	case "file":
		closeMethod = alloc.VarName + ".Close()"
	case "http":
		closeMethod = alloc.VarName + ".Body.Close()"
	case "database":
		closeMethod = alloc.VarName + ".Close()"
	case "network":
		closeMethod = alloc.VarName + ".Close()"
	case "process":
		closeMethod = alloc.VarName + ".Wait()"
	default:
		closeMethod = alloc.VarName + ".Close()"
	}

	return "在资源分配后立即添加 defer " + closeMethod + " 以确保资源被正确释放"
}
