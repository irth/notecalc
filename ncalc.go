package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/kr/pretty"
)

func main() {
	// read command line arguments
	input_filename := flag.String("input", "/dev/stdin", "Input filename")
	flag.Parse()

	code, _ := ioutil.ReadFile(*input_filename)
	for _, exp := range Parse(code) {
		fmt.Printf("%# v\n", pretty.Formatter(Evaluate(exp, make(map[string]Value))))
	}
}
