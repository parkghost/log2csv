package log2csv

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestConverter(t *testing.T) {
	testdata := []string{
		"testdata/testdata_go_1_0_3",
		"testdata/testdata_go_1_1",
		"testdata/testdata_go_1_2",
		"testdata/testdata_go_1_3",
		"testdata/testdata_go_1_4",
	}

	for _, item := range testdata {
		in, err := os.Open(item + ".log")
		if err != nil {
			t.Fatal(err)
		}
		out := new(bytes.Buffer)

		c := NewConverter(in, out, false)
		err = c.Run()
		if err != nil {
			t.Fatalf("convert %s.log to csv format failed: %s", item, err)
		}

		expected, err := ioutil.ReadFile(item + ".csv")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, out.Bytes()) {
			t.Fatalf("expected\n[%s], got\n[%s]", expected, out.Bytes())
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
		if result := fmtFrac(item.time, item.prec); item.expected != result {
			t.Fatalf("expected %s, got %s", item.expected, result)
		}
	}
}
