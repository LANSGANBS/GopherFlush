package analyzer

import (
	"bufio"
	"os"
)

// TextAnalyzer 通用文本分析器（用于非Go语言）
type TextAnalyzer struct {
	language Language
}

// NewTextAnalyzer 创建文本分析器
func NewTextAnalyzer(language Language) *TextAnalyzer {
	return &TextAnalyzer{
		language: language,
	}
}

// AnalyzeFile 分析文件
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
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		info.Lines = append(info.Lines, line)
	}

	// 检测函数
	info.Functions = ta.detectFunctions(info.Lines)

	// 检测全局变量
	info.Globals = ta.detectGlobalVars(info.Lines)

	return info, scanner.Err()
}

// FileInfo 文件信息
type FileInfo struct {
	Path      string
	Language  Language
	Lines     []string
	Functions []*FunctionInfo
	Globals   []*GlobalVarInfo
}

// FunctionInfo 函数信息
type FunctionInfo struct {
	Name      string
	StartLine int
	EndLine   int
	LineCount int
}

// GlobalVarInfo 全局变量信息
type GlobalVarInfo struct {
	Name       string
	Line       int
	IsExported bool
}
