package analyzer

import (
	"regexp"
	"strings"
)

// C/C++ 函数检测
var cFuncPattern = regexp.MustCompile(`^\s*[\w\*\s]+\s+(\w+)\s*\([^)]*\)\s*\{?`)

func (ta *TextAnalyzer) detectCFunctions(lines []string) []*FunctionInfo {
	functions := []*FunctionInfo{}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过预处理指令和注释
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if matches := cFuncPattern.FindStringSubmatch(line); matches != nil {
			funcName := matches[1]

			// 排除一些常见的非函数模式
			if funcName == "if" || funcName == "while" || funcName == "for" ||
			   funcName == "switch" || funcName == "return" {
				continue
			}

			startLine := i + 1

			// 检测函数结束位置（通过大括号）
			endLine := ta.findBracketEnd(lines, i)

			functions = append(functions, &FunctionInfo{
				Name:      funcName,
				StartLine: startLine,
				EndLine:   endLine,
				LineCount: endLine - startLine + 1,
			})
		}
	}

	return functions
}

// C/C++ 全局变量检测
var cGlobalPattern = regexp.MustCompile(`^\s*(?:extern\s+)?(?:static\s+)?[\w\*]+\s+(\w+)\s*[;=]`)

func (ta *TextAnalyzer) detectCGlobals(lines []string) []*GlobalVarInfo {
	globals := []*GlobalVarInfo{}
	inFunction := false
	braceLevel := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过预处理指令和注释
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		// 跟踪大括号层级
		braceLevel += strings.Count(line, "{") - strings.Count(line, "}")

		// 检测是否在函数内
		if cFuncPattern.MatchString(line) && strings.Contains(line, "{") {
			inFunction = true
		}

		// 如果大括号层级回到0，说明退出了函数
		if inFunction && braceLevel == 0 {
			inFunction = false
		}

		// 只检测文件级别的变量（不在函数内）
		if !inFunction && braceLevel == 0 && cGlobalPattern.MatchString(trimmed) {
			matches := cGlobalPattern.FindStringSubmatch(trimmed)
			if matches != nil {
				varName := matches[1]

				// 排除一些常见的非变量模式
				if varName != "typedef" && varName != "struct" && varName != "enum" {
					globals = append(globals, &GlobalVarInfo{
						Name:       varName,
						Line:       i + 1,
						IsExported: !strings.Contains(line, "static"),
					})
				}
			}
		}
	}

	return globals
}

// findBracketEnd 查找大括号结束位置（通用方法）
func (ta *TextAnalyzer) findBracketEnd(lines []string, startIdx int) int {
	braceCount := 0
	foundStart := false

	for i := startIdx; i < len(lines); i++ {
		line := lines[i]

		for _, ch := range line {
			if ch == '{' {
				braceCount++
				foundStart = true
			} else if ch == '}' {
				braceCount--
				if foundStart && braceCount == 0 {
					return i + 1
				}
			}
		}
	}

	return len(lines)
}
