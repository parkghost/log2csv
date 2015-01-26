package log2csv

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
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
		w := new(bytes.Buffer)
		cw := NewCSVWriter(w, false, true)

		c := NewConverter(in, cw)
		err = c.Convert()
		if err != nil {
			t.Fatalf("convert %s.log failed: %s", item, err)
		}

		expected, err := ioutil.ReadFile(item + ".csv")
		if err != nil {
			t.Fatal("unexpected error on ReadFile:", err)
		}

		if !bytes.Equal(expected, w.Bytes()) {
			t.Errorf("convert %s.log, expected\n[%s], got\n[%s]", item, expected, w.Bytes())
		}
	}
}

var errWriteTest = errors.New("Write Test")

type errorWriter struct{}

func (e *errorWriter) Write(log *Log) error {
	return errWriteTest
}

func TestConvertError(t *testing.T) {
	cw := NewCSVWriter(ioutil.Discard, false, false)
	c := NewConverter(&errorReader{}, cw)
	if err := c.Convert(); err != errReadTest {
		t.Fatalf("expected errReadTest, got %v", err)
	}

	testdata := "gc14(2): 1+1+0 ms 10 -> 5 MB 58439 -> 8912 (573381-564469) objects 184 handoff"
	c = NewConverter(strings.NewReader(testdata), &errorWriter{})
	if err := c.Convert(); err != errWriteTest {
		t.Fatalf("expected errWriteTest, got %v", err)
	}
}

func BenchmarkConvert(b *testing.B) {
	b.ReportAllocs()

	testdata :=
		`gc1(1): 3+0+125+1 us, 0 -> 0 MB, 21 (21-0) objects, 2 goroutines, 15/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields
gc2(1): 0+0+99+0 us, 0 -> 0 MB, 48 (49-1) objects, 3 goroutines, 19/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields
gc3(1): 1+0+102+0 us, 0 -> 0 MB, 178 (197-19) objects, 5 goroutines, 25/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields
gc4(1): 1+0+139+1 us, 0 -> 0 MB, 302 (384-82) objects, 5 goroutines, 33/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields
gc5(2): 120+4+2438+5 us, 0 -> 0 MB, 541 (677-136) objects, 8 goroutines, 50/0/0 sweeps, 0(0) handoff, 1(3) steal, 16/10/2 yields`
	r := strings.NewReader(testdata)
	cw := NewCSVWriter(ioutil.Discard, false, true)

	for i := 0; i < b.N; i++ {
		c := NewConverter(r, cw)
		if err := c.Convert(); err != nil {
			b.Fatalf("convert log to csv failed: %s", err)
		}

		r.Seek(0, 0)
	}
}
