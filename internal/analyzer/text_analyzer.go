package analyzer

import (
	"bufio"
	"os"
)

type TextAnalyzer struct {
	language       Language
	patterns       *DetectorPatterns
	patternRegistry *PatternRegistry
}

func NewTextAnalyzer(language Language) *TextAnalyzer {
	registry := NewPatternRegistry()
	return &TextAnalyzer{
		language:        language,
		patterns:        registry.Get(language),
		patternRegistry: registry,
	}
}

func NewTextAnalyzerWithRegistry(language Language, registry *PatternRegistry) *TextAnalyzer {
	return &TextAnalyzer{
		language:        language,
		patterns:        registry.Get(language),
		patternRegistry: registry,
	}
}

func (ta *TextAnalyzer) AnalyzeFile(filePath string) (*FileInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := &FileInfo{
		Path:      filePath,
		Language:  ta.language,
		Lines:     []string{},
		Functions: []*FunctionInfo{},
		Globals:   []*GlobalVarInfo{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		info.Lines = append(info.Lines, line)
	}

	info.Functions = ta.detectFunctions(info.Lines)
	info.Globals = ta.detectGlobalVars(info.Lines)

	return info, scanner.Err()
}

func (ta *TextAnalyzer) GetPatterns() *DetectorPatterns {
	return ta.patterns
}

func (ta *TextAnalyzer) GetLanguage() Language {
	return ta.language
}

type FileInfo struct {
	Path      string
	Language  Language
	Lines     []string
	Functions []*FunctionInfo
	Globals   []*GlobalVarInfo
}

type FunctionInfo struct {
	Name      string
	StartLine int
	EndLine   int
	LineCount int
}

type GlobalVarInfo struct {
	Name       string
	Line       int
	IsExported bool
}
