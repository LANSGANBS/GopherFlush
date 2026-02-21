package analyzer

import (
	"strings"
)

func (ta *TextAnalyzer) detectJavaFunctions(lines []string) []*FunctionInfo {
	if ta.patterns == nil || ta.patterns.FuncPattern == nil {
		return nil
	}

	functions := []*FunctionInfo{}
	pattern := ta.patterns.FuncPattern

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
			continue
		}

		cleanLine := StripComments(line, "java")

		if matches := pattern.FindStringSubmatch(cleanLine); matches != nil {
			funcName := matches[1]

			if IsKeyword(funcName) {
				continue
			}

			if !IsValidIdentifier(funcName) {
				continue
			}

			if strings.Contains(cleanLine, "class ") && strings.Contains(cleanLine, "{") {
				if !strings.Contains(cleanLine, "(") {
					continue
				}
			}

			startLine := i + 1
			endLine := ta.findBracketEndEnhanced(lines, i)

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

func (ta *TextAnalyzer) detectJavaGlobals(lines []string) []*GlobalVarInfo {
	if ta.patterns == nil || ta.patterns.GlobalPattern == nil {
		return nil
	}

	globals := []*GlobalVarInfo{}
	pattern := ta.patterns.GlobalPattern

	braceStack := []int{}
	inMethod := false
	classDepth := 0
	inClass := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
			continue
		}

		cleanLine := StripComments(line, "java")

		if strings.Contains(cleanLine, "class ") || strings.Contains(cleanLine, "interface ") ||
			strings.Contains(cleanLine, "enum ") || strings.Contains(cleanLine, "@interface ") {
			for _, ch := range cleanLine {
				if ch == '{' {
					classDepth++
					inClass = true
				}
			}
		}

		openBraces := strings.Count(cleanLine, "{")
		closeBraces := strings.Count(cleanLine, "}")

		for j := 0; j < openBraces; j++ {
			if len(braceStack) == 0 && ta.patterns.FuncPattern != nil {
				if ta.patterns.FuncPattern.MatchString(cleanLine) {
					inMethod = true
				}
			}
			braceStack = append(braceStack, i)
		}

		for j := 0; j < closeBraces; j++ {
			if len(braceStack) > 0 {
				braceStack = braceStack[:len(braceStack)-1]
			}
			if len(braceStack) <= 1 {
				inMethod = false
			}
			if len(braceStack) == 0 {
				classDepth--
				if classDepth <= 0 {
					inClass = false
					classDepth = 0
				}
			}
		}

		if inClass && !inMethod && len(braceStack) == 1 {
			if matches := pattern.FindStringSubmatch(trimmed); matches != nil {
				varName := matches[1]

				if IsKeyword(varName) {
					continue
				}

				if !IsValidIdentifier(varName) {
					continue
				}

				if strings.Contains(cleanLine, "(") && strings.Contains(cleanLine, ")") {
					continue
				}

				globals = append(globals, &GlobalVarInfo{
					Name:       varName,
					Line:       i + 1,
					IsExported: strings.Contains(line, "public"),
				})
			}
		}
	}

	return globals
}
