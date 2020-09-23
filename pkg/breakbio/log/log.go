package log

import (
	"fmt"
	"io"
)

const (
	green     = `[1;32m`
	red       = `[1;31m`
	cyan      = `[1;36m`
	nocolor   = ""
	stopColor = `[1;0m`

	msgPrefix = "[BreakBio] "
)

func Greenf(msg string, args ...interface{}) {
	Green(fmt.Sprintf(msg, args...))
}

func Green(msg string) {
	fmt.Println(formatMessage(msg, green))
}

func Cyanf(msg string, args ...interface{}) {
	Cyan(fmt.Sprintf(msg, args...))
}

func Cyan(msg string) {
	fmt.Println(formatMessage(msg, cyan))
}

func Redf(msg string, args ...interface{}) {
	Red(fmt.Sprintf(msg, args...))
}

func Red(msg string) {
	fmt.Println(formatMessage(msg, red))
}

func NoColorf(msg string, args ...interface{}) {
	NoColor(fmt.Sprintf(msg, args...))
}

func NoColor(msg string) {
	fmt.Println(formatMessage(msg, nocolor))
}

func formatMessage(msg string, color string) string {
	formattedMsg := color + msgPrefix + msg + stopColor
	return formattedMsg
}

func TeeGreen(writer io.Writer, msg string) {
	tee(writer, msg, green)
}

func TeeRed(writer io.Writer, msg string) {
	tee(writer, msg, red)
}

func TeeNoColor(writer io.Writer, msg string) {
	tee(writer, msg, nocolor)
}

func tee(writer io.Writer, msg string, color string) {
	formattedMsg := formatMessage(msg, color)
	_, _ = fmt.Fprintln(writer, formattedMsg)
	fmt.Println(formattedMsg)
}
