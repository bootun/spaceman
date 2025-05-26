package text

import "fmt"

// ANSIColor 表示 ANSI 颜色枚举
type ANSIColor int

const (
	// 前景色
	Black ANSIColor = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	// 背景色从 40 开始
	BgBlack ANSIColor = iota + 40 - 8
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Colorize 返回带有颜色的字符串
func Colorize(text string, fg ANSIColor, bg ANSIColor) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", fg, bg, text)
}
