package fancy

import (
	"fmt"
)

type Style int

const (
	None Style = 0
)

const (
	Black Style = 30 + iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func Print(s Style, a ...interface{}) {
	changeStyle(s)
	fmt.Print(a...)
	changeStyle(None)
}

func Println(s Style, a ...interface{}) {
	changeStyle(s)
	fmt.Println(a...)
	changeStyle(None)
}

func Printf(s Style, format string, a ...interface{}) {
	changeStyle(s)
	fmt.Printf(format, a...)
	changeStyle(None)
}

func changeStyle(s Style) {
	fmt.Printf("\x1b[%dm", int(s))
}
