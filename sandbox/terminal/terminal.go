package terminal

import (
	"os"

	term "golang.org/x/crypto/ssh/terminal"
)

const (
	reset   = "\u001b[0m"
	black   = "\u001b[30m"
	red     = "\u001b[31m"
	green   = "\u001b[32m"
	yellow  = "\u001b[33m"
	blue    = "\u001b[34m"
	magenta = "\u001b[35m"
	cyan    = "\u001b[36m"
	white   = "\u001b[37m"
)

var isatty = term.IsTerminal(int(os.Stdout.Fd()))

func color(c string, s string) string {
	if isatty {
		return c + s + reset
	}
	return s
}

func Black(s string) string {
	return color(black, s)
}

func Red(s string) string {
	return color(red, s)
}

func Green(s string) string {
	return color(green, s)
}

func Yellow(s string) string {
	return color(yellow, s)
}

func Blue(s string) string {
	return color(blue, s)
}

func Magenta(s string) string {
	return color(magenta, s)
}

func Cyan(s string) string {
	return color(cyan, s)
}

func White(s string) string {
	return color(white, s)
}
