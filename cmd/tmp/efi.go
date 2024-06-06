package test

import (
	"fmt"
	"os"
	"strings"
)

const efivar_path string = "/sys/firmware/efi/efivars/"

type EfiVar struct {
	Name     string
	FullName string
	Data     []byte
}

type EfiVars []EfiVar

func (e EfiVars) Search(name string) (EfiVar, bool) {
	for _, efivar := range e {
		if efivar.Name == name {
			return efivar, true
		}
	}

	return EfiVar{}, false
}

func ReadEfiVars() EfiVars {
	dirents, err := os.ReadDir(efivar_path)
	check(err)

	var efivars []EfiVar
	for _, dirent := range dirents {
		var efivar EfiVar

		efivar.FullName = dirent.Name()
		efivar.Name = strings.Split(efivar.FullName, "-")[0]
		efivar.Data = readEfiVar(efivar.FullName)

		efivars = append(efivars, efivar)
	}

	return efivars
}

func readEfiVar(v string) []byte {
	fullpath := efivar_path + v
	data, err := os.ReadFile(fullpath)
	check(err)

	return data
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func err(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
