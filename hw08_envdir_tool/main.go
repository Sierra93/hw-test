package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("invalid arguments number: want - 2, got - %d\n", len(os.Args))
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
