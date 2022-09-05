package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
	"bufio"
)

// Types
const (
	NUMBER		int	= 1
	TEXT			= 2
	CHARACTER		= 3
	CONDITION		= 4
	TYPE			= 5
	NONE			= 6
)

// The variable structure
type Variable struct {
	Type		int

	VarNumber	int
	VarText		string
	VarChar		byte
	VarType		int
	VarCondition	bool
}

// Global variables
var ProgramCounter int
var Variables map[string]Variable

// parseString parses the code string, separates it into lines and words
func parseString(text string) [][]string {
	// Separate the program into lines
	lines := strings.Split(text, "\n")

	// Separate lines into actual code and comments
	for i, l := range lines {
		split := strings.Split(l, ";")
		lines[i] = split[0] // ub?
	}

	code := make([][]string, len(lines))

	// Separate the program into words
	for l := 0; l < len(lines); l++ {
		code[l] = strings.Fields(lines[l])
	}

	return code
}

// This function creates a map of the program, which keys are the row numbers 
func parseCode(code [][]string) (rows map[int][]string, max int) {
	rows = make(map[int][]string, len(code))
	max = 0

	for i := 0; i < len(code) - 1; i++ {
		r, err := strconv.ParseInt(code[i][0], 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		rows[int(r)] = code[i][1:]
		if int(r) > max {
			max = int(r)
		}
	}

	return
}

// This function reads a string/TEXT during execution
// It may require executing a command, reading a variable value
// or just reading a string ("This is a string!")
// Returns a boolean value based on if it actually worked, the actual string and
// a list of arguments that weren't used
func readString(args []string) (bool, string, []string) {
	var s string

	// Not a "traditional" string if it doesn't start with quotation marks
	if args[0][0] != '"' {
		_, ok := Variables[args[0]]
		if !ok {
			// As a final backup, we read it as a command
			N, remaining := runCommand(args[0], args[1:])

			return true, N.VarText, remaining
		}

		return true, Variables[args[0]].VarText, args[1:]
	}

	// If this is true -> it's a one word string
	if args[0][len(args[0]) - 1] == '"' {
		s = args[0][1:(len(args[0]) - 1)]

		return true, s, args[1:]
	}

	s = args[0][1:]
	s += " "

	for i, arg := range args[1:] {
		// If it is a multi-word string we loop throught it till we find the end (")
		if arg[len(arg) - 1] == '"' {
			s += arg[:(len(arg) - 1)]

			return true, s, args[(i+2):]
		}

		s += arg
		s += " "
	}

	return false, s, args
}

func readType(str string) int {
	switch(str) {
	case "NUMBER":
		return NUMBER
	case "TEXT":
		return TEXT
	case "TYPE":
		return TYPE
	case "CHARACTER":
		return CHARACTER
	case "CONDITION":
		return CONDITION
	default:
		return NONE
	}
}

func readNumber(args []string) (bool, int, []string) {
	n, err := strconv.ParseInt(args[0], 10, 32)

	// No errors while parsing -> we got a number straight up
	if err == nil {
		return true, int(n), args[1:]
	}

	// Otherwise, it may be a variable
	_, ok := Variables[args[0]]
	if ok {
		return true, Variables[args[0]].VarNumber, args[1:]
	}

	// As a final backup, we read it as a command
	N, remaining := runCommand(args[0], args[1:])

	return true, N.VarNumber, remaining
}

func readCondition(args []string) (bool, bool, []string) {
	// Sometimes a condition is simply one argument: true or false
	if args[0] == "TRUE" {
		return true, true, args[1:]
	} else if args[0] == "FALSE" {
		return true, false, args[1:]
	}

	// Sometimes it's a variable
	_, ok := Variables[args[0]]
	if ok {
		return true, Variables[args[0]].VarCondition, args[1:]
	}

	// Most often it'll be a command's return value 
	N, remaining := runCommand(args[0], args[1:])

	return true, N.VarCondition, remaining
}

func commandPrint(command string, args []string) (remaining []string) {
	var temp bool
	var str string
	temp, str, remaining = readString(args)
	if temp != true {
		log.Fatal("Expected string!")
	}

	switch((command)) {
	case "SAY":
		fmt.Println(str)
		break
	case "SHOUT":
		fmt.Println(strings.ToUpper(str) + "!")
		break
	case "WHISPER":
		fmt.Println(strings.ToLower(str))
		break
	case "PRINT":
		fmt.Printf(str)
		break
	}

	return remaining
}

func commandMath(command string, args []string) (r Variable, remaining []string) {
	r.Type = NUMBER

	var n1, n2 int
	var ok1, ok2 bool

	ok1, n1, remaining = readNumber(args)
	ok2, n2, remaining = readNumber(remaining)

	if !ok1 || !ok2 {
		log.Fatal("Failed to convert to number")
	}

	switch(command) {
	case "SUM":
		r.VarNumber = (n1 + n2)
	case "DIVIDE":
		r.VarNumber = (n1 / n2)
	case "MULTIPLY":
		r.VarNumber = (n1 * n2)
	case "MODULO":
		r.VarNumber = (n1 % n2)
	}

	return
}

func commandAsk(t string, args []string) (r Variable, remaining []string) {
	var T int
	var str string
	var answer string
	var err error

	T = readType(t)

	r.Type = T

	_, str, remaining = readString(args)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(str)
	answer, err = reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	answer = answer[0:(len(answer)-1)]

	switch(T) {
	case TEXT:
		r.VarText = answer
	case NUMBER:
		temp, err := strconv.ParseInt(answer, 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		r.VarNumber = int(temp)
	default:
		log.Fatal("Can't use type: ", t)
	}

	return
}

func runCommand(command string, args []string) (r Variable, remaining []string) {
	r.Type = NONE

	// Parse through the command + arguments
	// Only exception is the "END" command, which is checked in runProgram
	switch(strings.ToUpper(command)) {
	case "SAY":
		remaining = commandPrint(strings.ToUpper(command), args)
	case "SHOUT":
		remaining = commandPrint(strings.ToUpper(command), args)
	case "WHISPER":
		remaining = commandPrint(strings.ToUpper(command), args)
	case "PRINT":
		remaining = commandPrint(strings.ToUpper(command), args)
	case "PRINTNUM":
		var num int
		_, num, remaining = readNumber(args)
		fmt.Printf("%v", num)
	case "JUMP":
		_, ProgramCounter, remaining = readNumber(args)
	case "RESERVE":
		v := Variable{
			Type:		readType(args[2]),
			VarNumber:	0,
			VarText:	"",
			VarChar:	0,
			VarType:	NONE,
		}

		Variables[args[0]] = v

		remaining = args[3:]

	case "INCREMENT":
		Variables[args[0]] = Variable{
			Type:		NUMBER,
			VarNumber:	Variables[args[0]].VarNumber+1,
		}

		remaining = args[1:]

	case "DECREMENT":
		Variables[args[0]] = Variable{
			Type:		NUMBER,
			VarNumber:	Variables[args[0]].VarNumber-1,
		}

		remaining = args[1:]

	case "SUM":
		r, remaining = commandMath(command, args)
	case "DIVIDE":
		r, remaining = commandMath(command, args)
	case "MULTIPLY":
		r, remaining = commandMath(command, args)
	case "MODULO":
		r, remaining = commandMath(command, args)

	case "IF":
		var cond bool
		_, cond, remaining = readCondition(args)

		if remaining[0] != "JUMP" {
			log.Fatal("IF doesn't follow a jump")
		}

		if cond {
			_, ProgramCounter, remaining = readNumber(remaining[1:])
		} else {
			remaining = nil
		}

	case "EQUALS":
		r.Type = CONDITION
		// TODO: Make equals work with any variable type
		var val1, val2 int
		_, val1, remaining = readNumber(args)
		_, val2, remaining = readNumber(remaining)

		r.VarCondition = (val1 == val2)

	case "NOT":
		r.Type = CONDITION
		var cond bool
		_, cond, remaining = readCondition(args)

		// Turn the condition to the opposite
		// TRUE 	-> FALSE
		// FALSE 	-> TRUE
		r.VarCondition = !cond

	case "ASKFOR":
		r, remaining = commandAsk(args[0], args[1:])
	case "ASK":
		r, remaining = commandAsk("TEXT", args)

	default:
		// In this case, the only option is a "variable IS value" type of statement so the command is the variable name
		v, ok := Variables[strings.ToUpper(command)]
		if !ok && args[0] != "IS"{
			log.Fatal("unrecognized command: ", command)
		}

		switch(v.Type) {
		case NUMBER:
			ok, v.VarNumber, remaining = readNumber(args[1:])
		case TEXT:
			ok, v.VarText, remaining = readString(args[1:])
		case CONDITION:
			ok, v.VarCondition, remaining = readCondition(args[1:])
		case TYPE:
			v.VarType = readType(args[1])
			ok = true
			remaining = args[2:]
		default:
			ok = false
		}

		if !ok {
			log.Fatal("unrecognized value")
		}

		Variables[strings.ToUpper(command)] = v
	}

	return r, remaining
}

func runProgram(rows map[int][]string, max int) {
	ProgramCounter = 1
	Variables = make(map[string]Variable, 16)

	for {
		row, ok := rows[ProgramCounter]
		ProgramCounter += 1

		// If the ok check fails, we just go to the next one
		if !ok {
			if ProgramCounter >= max {
				return
			}

			continue
		}

		// Check the 'END' command here because it exits here
		if row[0] == "END" {
			return
		}

		// First item of the row slice is the command (at least this is assumed) and after that
		// the command's arguments typically follow
		_, remaining := runCommand(row[0], row[1:])

		if len(remaining) != 0 {
			log.Fatal("Command doesn't use all arguments")
		}
	}
}

func main() {
	var err error
	var content string
	var code [][]string
	var rows map[int][]string
	var max int

	// If there's only one given argument, we let the user write to it directly to standard input
	if len(os.Args) == 1 {
		content = ""
		var l string

		reader := bufio.NewReader(os.Stdin)

		run := true

		for run {
			fmt.Printf(">>> ")
			l, err = reader.ReadString('\n')
			if err != nil {
				log.Fatal(err, " inputted line: ", l)
			}

			switch(l[:(len(l)-1)]) {
			// In this case we must execute the program (runProgram)
			case "RUN":
				code = parseString(content)

				rows, max = parseCode(code)

				runProgram(rows, max)
			case "EXIT":
				os.Exit(0)
			case "RUNANDEXIT":
				run = false
			default:
				content += l
			}

			// Otherwise we add the inputted line
		}
	} else {
		var fileContent []byte
		// Start reading the program file
		fileContent, err = ioutil.ReadFile(os.Args[1])

		if err != nil {
			log.Fatal(err)
		}

		content = string(fileContent)
	}

	code = parseString(string(content))

	rows, max = parseCode(code)

	runProgram(rows, max)
}
