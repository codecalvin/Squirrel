package main

import "fmt"

type KK interface {
	K(string)
}

type H func(int)

func (h H) K(c string) {
	fmt.Println(c)
}

func main() {
	print("hello world")
	e := func (int) {}
	H(e).K("s")
}
