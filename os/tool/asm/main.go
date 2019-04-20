package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/nanasi880/til/os/tool/asm/assembler"
)

var (
	sourceFileName = flag.String("f", "", "source file name or path (stdin by default)")
	outputFileName = flag.String("o", "", "output file name or path (stdout by default)")
)

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
		defer fclose(f)
	}
	if *outputFileName != "" {
		f, err := os.OpenFile(*outputFileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		outputFile = f
		defer fclose(f)
	}

	if err := assembler.New().Exec(sourceFile, outputFile); err != nil {
		log.Fatal(err)
	}
}

// パラメーターエラーがある場合、コマンドの使用方法を表示してプロセスをExitする
func showUsageIfError() {
	if *sourceFileName == "" {
		flag.Usage()
		os.Exit(1)
	}
}

// io.Closerをクローズし、もしエラーが発生した場合はそれをログする
//
// @param closer --- クローズ対象
func fclose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println(err)
	}
}
