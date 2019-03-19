package main

import (
	"flag"
	"log"
	"os"
)

var (
	sourceFileName = flag.String("f", "", "source file name or path (stdin by default)")
	outputFileName = flag.String("o", "", "output file name or path (stdout by default)")
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

	var (
		sourceFile = os.Stdin
		outputFile = os.Stdout
	)
	if *sourceFileName != "" {
		f, err := os.Open(*sourceFileName)
		if err != nil {
			log.Fatal(err)
		}
		sourceFile = f
		defer f.Close()
	}
	if *outputFileName != "" {
		f, err := os.OpenFile(*outputFileName, os.O_RDWR | os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		outputFile = f
		defer f.Close()
	}

	asm := new(assembler)
	if err := asm.asm(sourceFile, outputFile); err != nil {
		log.Fatal(err)
	}
}
