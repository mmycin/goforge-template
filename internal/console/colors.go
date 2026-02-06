package console

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func sprintColor(color, s string) string {
	return color + s + colorReset
}

func Success(format string, a ...interface{}) {
	fmt.Printf(sprintColor(colorGreen, "✓ ")+format+"\n", a...)
}

func Error(format string, a ...interface{}) {
	fmt.Printf(sprintColor(colorRed, "✗ ")+format+"\n", a...)
}

func Info(format string, a ...interface{}) {
	fmt.Printf(sprintColor(colorCyan, "ℹ ")+format+"\n", a...)
}

func Warning(format string, a ...interface{}) {
	fmt.Printf(sprintColor(colorYellow, "⚠ ")+format+"\n", a...)
}

func Color(colorCode, s string) string {
	return colorCode + s + colorReset
}
