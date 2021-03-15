package main

import (
	"os"

	"github.com/santiagocezar/cssmod"
)

func main() {
	i, err := os.Open("test/test.module.css")
	if err != nil {
		panic(err)
	}
	defer i.Close()
	o, err := os.Create("test/test_output.css")
	if err != nil {
		panic(err)
	}
	defer o.Close()

	b, _ := cssmod.Transform(i, "test.module.css")

	o.Write(b)
}
