package builtin

import (
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"regexp"
	"strings"
)

// CommentedCodeRule 注释掉的代码检测规则
type CommentedCodeRule struct {
	enabled bool
}

// NewCommentedCodeRule 创建注释掉的代码检测规则
func NewCommentedCodeRule() *CommentedCodeRule {
	return &CommentedCodeRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *CommentedCodeRule) Name() string {
	return "commented-code"
}

// Description 返回规则描述
func (r *CommentedCodeRule) Description() string {
	return "检测被注释掉的代码（而不是真正的注释）"
}

// Enabled 返回规则是否启用
func (r *CommentedCodeRule) Enabled() bool {
	return r.enabled
}

// Check 检查文件中的注释掉的代码
func (r *CommentedCodeRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 遍历所有注释组
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			// 检查注释内容是否是被注释掉的代码
			if r.isCommentedCode(comment.Text) {
				position := fset.Position(comment.Pos())

				// 确定严重程度
				severity := r.determineSeverity(comment.Text)

				violation := &types.Violation{
					RuleName:   r.Name(),
					Severity:   severity,
					FilePath:   filePath,
					Line:       position.Line,
					Column:     position.Column,
					Message:    "检测到被注释掉的代码，建议删除而不是注释",
					Suggestion: "如果代码不再需要，应该直接删除；如果需要保留历史，使用版本控制系统（如Git）",
				}

				violations = append(violations, violation)
			}
		}
	}

	return violations
}

// isCommentedCode 判断注释是否是被注释掉的代码
func (r *CommentedCodeRule) isCommentedCode(commentText string) bool {
	// 去除注释符号
	text := strings.TrimPrefix(commentText, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	// 排除空注释
	if text == "" {
		return false
	}

	// 排除常见的注释标记
	if r.isCommonCommentMarker(text) {
		return false
	}

	// 检查是否包含代码特征
	return r.hasCodePatterns(text)
}

// isCommonCommentMarker 判断是否是常见的注释标记
func (r *CommentedCodeRule) isCommonCommentMarker(text string) bool {
	// 常见的注释标记
	markers := []string{
		"TODO", "FIXME", "NOTE", "HACK", "XXX",
		"BUG", "OPTIMIZE", "REFACTOR", "REVIEW",
	}

	upperText := strings.ToUpper(text)
	for _, marker := range markers {
		if strings.HasPrefix(upperText, marker) {
			return true
		}
	}

	return false
}

// hasCodePatterns 判断是否包含代码模式
func (r *CommentedCodeRule) hasCodePatterns(text string) bool {
	// Go代码关键字模式
	codeKeywords := []string{
		`\bfunc\s+\w+\s*\(`,           // 函数定义
		`\bvar\s+\w+`,                 // 变量声明
		`\bconst\s+\w+`,               // 常量声明
		`\btype\s+\w+`,                // 类型定义
		`\bimport\s+`,                 // import语句
		`\bpackage\s+\w+`,             // package声明
		`\breturn\s+`,                 // return语句
		`\w+\s*:=\s*`,                 // 短变量声明
		`\w+\s*=\s*\w+\s*\(`,          // 赋值+函数调用
		`\bif\s+.*\{`,                 // if语句
		`\bfor\s+.*\{`,                // for循环
		`\}\s*$`,                      // 代码块结束
	}

	for _, pattern := range codeKeywords {
		matched, _ := regexp.MatchString(pattern, text)
		if matched {
			return true
		}
	}

	return false
}

// determineSeverity 确定严重程度
func (r *CommentedCodeRule) determineSeverity(commentText string) types.Severity {
	text := strings.TrimPrefix(commentText, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	// 被注释掉的函数定义 - 高
	if matched, _ := regexp.MatchString(`\bfunc\s+\w+\s*\(`, text); matched {
		return types.SeverityHigh
	}

	// 被注释掉的类型定义 - 中等
	if matched, _ := regexp.MatchString(`\btype\s+\w+`, text); matched {
		return types.SeverityMedium
	}

	// 被注释掉的变量/常量声明 - 中等
	if matched, _ := regexp.MatchString(`\b(var|const)\s+\w+`, text); matched {
		return types.SeverityMedium
	}

	// 其他被注释掉的代码 - 低
	return types.SeverityLow
}
