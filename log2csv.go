package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const LOG_PATTERN = `gc(\d+)\((\d+)\):\s(\d+)\+(\d+)\+(\d+)\s\w+,\s(\d+)\s->\s(\d+)\s\w+\s+(\d+)\s->\s(\d+)\s\((\d+)-(\d+)\)\sobjects,\s(\d+)\((\d+)\)\shandoff,\s(\d+)\((\d+)\)\ssteal,\s(\d+)\/(\d+)\/(\d+)\syields`

var (
	logRegex   = regexp.MustCompile(LOG_PATTERN)
	inputFile  = flag.String("i", "", "The input file (default: Stdin)")
	outputFile = flag.String("o", "", "The output file (default: Stdout)")
	timestamp  = flag.Bool("t", false, "Add timestamp at line head(Stdin input only)")
	help       = flag.Bool("h", false, "Show Usage")
	isStdin    = false
)

func convert(input string) (output string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("unmatched uint string: %s => %s", input, e))
		}
	}()

	if matched := logRegex.FindStringSubmatch(input); matched == nil {
		err = errors.New(fmt.Sprintf("unmatched string: %s", input))
	} else {
		output = strings.Join(matched[1:], ",")
	}
	return
}

func run(in, out *os.File) {

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)

	prefix := ""
	if *timestamp {
		prefix = "starttime,"
	}

	writer.WriteString(prefix + "numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep\n")
	for {
		if line, err := reader.ReadString('\n'); err != nil {
			break
		} else {
			if output, err := convert(line); err == nil {
				prefix := ""

				if *timestamp {
					prefix = strconv.FormatInt(time.Now().Unix(), 10) + ","
				}

				writer.WriteString(prefix + output + "\n")
				if isStdin {
					writer.Flush()
				}
			}
		}

	}
	writer.Flush()

}

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Println("Usage1: log2csv -i gc.log -o gc.csv")
		fmt.Println("Usage2: GCTRACE=1 your-go-program 2>&1 | log2csv -o gc.csv")
		flag.PrintDefaults()
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var in, out *os.File

	if *inputFile != "" {
		var err error
		in, err = os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open input file: %s", err)
			os.Exit(1)
		}
		defer in.Close()
	} else {
		in = os.Stdin
		isStdin = true
	}

	if *outputFile != "" {
		var err error
		out, err = os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create output file: %s", err)
			os.Exit(1)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	run(in, out)
}
