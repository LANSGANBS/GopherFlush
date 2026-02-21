package analyzer

import (
	"strings"
)

func (ta *TextAnalyzer) detectPythonFunctions(lines []string) []*FunctionInfo {
	if ta.patterns == nil || ta.patterns.FuncPattern == nil {
		return nil
	}

	functions := []*FunctionInfo{}
	pattern := ta.patterns.FuncPattern

	for i, line := range lines {
		if matches := pattern.FindStringSubmatch(line); matches != nil {
			funcName := matches[1]

			if IsKeyword(funcName) {
				continue
			}

			if !IsValidIdentifier(funcName) {
				continue
			}

			startLine := i + 1
			endLine := ta.findPythonFunctionEndEnhanced(lines, i)

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

func (ta *TextAnalyzer) findPythonFunctionEndEnhanced(lines []string, startIdx int) int {
	if startIdx >= len(lines) {
		return startIdx + 1
	}

	defLine := lines[startIdx]
	baseIndent := CalculateIndent(defLine)

	funcBodyIndent := -1
	foundBody := false

	for i := startIdx + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimRight(line, " \t\r\n")

		if trimmed == "" || strings.TrimSpace(trimmed) == "" {
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(trimmed), "#") {
			continue
		}

		currentIndent := CalculateIndent(line)

		if !foundBody {
			if currentIndent > baseIndent {
				funcBodyIndent = currentIndent
				foundBody = true
			} else if currentIndent <= baseIndent {
				return i
			}
			continue
		}

		if currentIndent <= baseIndent {
			if strings.HasPrefix(trimmed, "@") {
				continue
			}
			return i
		}

		if currentIndent < funcBodyIndent && !strings.HasPrefix(trimmed, "else:") &&
			!strings.HasPrefix(trimmed, "elif ") && !strings.HasPrefix(trimmed, "except") &&
			!strings.HasPrefix(trimmed, "finally:") {
			if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "for ") ||
				strings.HasPrefix(trimmed, "while ") || strings.HasPrefix(trimmed, "try:") ||
				strings.HasPrefix(trimmed, "with ") || strings.HasPrefix(trimmed, "class ") {
				continue
			}
		}
	}

	return len(lines)
}

func (ta *TextAnalyzer) detectPythonGlobals(lines []string) []*GlobalVarInfo {
	if ta.patterns == nil || ta.patterns.GlobalPattern == nil {
		return nil
	}

	globals := []*GlobalVarInfo{}
	pattern := ta.patterns.GlobalPattern

	indentStack := []int{0}
	inFunction := false
	functionIndent := -1
	inClass := false
	classIndent := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		currentIndent := CalculateIndent(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "def ") {
			if currentIndent == 0 || (inClass && currentIndent > classIndent) {
				inFunction = true
				functionIndent = currentIndent
			}
			continue
		}

		if strings.HasPrefix(trimmed, "class ") {
			inClass = true
			classIndent = currentIndent
			continue
		}

		if inFunction && currentIndent <= functionIndent {
			inFunction = false
			functionIndent = -1
		}

		if inClass && currentIndent <= classIndent {
			inClass = false
			classIndent = -1
		}

		if !inFunction && currentIndent == 0 {
			if matches := pattern.FindStringSubmatch(trimmed); matches != nil {
				varName := matches[1]

				if IsKeyword(varName) {
					continue
				}

				if !IsValidIdentifier(varName) {
					continue
				}

				if strings.HasPrefix(trimmed, "if ") ||
					strings.HasPrefix(trimmed, "for ") ||
					strings.HasPrefix(trimmed, "while ") ||
					strings.HasPrefix(trimmed, "with ") ||
					strings.HasPrefix(trimmed, "try:") ||
					strings.HasPrefix(trimmed, "except") ||
					strings.HasPrefix(trimmed, "elif ") ||
					strings.HasPrefix(trimmed, "else:") {
					continue
				}

				globals = append(globals, &GlobalVarInfo{
					Name:       varName,
					Line:       i + 1,
					IsExported: !strings.HasPrefix(varName, "_"),
				})
			}
		}

		indentStack = append(indentStack, currentIndent)
		if len(indentStack) > 100 {
			indentStack = indentStack[1:]
		}
	}

	return globals
}
