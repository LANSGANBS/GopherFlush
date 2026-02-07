package analyzer

import (
	"path/filepath"
	"strings"
)

// Language 编程语言类型
type Language int

const (
	LanguageUnknown Language = iota
	LanguageGo
	LanguagePython
	LanguageJava
	LanguageC
	LanguageCPP
)

// String 返回语言名称
func (l Language) String() string {
	switch l {
	case LanguageGo:
		return "Go"
	case LanguagePython:
		return "Python"
	case LanguageJava:
		return "Java"
	case LanguageC:
		return "C"
	case LanguageCPP:
		return "C++"
	default:
		return "Unknown"
	}
}

// DetectLanguage 根据文件扩展名检测编程语言
func DetectLanguage(filePath string) Language {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".go":
		return LanguageGo
	case ".py":
		return LanguagePython
	case ".java":
		return LanguageJava
	case ".c", ".h":
		return LanguageC
	case ".cpp", ".cc", ".cxx", ".hpp", ".hh", ".hxx":
		return LanguageCPP
	default:
		return LanguageUnknown
	}
}
