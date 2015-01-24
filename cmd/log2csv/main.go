package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/parkghost/log2csv"
)

var (
	inputFile  = flag.String("i", "", "The input file (default: standard input)")
	outputFile = flag.String("o", "", "The output file")
	timestamp  = flag.Bool("t", true, "Add timestamp at line head (the input file must be standard input)")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage1: log2csv -i gc.log -o gc.csv")
		fmt.Println("Usage2: GODEBUG=gctrace=1 your-go-program 2>&1 | log2csv -o gc.csv")
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() != 0 || flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	r, err := newReader(*inputFile)
	checkError(err)
	defer r.Close()

	w, err := newWriter(*outputFile)
	checkError(err)
	defer w.Close()

	cw := log2csv.NewCSVWriter(w, *timestamp && isTTY(), !isTTY())
	converter := log2csv.NewConverter(r, cw)
	checkError(converter.Convert())

}

func newReader(file string) (io.ReadCloser, error) {
	if isTTY() {
		return os.Stdin, nil
	}

	return os.Open(file)
}

func newWriter(file string) (io.WriteCloser, error) {
	if file == "" {
		return nil, errors.New("required output file parameter")
	}

	return os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func isTTY() bool {
	return *inputFile == ""
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
