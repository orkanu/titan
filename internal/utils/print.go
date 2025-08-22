package utils

// https://gist.github.com/vratiu/9780109

import "fmt"

func PrintlnBlack(a ...any) {
	fmt.Printf("\x1b[30;1m%v\x1b[0m\n", a...)
}

func PrintlnRed(a ...any) {
	fmt.Printf("\x1b[31;1m%v\x1b[0m\n", a...)
}

func PrintlnGreen(a ...any) {
	fmt.Printf("\x1b[32;1m%v\x1b[0m\n", a...)
}

func PrintlnYellow(a ...any) {
	fmt.Printf("\x1b[33;1m%v\x1b[0m\n", a...)
}

func PrintlnBlue(a ...any) {
	fmt.Printf("\x1b[34;1m%v\x1b[0m\n", a...)
}

func PrintlnPurple(a ...any) {
	fmt.Printf("\x1b[35;1m%v\x1b[0m\n", a...)
}

func PrintlnCyan(a ...any) {
	fmt.Printf("\x1b[36;1m%v\x1b[0m\n", a...)
}

func PrintlnWhite(a ...any) {
	fmt.Printf("\x1b[37;1m%v\x1b[0m\n", a...)
}
