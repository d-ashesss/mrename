package main

import (
	"fmt"
	"os"
)

func main() {
	p := Processor{}
	result := p.Process(os.Args[1:])
	fmt.Println(result)
}
