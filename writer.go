package log2csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Writer interface {
	Write(*Log) error
}

type Flusher interface {
	Flush() error
}

type csvWriter struct {
	w         *csv.Writer
	timestamp bool
	bufferred bool

	wroteHeader bool
}

func (cw *csvWriter) Write(log *Log) error {
	if !cw.wroteHeader {
		if err := cw.writeHeader(log); err != nil {
			return err
		}
		cw.wroteHeader = true
	}

	if err := cw.writeLog(log); err != nil {
		return err
	}
	if !cw.bufferred {
		cw.w.Flush()
	}

	return cw.w.Error()
}

func (cw *csvWriter) writeHeader(log *Log) error {
	header := log.Format.Header
	if cw.timestamp {
		header = "unixtime," + header
	}

	return cw.w.Write(strings.Split(header, ","))
}

func (cw *csvWriter) writeLog(log *Log) error {
	fields := log.Fields
	if cw.timestamp {
		fields = append([]string{fmtFrac(log.Timestamp, 6)}, fields...)
	}

	return cw.w.Write(fields)
}

func (cw *csvWriter) Flush() error {
	cw.w.Flush()

	return cw.w.Error()
}

func fmtFrac(t time.Time, prec int) string {
	unixNano := t.UnixNano()
	fmtStr := "%." + strconv.Itoa(prec) + "f"

	return fmt.Sprintf(fmtStr, float64(unixNano)/10e8)
}

func NewCSVWriter(w io.Writer, timestamp bool, bufferred bool) *csvWriter {
	cw := new(csvWriter)
	cw.w = csv.NewWriter(w)
	cw.timestamp = timestamp
	cw.bufferred = bufferred

	return cw
}
