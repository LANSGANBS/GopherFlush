package reporter

import (
	"gopherflush/pkg/types"
)

// Reporter 报告生成器接口
type Reporter interface {
	// Generate 生成报告
	Generate(report *types.Report) error
}
