package reporter

import (
	"fmt"
	"gopherflush/pkg/types"
)

// ConsoleReporter 控制台报告器
type ConsoleReporter struct{}

// fileStats 文件统计信息
type fileStats struct {
	total    int
	critical int
	high     int
	medium   int
	low      int
}

// NewConsoleReporter 创建控制台报告器
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// Generate 生成控制台报告
func (r *ConsoleReporter) Generate(report *types.Report) error {
	fmt.Println("\n========================================")
	fmt.Println("检测结果")
	fmt.Println("========================================")

	// 如果没有违规记录，显示成功信息
	if report.TotalViolations == 0 {
		successMsg := fmt.Sprintf("✓ 检测完成: 共检测 %d 个文件，未发现问题", report.TotalFiles)
		fmt.Println(Colorize(successMsg, ColorGreen))
		return nil
	}

	// 统计各严重程度的违规数量
	severityStats := r.calculateSeverityStats(report)

	// 显示总览
	fmt.Printf("检测完成: 共检测 %d 个文件，发现 %d 个问题\n\n",
		report.TotalFiles, report.TotalViolations)

	fmt.Println("严重程度统计:")
	fmt.Printf("  %s: %d\n", Colorize("极其严重", ColorRed), severityStats[types.SeverityCritical])
	fmt.Printf("  %s:     %d\n", Colorize("严重", ColorOrange), severityStats[types.SeverityHigh])
	fmt.Printf("  %s:     %d\n", Colorize("中等", ColorYellow), severityStats[types.SeverityMedium])
	fmt.Printf("  %s:       %d\n", Colorize("低", ColorGreen), severityStats[types.SeverityLow])

	// 按文件分组统计
	fileStats := r.calculateFileStats(report)

	fmt.Println("\n文件详情:")
	for filePath, stats := range fileStats {
		fmt.Printf("\n  %s: %d 个问题\n", filePath, stats.total)
		if stats.critical > 0 {
			fmt.Printf("    %s: %d\n", Colorize("极其严重", ColorRed), stats.critical)
		}
		if stats.high > 0 {
			fmt.Printf("    %s:     %d\n", Colorize("严重", ColorOrange), stats.high)
		}
		if stats.medium > 0 {
			fmt.Printf("    %s:     %d\n", Colorize("中等", ColorYellow), stats.medium)
		}
		if stats.low > 0 {
			fmt.Printf("    %s:       %d\n", Colorize("低", ColorGreen), stats.low)
		}
	}

	fmt.Println("\n========================================")
	return nil
}

// calculateSeverityStats 统计各严重程度的违规数量
func (r *ConsoleReporter) calculateSeverityStats(report *types.Report) map[types.Severity]int {
	stats := make(map[types.Severity]int)
	for _, violation := range report.Violations {
		stats[violation.Severity]++
	}
	return stats
}

// calculateFileStats 按文件分组统计违规
func (r *ConsoleReporter) calculateFileStats(report *types.Report) map[string]*fileStats {
	stats := make(map[string]*fileStats)

	for _, violation := range report.Violations {
		if _, exists := stats[violation.FilePath]; !exists {
			stats[violation.FilePath] = &fileStats{}
		}

		fs := stats[violation.FilePath]
		fs.total++

		switch violation.Severity {
		case types.SeverityCritical:
			fs.critical++
		case types.SeverityHigh:
			fs.high++
		case types.SeverityMedium:
			fs.medium++
		case types.SeverityLow:
			fs.low++
		}
	}

	return stats
}
