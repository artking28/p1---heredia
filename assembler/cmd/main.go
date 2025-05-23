package main

import (
	"ASM/compiler"
	"ASM/models"
	neander "ASM/neanderExecutor"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const DEBUG = ""

func main() {

	const isDebug = len(DEBUG) > 0
	count := len(os.Args)
	if !isDebug {
		if count < 2 || count > 3 {
			log.Fatal(errors.New("error: inválid arguments"))
		}
	}

	inputFile := DEBUG
	if !isDebug {
		inputFile = os.Args[1]
	}
	if !strings.HasSuffix(inputFile, ".asm") {
		log.Fatal(errors.New("error: this file " + inputFile + " is not a ASM file. Please, rename it before compiling"))
	}

	outputFile := strings.Split(inputFile, ".asm")[0] + ".bin"
	if count == 3 && !isDebug {
		outputFile = os.Args[2]
		if !strings.HasSuffix(outputFile, ".bin") {
			log.Fatal(errors.New("error: this is not a bin output file. Please, choose another name to output file"))
		}
	}
	tokens, err := compiler.Tokenize(inputFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	parser := models.NewParser(inputFile, tokens)
	err = compiler.ParseAll(&parser)
	if err != nil {
		log.Fatal(err.Error())
	}

	content, err := parser.WriteProgram()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.WriteFile(outputFile, content, 0744)
	if err != nil {
		log.Fatal(err.Error())
	}

	//parser.Inspect()
	//InterpreterTest()
}

func InterpreterTest() {
	bytes, err := os.ReadFile("misc/test.bin")
	if err != nil {
		log.Fatalf(err.Error())
	}
	neander.PrintProgram(bytes, false, false, false)

	pr, _ := neander.RunProgram(bytes, false, false)
	fmt.Printf("\n\nResult:\n\tAc = %d, Pc = %d, Z = %v, N = %v\n\n", int8(pr.Ac), pr.Pc, pr.Ac == 0, int8(pr.Ac) < 0)
}
