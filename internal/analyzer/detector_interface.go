package analyzer

type FunctionDetector interface {
	DetectFunctions(lines []string) []*FunctionInfo
	DetectGlobals(lines []string) []*GlobalVarInfo
}

type BaseDetector struct {
	patterns *DetectorPatterns
	analyzer *TextAnalyzer
}

func NewBaseDetector(patterns *DetectorPatterns, analyzer *TextAnalyzer) *BaseDetector {
	return &BaseDetector{
		patterns: patterns,
		analyzer: analyzer,
	}
}

func (d *BaseDetector) GetPatterns() *DetectorPatterns {
	return d.patterns
}

func (d *BaseDetector) GetAnalyzer() *TextAnalyzer {
	return d.analyzer
}
