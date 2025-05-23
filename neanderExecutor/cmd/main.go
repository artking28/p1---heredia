package main

import (
    "fmt"
    "neander"
    "os"
    "strings"
)

func main() {

    if len(os.Args) != 2 {
        panic("error, the input must match the pattern: exp <inputName>.lpn <outputName>.asm")
    }

    if !strings.HasSuffix(os.Args[1], ".bin") {
        panic("error" + os.Args[1] + "input is not a bin, rename it please")
    }

    bytes, err := os.ReadFile(os.Args[1])
    if err != nil {
        panic(err.Error())
    }

    neander.PrintProgram(bytes, false, false, false)
    pr, _ := neander.RunProgram(bytes, false, false)
    fmt.Printf("\n\nResult:\n\tAc = %d, Pc = %d, Z = %v, N = %v\n\n", pr.Ac, pr.Pc, pr.Ac == 0, int8(pr.Ac) < 0)
}
