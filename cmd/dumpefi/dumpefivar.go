package main

import (
	"fmt"
	"os"

	"github.com/vilroi/efitools"
)

var searchFlag bool

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
	}

	efivars := efitools.ReadEfiVars()
	efivar, ok := efivars.Search(os.Args[1])
	if !ok {
		err("'%s' cannot be found\n", os.Args[1])
	}

	fmt.Println(efivar.Name)
	fmt.Println(efivar.Data)
}

func usage(prog string) {
	fmt.Fprintf(os.Stderr, "usage: %s efivar\n", prog)
	os.Exit(1)
}

func err(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
