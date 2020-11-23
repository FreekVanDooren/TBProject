package hi

import "fmt"

type greeter struct {
	name string
}

func (g greeter) Hello() string {
	return fmt.Sprintf("Hi %s!", g.name)
}

func Tom() greeter {
	return greeter{name: "Tom"}
}
