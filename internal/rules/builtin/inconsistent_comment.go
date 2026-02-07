package builtin

import (
	"bufio"
	"go/ast"
	"go/token"
	"gopherflush/pkg/types"
	"os"
	"regexp"
	"strings"
)

// InconsistentCommentRule 不一致注释检测规则
type InconsistentCommentRule struct {
	enabled bool
}

// NewInconsistentCommentRule 创建不一致注释检测规则
func NewInconsistentCommentRule() *InconsistentCommentRule {
	return &InconsistentCommentRule{
		enabled: true,
	}
}

// Name 返回规则名称
func (r *InconsistentCommentRule) Name() string {
	return "inconsistent-comment"
}

// Description 返回规则描述
func (r *InconsistentCommentRule) Description() string {
	return "检测与代码不一致的注释"
}

// Enabled 返回规则是否启用
func (r *InconsistentCommentRule) Enabled() bool {
	return r.enabled
}

// Check 检查文件中的不一致注释
func (r *InconsistentCommentRule) Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation {
	violations := []*types.Violation{}

	// 读取源文件内容
	sourceLines := r.readSourceLines(filePath)

	// 遍历所有注释
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			position := fset.Position(comment.Pos())
			commentText := comment.Text

			// 获取注释所在行的代码
			codeLine := r.getCodeLine(sourceLines, position.Line, commentText)
			if codeLine == "" {
				continue
			}

			// 检查注释与代码是否一致
			if inconsistency := r.checkConsistency(commentText, codeLine); inconsistency != "" {
				violation := &types.Violation{
					RuleName:   r.Name(),
					Severity:   types.SeverityMedium,
					FilePath:   filePath,
					Line:       position.Line,
					Column:     position.Column,
					Message:    "检测到注释与代码不一致：" + inconsistency,
					Suggestion: "更新注释以匹配代码的实际行为，或修正代码以匹配注释的描述",
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

// readSourceLines 读取源文件的所有行
func (r *InconsistentCommentRule) readSourceLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

// getCodeLine 获取注释所在行的代码（去除注释部分）
func (r *InconsistentCommentRule) getCodeLine(sourceLines []string, lineNum int, commentText string) string {
	if lineNum < 1 || lineNum > len(sourceLines) {
		return ""
	}

	line := sourceLines[lineNum-1] // 行号从1开始，数组索引从0开始

	// 去除注释部分
	commentStart := strings.Index(line, "//")
	if commentStart >= 0 {
		line = line[:commentStart]
	}

	// 去除块注释
	commentStart = strings.Index(line, "/*")
	if commentStart >= 0 {
		line = line[:commentStart]
	}

	codeLine := strings.TrimSpace(line)

	// 如果当前行没有代码（注释在单独一行），查找下一行
	if codeLine == "" && lineNum < len(sourceLines) {
		nextLine := sourceLines[lineNum] // 下一行
		return strings.TrimSpace(nextLine)
	}

	return codeLine
}

// checkConsistency 检查注释与代码是否一致
func (r *InconsistentCommentRule) checkConsistency(commentText, codeLine string) string {
	// 去除注释符号
	comment := strings.TrimPrefix(commentText, "//")
	comment = strings.TrimPrefix(comment, "/*")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimSpace(comment)

	// 如果代码行为空，跳过
	if codeLine == "" {
		return ""
	}

	// 检查数字不一致
	if inconsistency := r.checkNumberInconsistency(comment, codeLine); inconsistency != "" {
		return inconsistency
	}

	return ""
}

// checkNumberInconsistency 检查注释和代码中的数字是否一致
func (r *InconsistentCommentRule) checkNumberInconsistency(comment, codeLine string) string {
	// 提取注释中的数字
	commentNumbers := r.extractNumbers(comment)
	if len(commentNumbers) == 0 {
		return ""
	}

	// 提取代码中的数字
	codeNumbers := r.extractNumbers(codeLine)
	if len(codeNumbers) == 0 {
		return ""
	}

	// 检查注释中的数字是否在代码中出现
	for _, commentNum := range commentNumbers {
		found := false
		for _, codeNum := range codeNumbers {
			if commentNum == codeNum {
				found = true
				break
			}
		}
		if !found {
			return "注释中提到的数字 '" + commentNum + "' 与代码中的实际值不匹配"
		}
	}

	return ""
}

// extractNumbers 从文本中提取所有数字
func (r *InconsistentCommentRule) extractNumbers(text string) []string {
	// 匹配整数和小数
	numberPattern := regexp.MustCompile(`\d+\.?\d*`)
	matches := numberPattern.FindAllString(text, -1)

	// 去重
	seen := make(map[string]bool)
	result := []string{}
	for _, match := range matches {
		if !seen[match] {
			seen[match] = true
			result = append(result, match)
		}
	}

	return result
}
