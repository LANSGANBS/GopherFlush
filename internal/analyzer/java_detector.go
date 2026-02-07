package analyzer

import (
	"regexp"
	"strings"
)

// Java 函数检测
var javaFuncPattern = regexp.MustCompile(`^\s*(?:public|private|protected|static|\s)+[\w<>\[\]]+\s+(\w+)\s*\(`)

func (ta *TextAnalyzer) detectJavaFunctions(lines []string) []*FunctionInfo {
	functions := []*FunctionInfo{}

	for i, line := range lines {
		if matches := javaFuncPattern.FindStringSubmatch(line); matches != nil {
			funcName := matches[1]
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

// Java 全局变量检测（类字段）
var javaFieldPattern = regexp.MustCompile(`^\s*(?:public|private|protected|static|\s)+[\w<>\[\]]+\s+(\w+)\s*[;=]`)

func (ta *TextAnalyzer) detectJavaGlobals(lines []string) []*GlobalVarInfo {
	globals := []*GlobalVarInfo{}
	inMethod := false
	braceLevel := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跟踪大括号层级
		braceLevel += strings.Count(line, "{") - strings.Count(line, "}")

		// 检测是否在方法内
		if javaFuncPattern.MatchString(line) {
			inMethod = true
		}

		// 如果大括号层级回到1，说明退出了方法
		if inMethod && braceLevel <= 1 {
			inMethod = false
		}

		// 只检测类级别的字段（不在方法内）
		if !inMethod && braceLevel == 1 && javaFieldPattern.MatchString(trimmed) {
			matches := javaFieldPattern.FindStringSubmatch(trimmed)
			if matches != nil {
				varName := matches[1]
				isExported := strings.Contains(line, "public")

				globals = append(globals, &GlobalVarInfo{
					Name:       varName,
					Line:       i + 1,
					IsExported: isExported,
				})
			}
		}
	}

	return globals
}
