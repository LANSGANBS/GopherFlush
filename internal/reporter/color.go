package reporter

import (
	"runtime"
)

// Color ANSI 颜色代码
type Color string

const (
	ColorReset   Color = "\033[0m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorWhite   Color = "\033[37m"
	ColorOrange  Color = "\033[38;5;208m" // 橙色（256色）
)

// Colorize 给文本添加颜色
func Colorize(text string, color Color) string {
	// Windows 终端可能不支持颜色，需要特殊处理
	if runtime.GOOS == "windows" {
		// Windows 10+ 支持 ANSI 颜色，但为了兼容性，可以选择禁用
		// 这里我们仍然尝试使用颜色
	}
	return string(color) + text + string(ColorReset)
}

// GetSeverityColor 根据严重程度获取颜色
func GetSeverityColor(severity string) Color {
	switch severity {
	case "极其严重":
		return ColorRed
	case "严重":
		return ColorOrange
	case "中等":
		return ColorYellow
	case "低":
		return ColorGreen
	default:
		return ColorWhite
	}
}
