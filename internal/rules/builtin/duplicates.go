package builtin

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"gopherflush/pkg/types"
	"strings"

	"github.com/cespare/xxhash/v2"
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

// FunctionInfo 函数信息
type FunctionInfo struct {
	Name    string
	Line    int
	Content string
	Hash    string
}

func (r *DuplicatesRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 收集所有函数
	functions := r.collectFunctions(file, fset)

	// 检测重复函数
	duplicates := r.findDuplicateFunctions(functions)

	// 生成违规记录
	for _, funcs := range duplicates {
		if len(funcs) < 2 {
			continue
		}

		// 收集函数名
		funcNames := []string{}
		for _, fn := range funcs {
			funcNames = append(funcNames, fn.Name)
		}

		// 确定严重程度
		severity := types.SeverityMedium
		if len(funcs) > 3 {
			severity = types.SeverityHigh
		}

		// 只为第一个函数创建违规记录
		fn := funcs[0]
		message := fmt.Sprintf("检测到 %d 个重复函数: %s", len(funcs), strings.Join(funcNames, ", "))

		violation := &types.Violation{
			RuleName:   r.Name(),
			Severity:   severity,
			FilePath:   filePath,
			Line:       fn.Line,
			Column:     1,
			Message:    message,
			Suggestion: "考虑提取公共逻辑到一个函数中，避免代码重复",
		}
		violations = append(violations, violation)
	}

	return violations
}

func (r *DuplicatesRule) Enabled() bool {
	return r.enabled
}

// collectFunctions 收集文件中的所有函数
func (r *DuplicatesRule) collectFunctions(file *ast.File, fset *token.FileSet) []*FunctionInfo {
	functions := []*FunctionInfo{}

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// 获取函数内容
		content := r.getFunctionContent(funcDecl, fset)

		// 标准化内容（去除空白和注释）
		normalized := r.normalizeContent(content)

		// 计算哈希值
		hash := r.calculateHash(normalized)

		funcInfo := &FunctionInfo{
			Name:    funcDecl.Name.Name,
			Line:    fset.Position(funcDecl.Pos()).Line,
			Content: content,
			Hash:    hash,
		}

		functions = append(functions, funcInfo)
		return true
	})

	return functions
}

// getFunctionContent 获取函数的内容
func (r *DuplicatesRule) getFunctionContent(funcDecl *ast.FuncDecl, fset *token.FileSet) string {
	if funcDecl.Body == nil {
		return ""
	}

	var buf strings.Builder
	printer.Fprint(&buf, fset, funcDecl.Body)
	return buf.String()
}

// normalizeContent 标准化内容（去除空白、换行等）
func (r *DuplicatesRule) normalizeContent(content string) string {
	// 去除所有空白字符
	normalized := strings.ReplaceAll(content, " ", "")
	normalized = strings.ReplaceAll(normalized, "\t", "")
	normalized = strings.ReplaceAll(normalized, "\n", "")
	normalized = strings.ReplaceAll(normalized, "\r", "")
	return normalized
}

// calculateHash 计算内容的哈希值
func (r *DuplicatesRule) calculateHash(content string) string {
	hash := xxhash.Sum64String(content)
	return fmt.Sprintf("%016x", hash)
}

// findDuplicateFunctions 查找重复的函数
func (r *DuplicatesRule) findDuplicateFunctions(functions []*FunctionInfo) map[string][]*FunctionInfo {
	hashMap := make(map[string][]*FunctionInfo)

	for _, fn := range functions {
		// 跳过空函数
		if fn.Hash == r.calculateHash("{}") || fn.Hash == r.calculateHash("") {
			continue
		}

		hashMap[fn.Hash] = append(hashMap[fn.Hash], fn)
	}

	// 只返回有重复的
	duplicates := make(map[string][]*FunctionInfo)
	for hash, funcs := range hashMap {
		if len(funcs) >= 2 {
			duplicates[hash] = funcs
		}
	}

	return duplicates
}
