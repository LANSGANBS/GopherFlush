package types

// Severity 违规严重程度
type Severity int

const (
	SeverityLow      Severity = iota // 低
	SeverityMedium                   // 中等
	SeverityHigh                     // 严重
	SeverityCritical                 // 极其严重
)

// String 返回严重程度的字符串表示
func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "低"
	case SeverityMedium:
		return "中等"
	case SeverityHigh:
		return "严重"
	case SeverityCritical:
		return "极其严重"
	default:
		return "未知"
	}
}

// Violation 代码违规记录
type Violation struct {
	RuleName    string   // 规则名称
	Severity    Severity // 严重程度
	FilePath    string   // 文件路径
	Line        int      // 行号
	Column      int      // 列号
	Message     string   // 违规描述
	Suggestion  string   // 修复建议
}

// Report 检测报告
type Report struct {
	TotalFiles      int          // 总文件数
	TotalViolations int          // 总违规数
	Violations      []*Violation // 违规列表
}
