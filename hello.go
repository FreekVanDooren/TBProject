package main

import (
	"fmt"
	"strings"
	"tbp.com/user/hello/hi"
)

type Greetable interface {
	Hello() string
}

func main() {
	fmt.Println(Hi())
	var tom Greetable = hi.Tom()
	hello := tom.Hello()
	fmt.Println(hello)
	hellos := []rune(hello)
	fmt.Printf("%s%s",strings.ToLower(string(hellos[0])), string(hellos[1:len(hellos)-1]))
}

func Hi() string {
	return "Hello"
}
