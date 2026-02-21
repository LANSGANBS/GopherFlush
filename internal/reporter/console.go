package reporter

import (
	"fmt"
	"gopherflush/pkg/types"
	"strings"
)

type ConsoleReporter struct{}

type fileStats struct {
	total    int
	critical int
	high     int
	medium   int
	low      int
}

func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

func (r *ConsoleReporter) Generate(report *types.Report) error {
	fmt.Println()

	if report.TotalViolations == 0 {
		r.printSuccessReport(report)
		return nil
	}

	r.printViolationReport(report)
	return nil
}

func (r *ConsoleReporter) printSuccessReport(report *types.Report) {
	printBox("✓ 检测完成", []string{
		fmt.Sprintf("共检测 %d 个文件，未发现问题", report.TotalFiles),
	}, ColorGreen)
	fmt.Println()
}

func (r *ConsoleReporter) printViolationReport(report *types.Report) {
	severityStats := r.calculateSeverityStats(report)

	var lines []string
	lines = append(lines, fmt.Sprintf("共检测 %d 个文件，发现 %d 个问题", report.TotalFiles, report.TotalViolations))
	lines = append(lines, "")
	lines = append(lines, "【严重程度统计】")

	if severityStats[types.SeverityCritical] > 0 {
		lines = append(lines, fmt.Sprintf("  极其严重: %d", severityStats[types.SeverityCritical]))
	}
	if severityStats[types.SeverityHigh] > 0 {
		lines = append(lines, fmt.Sprintf("  严重:     %d", severityStats[types.SeverityHigh]))
	}
	if severityStats[types.SeverityMedium] > 0 {
		lines = append(lines, fmt.Sprintf("  中等:     %d", severityStats[types.SeverityMedium]))
	}
	if severityStats[types.SeverityLow] > 0 {
		lines = append(lines, fmt.Sprintf("  低:       %d", severityStats[types.SeverityLow]))
	}

	lines = append(lines, "")
	lines = append(lines, "【文件详情】")

	fileStats := r.calculateFileStats(report)
	for filePath, stats := range fileStats {
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("  %s", truncatePath(filePath, 50)))
		lines = append(lines, fmt.Sprintf("    问题数: %d", stats.total))
		if stats.critical > 0 {
			lines = append(lines, fmt.Sprintf("      极其严重: %d", stats.critical))
		}
		if stats.high > 0 {
			lines = append(lines, fmt.Sprintf("      严重:     %d", stats.high))
		}
		if stats.medium > 0 {
			lines = append(lines, fmt.Sprintf("      中等:     %d", stats.medium))
		}
		if stats.low > 0 {
			lines = append(lines, fmt.Sprintf("      低:       %d", stats.low))
		}
	}

	printBox("检测结果", lines, ColorYellow)
	fmt.Println()
}

func printBox(title string, lines []string, titleColor Color) {
	width := 60

	fmt.Printf("%s%s%s\n", ColorCyan, strings.Repeat("─", width), ColorReset)

	centeredTitle := centerText(title, width)
	fmt.Printf("%s%s%s%s%s\n", titleColor, ColorBold, centeredTitle, ColorReset, ColorReset)

	fmt.Printf("%s%s%s\n", ColorCyan, strings.Repeat("─", width), ColorReset)

	for _, line := range lines {
		fmt.Printf("  %s\n", line)
	}

	fmt.Printf("%s%s%s\n", ColorCyan, strings.Repeat("─", width), ColorReset)
}

func centerText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text[:width]
	}
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-textLen-padding)
}

func (r *ConsoleReporter) calculateSeverityStats(report *types.Report) map[types.Severity]int {
	stats := make(map[types.Severity]int)
	for _, violation := range report.Violations {
		stats[violation.Severity]++
	}
	return stats
}

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

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}
