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
	Unknown = iota
	GO_1_0
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

	errVersionNotFound = errors.New("can't detected version")

	inputFile, outputFile string
	timestamp             bool
	isStdin               = false
)

func init() {
	flag.StringVar(&inputFile, "i", "", "The input file (default: Stdin)")
	flag.StringVar(&outputFile, "o", "", "The output file")
	flag.BoolVar(&timestamp, "t", false, "Add timestamp at line head(Stdin input only)")

	flag.Usage = func() {
		fmt.Println("Usage1: log2csv -i gc.log -o gc.csv")
		fmt.Println("Usage2: GOGCTRACE=1 your-go-program 2>&1 | log2csv -o gc.csv")
		flag.PrintDefaults()
		fmt.Println("  -h   : show help usage")
	}
}

func detectLogVersion(line string) (version int, err error) {
	for version, regexp := range regexes {
		if regexp.MatchString(line) {
			return version, nil
		}
	}

	err = errVersionNotFound
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

func process(reader io.Reader, writer io.Writer) {
	bufReader := bufio.NewReader(reader)
	bufWriter := bufio.NewWriter(writer)

	currentLogVersion := Unknown
	for {

		filtered := false
		if line, err := bufReader.ReadString('\n'); err != nil {
			break
		} else {
			if currentLogVersion == Unknown {
				if version, err := detectLogVersion(line); err == nil {
					currentLogVersion = version
					err := writeHeader(bufWriter, currentLogVersion)
					checkError(err)
				}

			}

			if currentLogVersion != Unknown {
				if output, err := convert(line, currentLogVersion); err == nil {
					err := writeBody(bufWriter, output)
					checkError(err)
					filtered = true
				}
			}

			if isStdin && filtered == false {
				fmt.Print(line)
			}
		}

	}
	err := bufWriter.Flush()
	checkError(err)
}

func writeHeader(writer *bufio.Writer, version int) (err error) {
	prefix := ""
	if isStdin && timestamp {
		prefix = "unixtime,"
	}
	_, err = writer.WriteString(prefix + header[version] + "\n")
	return
}

func writeBody(writer *bufio.Writer, output string) (err error) {
	prefix := ""
	if isStdin && timestamp {
		prefix = fmtFrac(time.Now(), 6) + ","
	}

	_, err = writer.WriteString(prefix + output + "\n")
	return
}

func fmtFrac(t time.Time, prec int) string {
	unixNano := t.UnixNano()
	fmtStr := "%." + strconv.Itoa(prec) + "f"
	return fmt.Sprintf(fmtStr, float64(unixNano)/10e8)
}

func getReader(file string) (reader io.Reader, err error) {
	if file == "" {
		reader = os.Stdin
		isStdin = true
	} else {
		reader, err = os.Open(file)
	}
	return
}

func getWriter(file string) (writer io.Writer, err error) {
	if file == "" {
		err = errors.New("required output file parameter")
	} else {
		writer, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	}
	return
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
}

func main() {
	flag.Parse()

	reader, err := getReader(inputFile)
	checkError(err)
	defer reader.(*os.File).Close()

	writer, err := getWriter(outputFile)
	checkError(err)
	defer writer.(*os.File).Close()

	process(reader, writer)
}
