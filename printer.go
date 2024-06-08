package main

import "fmt"

func Printfln(template string, args ...interface{}) {
	fmt.Printf(template+"\n", args...)
}
