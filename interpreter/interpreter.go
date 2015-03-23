package main

import (
	"flag"
	"fmt"
	"github.com/TalkTakesTime/brainfuck"
	"io/ioutil"
)

var (
	Program = ",[.,]"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		if len(args) > 1 {
			// error, show help message
			return
		}

		filename := args[0]
		program, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("File %s could not be used\n", filename)
			return
		}
		err = brainfuck.Run(string(program), false)
		if err != nil {
			fmt.Printf("File %s does not contain a valid Brainfuck program: %s\n",
				filename, err)
			return
		}

		return
	}

	//brainfuck.Run(Program, true)
}
