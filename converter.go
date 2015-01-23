package log2csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Converter struct {
	r         io.Reader
	cw        *csv.Writer
	timestamp bool
}

func (c *Converter) Run() error {
	sc := NewScanner(c.r, formats)
	wroteHeader := false

	for {
		log := sc.Scan()
		if log == nil {
			break
		}

		if !wroteHeader {
			if err := c.writeHeader(log); err != nil {
				return err
			}
			wroteHeader = true
		}

		if err := c.writeLog(log); err != nil {
			return err
		}
		c.cw.Flush()
	}

	if sc.Err() != nil {
		return sc.Err()
	}
	return c.cw.Error()
}

func (c *Converter) writeHeader(log *Log) error {
	header := log.Format.Header
	if c.timestamp {
		header = "unixtime," + header
	}

	return c.cw.Write(strings.Split(header, ","))
}

func (c *Converter) writeLog(log *Log) error {
	fields := log.Fields
	if c.timestamp {
		fields = append([]string{fmtFrac(log.Timestamp, 6)}, fields...)
	}

	return c.cw.Write(fields)
}

func fmtFrac(t time.Time, prec int) string {
	unixNano := t.UnixNano()
	fmtStr := "%." + strconv.Itoa(prec) + "f"

	return fmt.Sprintf(fmtStr, float64(unixNano)/10e8)
}

func NewConverter(r io.Reader, w io.Writer, timestamp bool) *Converter {
	c := new(Converter)
	c.r = r
	c.cw = csv.NewWriter(w)
	c.timestamp = timestamp

	return c
}
