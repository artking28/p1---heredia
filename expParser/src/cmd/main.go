package main

import (
	"exps-heredia/parser"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 3 {
		panic("error, the input must match the pattern: exp <inputName>.lpn <outputName>.asm")
	}

	if !strings.HasSuffix(os.Args[1], ".lpn") {
		panic("error" + os.Args[1] + "input is not a lpn, rename it please")
	}

	if !strings.HasSuffix(os.Args[2], ".asm") {
		panic("error" + os.Args[2] + "output is not a asm, rename it please")
	}

	//file := "../expression.lpn"
	p, err := parser.NewParser(os.Args[1], os.Args[2])
	if err != nil {
		panic(err.Error())
	}

	//out, err := p.ParseExpression(false)
	out, err := p.ParseScope(parser.RootScope)
	if err != nil {
		panic(err.Error())
	}

	//println(out.String())

	asm, err := out.WriteMemASM()
	if err != nil {
		panic(err.Error())
	}

	err = os.WriteFile(p.OutputFile, []byte(asm), 0766)
	if err != nil {
		panic(err.Error())
	}

	//fmt.Printf("%v\n", asm)
}
