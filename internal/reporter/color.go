package reporter

import (
	"runtime"
)

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
	ColorOrange  Color = "\033[38;5;208m"
	ColorBold    Color = "\033[1m"
	ColorDim     Color = "\033[90m"
)

func ApplyColor(text string, color Color) string {
	if runtime.GOOS == "windows" {
	}
	return string(color) + text + string(ColorReset)
}

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
