package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello")
	os.Exit(1) // want "использование прямого вызова os.Exit в функции main пакета main"
}
