package log2csv

import (
	"bytes"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	log := Log{
		time.Unix(1366274236, 919928546),
		&Format{Header: "c1,c2,c3"},
		[]string{"1", "2", "3"},
	}

	testdata := []struct {
		log       *Log
		timestamp bool
		expected  string
	}{
		{
			&log,
			false,
			"c1,c2,c3\n1,2,3\n",
		},
		{
			&log,
			true,
			"unixtime,c1,c2,c3\n1366274236.919929,1,2,3\n",
		},
	}

	out := new(bytes.Buffer)
	for _, item := range testdata {
		cw := NewCSVWriter(out, item.timestamp, false)
		if err := cw.Write(item.log); err != nil {
			t.Fatal(err)
		}

		actual := string(out.Bytes())
		if item.expected != actual {
			t.Fatalf("expected %x, got %x", item.expected, actual)
		}
		out.Reset()
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
		if result := fmtFrac(item.time, item.prec); item.expected != result {
			t.Fatalf("expected %s, got %s", item.expected, result)
		}
	}
}
