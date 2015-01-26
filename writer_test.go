package log2csv

import (
	"bytes"
	"testing"
	"time"
)

func TestCSVWriter(t *testing.T) {
	example := &Log{
		time.Unix(1366274236, 919928546),
		&Format{Header: "c1,c2,c3"},
		[]string{"1", "2", "3"},
	}

	var testdata = []struct {
		log       *Log
		timestamp bool
		expected  string
	}{
		{
			example,
			false,
			"c1,c2,c3\n1,2,3\n",
		},
		// enable timestamp
		{
			example,
			true,
			"unixtime,c1,c2,c3\n1366274236.919929,1,2,3\n",
		},
	}

	for _, item := range testdata {
		w := new(bytes.Buffer)
		cw := NewCSVWriter(w, item.timestamp, false)
		if err := cw.Write(item.log); err != nil {
			t.Error("unexpected error on Write:", err)
			continue
		}

		actual := string(w.Bytes())
		if item.expected != actual {
			t.Fatalf("expected %s, got %s", item.expected, actual)
		}
	}
}

func TestFmtFrac(t *testing.T) {
	testData := []struct {
		time     time.Time
		prec     int
		expected string
	}{
		{
			time.Unix(1366274236, 919928546),
			6,
			"1366274236.919929",
		},
		{
			time.Unix(1366274236, 919928546),
			3,
			"1366274236.920",
		},
	}

	for _, item := range testData {
		actual := fmtFrac(item.time, item.prec)
		if item.expected != actual {
			t.Fatalf("fmtFrac(%s, %d) => %s, want %s", item.time, item.prec, actual, item.expected)
		}
	}
}
