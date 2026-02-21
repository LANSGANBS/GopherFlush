package analyzer

import (
	"strings"
)

func (ta *TextAnalyzer) detectCFunctions(lines []string) []*FunctionInfo {
	if ta.patterns == nil || ta.patterns.FuncPattern == nil {
		return nil
	}

	functions := []*FunctionInfo{}
	pattern := ta.patterns.FuncPattern

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		cleanLine := StripComments(line, "c")

		if matches := pattern.FindStringSubmatch(cleanLine); matches != nil {
			funcName := matches[1]

			if IsKeyword(funcName) {
				continue
			}

			if !IsValidIdentifier(funcName) {
				continue
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

func (ta *TextAnalyzer) detectCGlobals(lines []string) []*GlobalVarInfo {
	if ta.patterns == nil || ta.patterns.GlobalPattern == nil {
		return nil
	}

	globals := []*GlobalVarInfo{}
	pattern := ta.patterns.GlobalPattern

	braceStack := []int{}
	inFunction := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		cleanLine := StripComments(line, "c")

		openBraces := strings.Count(cleanLine, "{")
		closeBraces := strings.Count(cleanLine, "}")

		for j := 0; j < openBraces; j++ {
			if len(braceStack) == 0 && ta.patterns.FuncPattern != nil {
				if ta.patterns.FuncPattern.MatchString(cleanLine) || (i > 0 && ta.patterns.FuncPattern.MatchString(lines[i-1])) {
					inFunction = true
				}
			}
			braceStack = append(braceStack, i)
		}

		for j := 0; j < closeBraces; j++ {
			if len(braceStack) > 0 {
				braceStack = braceStack[:len(braceStack)-1]
			}
			if len(braceStack) == 0 {
				inFunction = false
			}
		}

		if !inFunction && len(braceStack) == 0 {
			if matches := pattern.FindStringSubmatch(trimmed); matches != nil {
				varName := matches[1]

				if IsKeyword(varName) {
					continue
				}

				if !IsValidIdentifier(varName) {
					continue
				}

				globals = append(globals, &GlobalVarInfo{
					Name:       varName,
					Line:       i + 1,
					IsExported: !strings.Contains(line, "static"),
				})
			}
		}
	}

	return globals
}

func (ta *TextAnalyzer) findBracketEndEnhanced(lines []string, startIdx int) int {
	braceCount := 0
	foundStart := false
	inMultiLineComment := false

	for i := startIdx; i < len(lines); i++ {
		line := lines[i]

		processedLine := ""
		j := 0
		for j < len(line) {
			if inMultiLineComment {
				if j+1 < len(line) && line[j:j+2] == "*/" {
					inMultiLineComment = false
					j += 2
					continue
				}
				j++
				continue
			}

			if j+1 < len(line) && line[j:j+2] == "/*" {
				inMultiLineComment = true
				j += 2
				continue
			}

			if j+1 < len(line) && line[j:j+2] == "//" {
				break
			}

			processedLine += string(line[j])
			j++
		}

		inString := false
		inChar := false
		escape := false

		for _, ch := range processedLine {
			if escape {
				escape = false
				continue
			}

			switch ch {
			case '\\':
				if inString || inChar {
					escape = true
				}
			case '"':
				if !inChar {
					inString = !inString
				}
			case '\'':
				if !inString {
					inChar = !inChar
				}
			case '{':
				if !inString && !inChar {
					braceCount++
					foundStart = true
				}
			case '}':
				if !inString && !inChar {
					braceCount--
					if foundStart && braceCount == 0 {
						return i + 1
					}
				}
			}
		}
	}

	return len(lines)
}
