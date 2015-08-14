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
	fmt.Printf("%# v\n", pretty.Formatter(Parse(code)))
}
