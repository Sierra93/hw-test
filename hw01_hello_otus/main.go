package main

import (
	"fmt"

	// Пакет для реверса строк.
	"golang.org/x/example/hello/reverse"
)

func main() {
	// Реверсим строку Hello, OTUS!.
	fmt.Println(reverse.String("Hello, OTUS!"))
}
