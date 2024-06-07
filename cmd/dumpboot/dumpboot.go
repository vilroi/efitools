package main

import (
	"fmt"
	"os"

	"github.com/vilroi/efitools"
)

func main() {
	efivars := efitools.GetEfiVars()

	bootvars := extractBootVars(efivars)
	for i, bootvar := range bootvars {
		fmt.Printf("%s: ", bootvar.Name)
		loadopt := efitools.ParseLoadOption(bootvar.Data)
		fmt.Printf("%d: %s\n", i, loadopt.Desc)
		//fmt.Printf("%d: %s\n", i, string(loadopt.OptData))
	}
}

func extractBootVars(efivars efitools.EfiVars) efitools.EfiVars {
	bootorder, ok := efivars.Search("BootOrder")
	if !ok {
		err("failed to get BootOrder")
	}

	var bootvars efitools.EfiVars
	var ordernum int
	for i := 0; i < len(bootorder.Data); i += 2 {
		val := uint16(bootorder.Data[i]) | uint16(bootorder.Data[i+1]<<8)
		varname := fmt.Sprintf("Boot00%02X", val)

		bootvar, ok := efivars.Search(varname)
		if !ok {
			fmt.Fprintf(os.Stderr, "Boot variable %s not found\n", varname)
			continue
		}
		bootvars = append(bootvars, bootvar)

		// Debug
		fmt.Printf("%d. %s\n", ordernum, varname)
		ordernum++
	}

	return bootvars

}

func err(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
