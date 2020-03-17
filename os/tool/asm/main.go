package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nanasi880/til/os/tool/asm/assembler"
)

var (
	sourceFileName string
	outputFileName string
)

func init() {
	flag.StringVar(&sourceFileName, "f", "", "source file name or path (stdin by default)")
	flag.StringVar(&outputFileName, "o", "", "output file name or path (stdout by default)")
}

func main() {
	os.Exit(_main())
}

func _main() int {
	flag.Parse()

	var (
		sourceFile = os.Stdin
		outputFile = os.Stdout
	)
	if sourceFileName != "" {
		f, err := os.Open(sourceFileName)
		if err != nil {
			errorln(err)
			return 1
		}
		sourceFile = f
		defer fclose(f)
	}
	if outputFileName != "" {
		f, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			errorln(err)
			return 1
		}
		outputFile = f
		defer fclose(f)
	}

	if err := assembler.New().Exec(sourceFile, outputFile); err != nil {
		errorln(err)
		return 1
	}

	return 0
}

func errorln(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}

// io.Closerをクローズし、もしエラーが発生した場合はそれをログする
//
// @param closer --- クローズ対象
func fclose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println(err)
	}
}
