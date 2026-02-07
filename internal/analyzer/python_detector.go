package analyzer

import (
	"regexp"
	"strings"
)

// Python 函数检测
var pythonFuncPattern = regexp.MustCompile(`^\s*def\s+(\w+)\s*\(`)

func (ta *TextAnalyzer) detectPythonFunctions(lines []string) []*FunctionInfo {
	functions := []*FunctionInfo{}

	for i, line := range lines {
		if matches := pythonFuncPattern.FindStringSubmatch(line); matches != nil {
			funcName := matches[1]
			startLine := i + 1

			// 检测函数结束位置（通过缩进）
			endLine := ta.findPythonFunctionEnd(lines, i)

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

func (ta *TextAnalyzer) findPythonFunctionEnd(lines []string, startIdx int) int {
	if startIdx >= len(lines) {
		return startIdx + 1
	}

	// 获取函数定义的缩进级别
	defLine := lines[startIdx]
	baseIndent := len(defLine) - len(strings.TrimLeft(defLine, " \t"))

	// 查找函数结束位置
	for i := startIdx + 1; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], " \t\r\n")

		// 跳过空行和注释
		if line == "" || strings.TrimSpace(line) == "" {
			continue
		}

		// 计算当前行的缩进
		currentIndent := len(lines[i]) - len(strings.TrimLeft(lines[i], " \t"))

		// 如果缩进回到或小于函数定义级别，说明函数结束
		if currentIndent <= baseIndent {
			return i
		}
	}

	return len(lines)
}

// Python 全局变量检测
var pythonGlobalPattern = regexp.MustCompile(`^(\w+)\s*=`)

func (ta *TextAnalyzer) detectPythonGlobals(lines []string) []*GlobalVarInfo {
	globals := []*GlobalVarInfo{}
	inFunction := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 检测是否进入函数
		if pythonFuncPattern.MatchString(line) {
			inFunction = true
			continue
		}

		// 检测是否退出函数（缩进回到顶层）
		if inFunction && len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
			inFunction = false
		}

		// 只检测顶层的变量赋值
		if !inFunction && pythonGlobalPattern.MatchString(trimmed) {
			matches := pythonGlobalPattern.FindStringSubmatch(trimmed)
			if matches != nil {
				varName := matches[1]
				// 排除一些常见的非全局变量模式
				if !strings.HasPrefix(trimmed, "if ") &&
				   !strings.HasPrefix(trimmed, "for ") &&
				   !strings.HasPrefix(trimmed, "while ") {
					globals = append(globals, &GlobalVarInfo{
						Name:       varName,
						Line:       i + 1,
						IsExported: true, // Python 没有明确的导出概念
					})
				}
			}
		}
	}

	return globals
}
