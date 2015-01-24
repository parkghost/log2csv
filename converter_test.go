package log2csv

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
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
		cw := NewCSVWriter(out, false, true)

		c := NewConverter(in, cw)
		err = c.Convert()
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
