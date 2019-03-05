package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	sourceFileName = flag.String("f", "", "source file name or path")
)

func showUsageIfError() {
	if *sourceFileName == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	showUsageIfError()

	sourceFile, err := ioutil.ReadFile(*sourceFileName)
	if err != nil {
		log.Fatal(err)
	}

	asm := new(assembler)
	if _, err := asm.asm(sourceFile); err != nil {
		log.Fatal(err)
	}
}