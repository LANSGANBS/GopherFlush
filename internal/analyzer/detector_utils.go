package analyzer

import (
	"regexp"
	"strings"
	"unicode"
)

func CalculateIndent(line string) int {
	indent := 0
	for _, ch := range line {
		if ch == ' ' {
			indent++
		} else if ch == '\t' {
			indent += 4
		} else {
			break
		}
	}
	return indent
}

func NormalizeIndent(line string) (string, int) {
	trimmed := strings.TrimLeft(line, " \t")
	indent := len(line) - len(trimmed)
	return trimmed, indent
}

func IsBlankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

func IsCommentLine(line, language string) bool {
	trimmed := strings.TrimSpace(line)
	switch language {
	case "go":
		return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*")
	case "python":
		return strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, `"""`) || strings.HasPrefix(trimmed, `'''`)
	case "c", "cpp", "java":
		return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*")
	}
	return false
}

func CountBrackets(line string) (open, close int) {
	inString := false
	inChar := false
	escape := false

	for i, ch := range line {
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
				open++
			}
		case '}':
			if !inString && !inChar {
				close++
			}
		}

		_ = i
	}
	return open, close
}

func FindMatchingBracket(lines []string, startLine, startCol int) int {
	braceCount := 0
	foundStart := false

	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		startPos := 0
		if i == startLine {
			startPos = startCol
		}

		inString := false
		inChar := false
		escape := false

		for j, ch := range line[startPos:] {
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
						return i
					}
				}
			}
			_ = j
		}
	}

	return len(lines) - 1
}

func ExtractFunctionName(pattern *regexp.Regexp, line string) string {
	if pattern == nil {
		return ""
	}
	matches := pattern.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func IsKeyword(name string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "for": true, "while": true,
		"switch": true, "case": true, "default": true, "return": true,
		"break": true, "continue": true, "goto": true, "sizeof": true,
		"typedef": true, "struct": true, "enum": true, "union": true,
		"class": true, "interface": true, "package": true, "import": true,
		"try": true, "catch": true, "finally": true, "throw": true,
		"new": true, "delete": true, "this": true, "super": true,
		"def": true, "lambda": true, "yield": true,
		"async": true, "await": true, "with": true, "as": true,
	}
	return keywords[name]
}

func IsValidIdentifier(name string) bool {
	if name == "" {
		return false
	}

	if !unicode.IsLetter(rune(name[0])) && name[0] != '_' {
		return false
	}

	for _, ch := range name[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}

	return true
}

func StripComments(line, language string) string {
	switch language {
	case "c", "cpp", "java", "go":
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
	case "python":
		inString := false
		stringChar := byte(0)
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if inString {
				if ch == stringChar && (i == 0 || line[i-1] != '\\') {
					inString = false
				}
			} else {
				if ch == '"' || ch == '\'' {
					inString = true
					stringChar = ch
				} else if ch == '#' {
					return line[:i]
				}
			}
		}
	}
	return line
}
