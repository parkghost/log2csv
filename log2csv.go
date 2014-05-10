package main

import (
	"bufio"
	"encoding/csv"
	"errors"
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
	GO_1_1_AND_1_2
)

var (
	regexes = map[int]*regexp.Regexp{
		GO_1_0:         regexp.MustCompile(`gc(\d+)\((\d+)\):\s(\d+)\+(\d+)\+(\d+)\s\w+\s(\d+)\s->\s(\d+)\s\w+\s(\d+)\s->\s(\d+)\s\((\d+)-(\d+)\)\sobjects\s(\d+)\shandoff`),
		GO_1_1_AND_1_2: regexp.MustCompile(`gc(\d+)\((\d+)\):\s(\d+)\+(\d+)\+(\d+)\s\w+,\s(\d+)\s->\s(\d+)\s\w+\s(\d+)\s->\s(\d+)\s\((\d+)-(\d+)\)\sobjects,\s(\d+)\((\d+)\)\shandoff,\s(\d+)\((\d+)\)\ssteal,\s(\d+)\/(\d+)\/(\d+)\syields`),
	}
	header = map[int]string{
		GO_1_0:         "numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff",
		GO_1_1_AND_1_2: "numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep",
	}

	errVersionNotFound = errors.New("can't detected version")
)

func newReader(file string) (reader io.Reader, err error) {
	if isTTY {
		reader = os.Stdin
	} else {
		reader, err = os.Open(file)
	}

	return
}

func newWriter(file string) (writer io.Writer, err error) {
	if file == "" {
		err = errors.New("required output file parameter")
	} else {
		writer, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	}

	return
}

func process(reader io.Reader, writer io.Writer) (err error) {
	scanner := bufio.NewScanner(reader)
	csvWriter := csv.NewWriter(writer)

	defer func() {
		if err == nil {
			csvWriter.Flush()
			err = csvWriter.Error()
		}
	}()

	currentLogVersion := Unknown
	for scanner.Scan() {

		filtered := false
		line := scanner.Text()

		// detect the log version if the current log version is Unknown
		if currentLogVersion == Unknown {
			if version, errVersion := detectLogVersion(line); errVersion == nil {
				currentLogVersion = version

				header := header[currentLogVersion]
				if isTTY && timestamp {
					header = "unixtime," + header
				}
				err = csvWriter.Write(strings.Split(header, ","))
				if err != nil {
					return
				}
			}
		}

		if currentLogVersion != Unknown {

			// parse and convert string from raw string to csv structure
			if record, errConvert := convert(line, currentLogVersion); errConvert == nil {

				if isTTY && timestamp {
					record = append([]string{fmtFrac(time.Now(), 6)}, record...)
				}

				err = csvWriter.Write(record)
				if err != nil {
					return
				}

				filtered = true
			}
		}

		if isTTY && filtered == false {
			fmt.Println(line)
		}

	}

	err = scanner.Err()
	if err != nil {
		return
	}

	return
}

func detectLogVersion(line string) (version int, err error) {
	found := false

	// Find the version from string and check all versions of log patterns
	for ver, regexp := range regexes {
		if regexp.MatchString(line) {
			if found == true {
				return Unknown, fmt.Errorf("ambiguous log version: %s", line)
			}
			found = true
			version = ver
		}
	}

	if !found {
		err = errVersionNotFound
	}

	return
}

func convert(input string, version int) (output []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("unmatched uint string: %s => %s", input, e)
		}
	}()

	if matched := regexes[version].FindStringSubmatch(input); matched == nil {
		err = fmt.Errorf("unmatched string: %s", input)
	} else {
		output = matched[1:]
	}

	return
}

func fmtFrac(t time.Time, prec int) string {
	unixNano := t.UnixNano()
	fmtStr := "%." + strconv.Itoa(prec) + "f"

	return fmt.Sprintf(fmtStr, float64(unixNano)/10e8)
}
