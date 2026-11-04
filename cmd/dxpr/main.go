package main

import (
	"fmt"
	"os"
	"strings"

	dxpr "github.com/shvixxl/dxpr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dxpr <expression>")
		os.Exit(1)
	}

	result, err := dxpr.Evaluate(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Fprintln(os.Stderr, "dxpr:", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
