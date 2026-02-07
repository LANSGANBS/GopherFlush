package reporter

import (
	"encoding/json"
	"gopherflush/pkg/types"
	"os"
)

// JSONReporter JSON报告器
type JSONReporter struct {
	outputPath string
}

// NewJSONReporter 创建JSON报告器
func NewJSONReporter(outputPath string) *JSONReporter {
	return &JSONReporter{
		outputPath: outputPath,
	}
}

// Generate 生成JSON报告
func (r *JSONReporter) Generate(report *types.Report) error {
	// TODO: 实现JSON报告输出
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.outputPath, data, 0644)
}
