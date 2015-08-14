package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	// read command line arguments
	input_filename := flag.String("input", "/dev/stdin", "Input filename")
	flag.Parse()

	code, _ := ioutil.ReadFile(*input_filename)
	fmt.Println(Tokenize(code))
}
