package parser

import (
	"fmt"
	"strconv"
)

func compare(left, right interface{}, op string) bool {
	switch l := left.(type) {
	case string:
		rs := fmt.Sprintf("%v", right)
		switch op {
		case "==":
			return l == rs
		case "!=":
			return l != rs
		case ">":
			return l > rs
		case "<":
			return l < rs
		case ">=":
			return l >= rs
		case "<=":
			return l <= rs
		}
	case float64:
		rf, _ := strconv.ParseFloat(fmt.Sprintf("%v", right), 64)
		switch op {
		case "==":
			return l == rf
		case "!=":
			return l != rf
		case ">":
			return l > rf
		case "<":
			return l < rf
		case ">=":
			return l >= rf
		case "<=":
			return l <= rf
		}
	}
	return false
}
