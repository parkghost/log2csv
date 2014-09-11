package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	inputFile, outputFile string
	timestamp             bool

	isTTY = true
)

func main() {
	flag.StringVar(&inputFile, "i", "", "The input file (default: Stdin)")
	flag.StringVar(&outputFile, "o", "", "The output file")
	flag.BoolVar(&timestamp, "t", false, "Add timestamp at line head(Stdin input only)")

	flag.Usage = func() {
		fmt.Println("Usage1: log2csv -i gc.log -o gc.csv")
		fmt.Println("Usage2: GODEBUG=gctrace=1 your-go-program 2>&1 | log2csv -o gc.csv\n" +
			"       (GO version below 1.2) GOGCTRACE=1 your-go-program 2>&1 | log2csv -o gc.csv")
		flag.PrintDefaults()
		fmt.Println("  -h   : show help usage")
	}
	flag.Parse()

	if flag.NArg() != 0 || flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	isTTY = inputFile == ""

	reader, err := newReader(inputFile)
	checkError(err)
	defer reader.Close()

	writer, err := newWriter(outputFile)
	checkError(err)
	defer writer.Close()

	err = process(reader, writer)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
}
