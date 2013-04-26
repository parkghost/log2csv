package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	GO_1_0 = iota
	GO_1_1
)

var (
	regexes = map[int]*regexp.Regexp{
		GO_1_0: regexp.MustCompile(`gc(\d+)\((\d+)\):\s(\d+)\+(\d+)\+(\d+)\s\w+\s(\d+)\s->\s(\d+)\s\w+\s(\d+)\s->\s(\d+)\s\((\d+)-(\d+)\)\sobjects\s(\d+)\shandoff`),
		GO_1_1: regexp.MustCompile(`gc(\d+)\((\d+)\):\s(\d+)\+(\d+)\+(\d+)\s\w+,\s(\d+)\s->\s(\d+)\s\w+\s(\d+)\s->\s(\d+)\s\((\d+)-(\d+)\)\sobjects,\s(\d+)\((\d+)\)\shandoff,\s(\d+)\((\d+)\)\ssteal,\s(\d+)\/(\d+)\/(\d+)\syields`),
	}
	header = map[int]string{
		GO_1_0: "numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff",
		GO_1_1: "numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep",
	}

	versionNotFound = errors.New("can't detected version")

	inputFile  = flag.String("i", "", "The input file (default: Stdin)")
	outputFile = flag.String("o", "", "The output file")
	timestamp  = flag.Bool("t", false, "Add timestamp at line head(Stdin input only)")
	isStdin    = false
)

func detectLogVersion(line string) (version int, err error) {
	for version, regexp := range regexes {
		if regexp.MatchString(line) {
			return version, nil
		}
	}

	err = versionNotFound
	return
}

func convert(input string, version int) (output string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("unmatched uint string: %s => %s", input, e))
		}
	}()

	if matched := regexes[version].FindStringSubmatch(input); matched == nil {
		err = errors.New(fmt.Sprintf("unmatched string: %s", input))
	} else {
		output = strings.Join(matched[1:], ",")
	}
	return
}

func run(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)

	currentLogVersion := -1
	for {

		filtered := false
		if line, err := reader.ReadString('\n'); err != nil {
			break
		} else {
			if currentLogVersion == -1 {
				if version, err := detectLogVersion(line); err == nil {
					currentLogVersion = version
					writeHeader(writer, currentLogVersion)
				}
			}

			if currentLogVersion != -1 {
				if output, err := convert(line, currentLogVersion); err == nil {
					writeBody(writer, output)
					filtered = true
				}
			}

			if isStdin && filtered == false {
				fmt.Print(line)
			}
		}

	}
	writer.Flush()
}

func writeHeader(writer *bufio.Writer, version int) {
	prefix := ""
	if isStdin && *timestamp {
		prefix = "unixtime,"
	}
	writer.WriteString(prefix + header[version] + "\n")
}

func writeBody(writer *bufio.Writer, output string) {
	prefix := ""
	if isStdin && *timestamp {
		prefix = fmtFrac(time.Now(), 6) + ","
	}

	writer.WriteString(prefix + output + "\n")
}

func fmtFrac(t time.Time, prec int) string {
	unixNano := t.UnixNano()
	fmtStr := "%." + strconv.Itoa(prec) + "f"
	return fmt.Sprintf(fmtStr, float64(unixNano)/10e8)
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage1: log2csv -i gc.log -o gc.csv")
		fmt.Println("Usage2: GOGCTRACE=1 your-go-program 2>&1 | log2csv -o gc.csv")
		flag.PrintDefaults()
		fmt.Println("  -h   : show help usage")
	}

	flag.Parse()

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
		fmt.Fprint(os.Stderr, "required output file parameter\n")
		flag.Usage()
		os.Exit(-1)
	}

	run(in, out)
}
