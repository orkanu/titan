package main

import (
	"fmt"
	"log"
)

func main() {
	flags, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(flags.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", cfg)
}
