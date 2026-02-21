package analyzer

import (
	"regexp"
	"sync"
)

var (
	globalPatternRegistry *PatternRegistry
	patternOnce           sync.Once
)

func GetGlobalPatternRegistry() *PatternRegistry {
	patternOnce.Do(func() {
		globalPatternRegistry = NewPatternRegistry()
	})
	return globalPatternRegistry
}

type DetectorPatterns struct {
	FuncPattern   *regexp.Regexp
	GlobalPattern *regexp.Regexp
}

func NewCPatterns() *DetectorPatterns {
	return &DetectorPatterns{
		FuncPattern: regexp.MustCompile(
			`^\s*(?:(?:static|inline|extern)\s+)*` +
				`(?:(?:const|volatile|signed|unsigned|short|long)\s+)*` +
				`[\w\*\s]+?\s+(\w+)\s*\([^;]*\)\s*\{?`,
		),
		GlobalPattern: regexp.MustCompile(
			`^\s*(?:extern\s+)?(?:static\s+)?(?:const\s+|volatile\s+)*` +
				`(?:unsigned\s+|signed\s+)?(?:short\s+|long\s+)*[\w\*]+\s+(\w+)\s*(?:\[[^\]]*\])?\s*[;=]`,
		),
	}
}

func NewPythonPatterns() *DetectorPatterns {
	return &DetectorPatterns{
		FuncPattern: regexp.MustCompile(
			`^\s*(?:(?:@\w+(?:\([^)]*\))?\s*)*)\s*(?:async\s+)?def\s+(\w+)\s*\(`,
		),
		GlobalPattern: regexp.MustCompile(
			`^(\w+)\s*=`,
		),
	}
}

func NewJavaPatterns() *DetectorPatterns {
	return &DetectorPatterns{
		FuncPattern: regexp.MustCompile(
			`^\s*(?:(?:public|private|protected|static|final|abstract|synchronized|native|strictfp)\s+)*` +
				`(?:<[^>]+>\s+)?[\w<>\[\],\s\?]+?\s+(\w+)\s*\([^)]*\)\s*(?:throws\s+[\w\s,]+)?\s*\{?`,
		),
		GlobalPattern: regexp.MustCompile(
			`^\s*(?:(?:public|private|protected|static|final|transient|volatile)\s+)*` +
				`[\w<>\[\],\s\?]+?\s+(\w+)\s*(?:\[[^\]]*\])?\s*[;=]`,
		),
	}
}

type PatternRegistry struct {
	patterns map[Language]*DetectorPatterns
}

func NewPatternRegistry() *PatternRegistry {
	return &PatternRegistry{
		patterns: map[Language]*DetectorPatterns{
			LanguageC:      NewCPatterns(),
			LanguageCPP:    NewCPatterns(),
			LanguagePython: NewPythonPatterns(),
			LanguageJava:   NewJavaPatterns(),
		},
	}
}

func (pr *PatternRegistry) Get(lang Language) *DetectorPatterns {
	if patterns, exists := pr.patterns[lang]; exists {
		return patterns
	}
	return nil
}
