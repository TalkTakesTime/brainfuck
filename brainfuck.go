/*
Brainfuck interpreter written in Go

Brainfuck is an esoteric language with the following instructions:
  >  Move the pointer to the right
  <  Move the pointer to the left
  +  Increment the memory cell under the pointer
  -  Decrement the memory cell under the pointer
  .  Output the character signified by the cell at the pointer
  ,  Input a character and store it in the cell at the pointer
  [  Jump past the matching ] if the cell under the pointer is 0
  ]  Jump back to the matching [ if the cell under the pointer is nonzero

TODO: change from a global tape to a single tape for each program

For more information on Brainfuck, see http://esolangs.org/wiki/Brainfuck
*/
package brainfuck

import (
	"errors"
	"fmt"
	"github.com/TalkTakesTime/stack"
	"regexp"
	"strconv"
)

// TapeLength is the default length for a standard Brainfuck tape, taken from
// the original implementation by Urban Müller.
const TapeLength = 30000

var (
	// the tape code executes on
	tape = make([]byte, TapeLength)
	// the position of the pointer on the tape
	pointer = 0
	// the position of the current token being executed
	index = 0
	// a stack used to store the positions of [ in code
	loopStack stack.Stack
)

// SpecialInstructionRegex is used to match the special runtime instructions
// given by "!! instruction"
var SpecialInstructionRegex = regexp.MustCompile("^!! ([a-z]+)([0-9]+)?")

// MoveLeft represents the Brainfuck instruction "<". It moves the pointer
// left by one cell, wrapping to the end of the tape if necessary
func MoveLeft() {
	if pointer == 0 {
		pointer = TapeLength - 1
	} else {
		pointer--
	}
}

// MoveRight represents the Brainfuck instruction ">". It moves the pointer
// right by one cell, wrapping to 0 if necessary
func MoveRight() {
	if pointer == TapeLength-1 {
		pointer = 0
	} else {
		pointer++
	}
}

// Increment represents the Brainfuck instruction "+". It increments the memory
// cell under the pointer
func Increment() {
	tape[pointer]++
}

// Decrement represents the Brainfuck instruction "-". It decrements the memory
// cell under the pointer
func Decrement() {
	tape[pointer]--
}

// Output represents the Brainfuck instruction ".". It prints the character
// value of the cell under the pointer
func Output() {
	fmt.Printf("%c", tape[pointer])
}

// Input represents the Brainfuck instruction ",". It takes a character from
// stdin and stores it in the cell under the pointer
func Input() {
	fmt.Scanf("%c", &tape[pointer])
}

// OpenLoop represents the Brainfuck instruction "[". It forms the opening
// part of a loop. If the cell under the pointer is 0, returns true to
// indicate to skip to the next ]
func OpenLoop(pos int) bool {
	if tape[pointer] == 0 {
		return true
	}
	loopStack.Push(pos)
	return false
}

// CloseLoop represents the Brainfuck instruction "]".
// It closes a loop and returns the index of the matching open brace in the
// code, if the cell under the pointer is not 0. Otherwise returns -1
func CloseLoop() int {
	p, err := loopStack.Pop()
	if err != nil {
		panic(err)
	}

	if tape[pointer] == 0 {
		return -1
	}
	return p.(int)
}

// RunSpecialInstruction executes a non-standard runtime instruction for
// various utilities not included in a standard Brainfuck interpreter, such as
// "!! clear", which clears the tape so that multiple Brainfuck programs
// can be run from a single file.
//
// All special instructions are of the form "!! instruction", and the valid
// instructions are as follows:
//   - clear: clears the tape and resets the pointer to position 0
//   - print: prints the contents of the 11 cells surrounding the pointer
//   - printn: prints the contents of the n cells surrounding the pointer
func RunSpecialInstruction(inst []string) {
	switch inst[0] {
	case "clear":
		ClearTape()
	case "print":
		fmt.Println(FormatCells(pointer-5, pointer+5))
	case "printn":
		n, err := strconv.Atoi(inst[1])
		if err != nil {
			return
		}
		fmt.Println(FormatCells(pointer-n/2, pointer+n/2))
	}
}

// FormatCells formats the cells from indices start to end (inclusive)
// into an easily human-readable format and returns the result as a string.
func FormatCells(start, end int) string {
	if start < 0 {
		start = TapeLength + start - 1
	}
	if end >= TapeLength {
		end = end - TapeLength
	}

	indicesText := ""
	cellsText := "[\t"
	for i := start; i != end+1; i++ {
		if i >= TapeLength {
			i = 0
		}
		indicesText += fmt.Sprintf("\t%d", i)
		cellsText += fmt.Sprintf("%d\t", tape[i])
	}
	return indicesText + "\n" + cellsText + "]"
}

// Validate tests to see if the given string of code contains any syntax
// errors -- namely, unmatched closing or opening braces.
func Validate(code string) error {
	var testStack stack.Stack
	for i, r := range code {
		c := string(r)
		if c == "[" {
			testStack.Push(i)
		} else if c == "]" {
			_, err := testStack.Pop()
			if err != nil {
				return errors.New("Syntax error: closing brace without " +
					"matched opening brace")
			}
		}
	}
	if testStack.Length() != 0 {
		return errors.New("Syntax error: opening brace without matching " +
			"closing brace")
	}

	return nil
}

// Run validates and runs the given Brainfuck program, clearing the tape
// before running if clearTape is true. Returns an error if the program
// is invalid, otherwise returns nil.
func Run(code string, clearTape bool) error {
	err := Validate(code)
	if err != nil {
		return err
	}

	if clearTape {
		ClearTape()
	}

	index = 0
	codeLen := len(code)

	for {
		switch string(code[index]) {
		case "<":
			MoveLeft()
		case ">":
			MoveRight()
		case "-":
			Decrement()
		case "+":
			Increment()
		case ".":
			Output()
		case ",":
			Input()
		case "[":
			var skipStack stack.Stack
			skipStack.Push(1)
			for skip := OpenLoop(index); skip; {
				index += 1
				if string(code[index]) == "[" {
					skipStack.Push(1)
				} else if string(code[index]) == "]" {
					skipStack.Pop()
					if skipStack.IsEmpty() {
						skip = false
					}
				}
			}
		case "]":
			if i := CloseLoop(); i != -1 {
				index = i - 1
			}
		// special runtime instructions can be prefixed with "!! " and a
		// human readable command
		case "!":
			match := SpecialInstructionRegex.FindStringSubmatch(code[index:])
			if len(match) > 0 {
				RunSpecialInstruction(match[1:])
			}
		}

		index += 1
		if index >= codeLen {
			break
		}
	}

	// add a new line to ensure nice ending
	fmt.Println()

	return nil
}

// ClearTape clears the tape and resets the pointer position to 0.
func ClearTape() {
	for i := range tape {
		tape[i] = 0
	}

	pointer = 0
}
