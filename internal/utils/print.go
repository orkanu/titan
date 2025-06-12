package utils

import "fmt"

func PrintlnRed(a ...any) {
	fmt.Printf("\x1b[31;1m%v\x1b[0m\n", a...)
}
