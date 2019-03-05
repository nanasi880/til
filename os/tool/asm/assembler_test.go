package main

import (
	"io/ioutil"
	"testing"
)

func TestAssembler_DB(t *testing.T) {

	asmFile, err := ioutil.ReadFile("testdata/db_only.asm.txt")
	if err != nil {
		t.Fatal(err)
	}

	a := new(assembler)
	b, err := a.asm(asmFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 14 {
		t.Fatal(b)
	}
}
