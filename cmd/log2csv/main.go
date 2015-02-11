package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/parkghost/log2csv"
)

var (
	input, output string
	timestamp     bool
)

func init() {
	flag.StringVar(&input, "i", "stdin", "The input file")
	flag.StringVar(&output, "o", "stdout", "The output file")
	flag.BoolVar(&timestamp, "t", true, "Add timestamp at line head (the input file must be `stdin`)")
	flag.Usage = func() {
		fmt.Println("Usage1: log2csv -i gc.log -o gc.csv")
		fmt.Println("Usage2: GODEBUG=gctrace=1 your-go-program 2>&1 | log2csv -o gc.csv")
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(-1)
	}

	r, err := file(input, false)
	checkError(err)
	defer r.Close()

	w, err := file(output, true)
	checkError(err)
	defer w.Close()

	timestamp = timestamp && input == "stdin"
	buffering := input != "stdin"
	cw := log2csv.NewCSVWriter(w, timestamp, buffering)
	sc := log2csv.NewScanner(r, log2csv.GCTraceFormats)
	converter := log2csv.NewConverter(sc, cw)

	err = converter.Convert()
	checkError(err)
}

func file(name string, create bool) (*os.File, error) {
	switch name {
	case "stdin":
		return os.Stdin, nil
	case "stdout":
		return os.Stdout, nil
	default:
		if create {
			return os.Create(name)
		}
		return os.Open(name)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
