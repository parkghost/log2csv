package log2csv

import (
	"errors"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	var testdata = []struct {
		formatName string
		text       string
		expected   string
	}{
		{
			"Go 1.0",
			"gc14(2): 1+1+0 ms 10 -> 5 MB 58439 -> 8912 (573381-564469) objects 184 handoff",
			"14,2,1,1,0,10,5,58439,8912,573381,564469,184",
		},
		{
			"Go 1.1",
			"gc13(2): 48+24+2 ms, 263 -> 124 MB 1891444 -> 938285 (6426929-5488644) objects, 1(8) handoff, 3(11016) steal, 21/2/0 yields",
			"13,2,48,24,2,263,124,1891444,938285,6426929,5488644,1,8,3,11016,21,2,0",
		},
		{
			"Go 1.1", // 1.2
			"gc63(2): 3+1+0 ms, 15 -> 7 MB 167805 -> 12894 (9983900-9971006) objects, 0(0) handoff, 4(350) steal, 16/2/0 yields",
			"63,2,3,1,0,15,7,167805,12894,9983900,9971006,0,0,4,350,16,2,0",
		},
		{
			"Go 1.3",
			"gc1(1): 5+0+186+0 us, 0 -> 0 MB, 18 (19-1) objects, 0/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields",
			"1,1,5,0,186,0,0,0,18,19,1,0,0,0,0,0,0,0,0,0,0",
		},
		{
			"Go 1.4",
			"gc1(1): 4+0+1097+3 us, 0 -> 0 MB, 21 (21-0) objects, 2 goroutines, 15/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields",
			"1,1,4,0,1097,3,0,0,21,21,0,2,15,0,0,0,0,0,0,0,0,0",
		},
		{
			"Go 1.4", // add non-gctrace data
			`gc1(1): 3+0+125+1 us, 0 -> 0 MB, 21 (21-0) objects, 2 goroutines, 15/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields
somthing else ...
gc2(1): 0+0+99+0 us, 0 -> 0 MB, 48 (49-1) objects, 3 goroutines, 19/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields`,
			`1,1,3,0,125,1,0,0,21,21,0,2,15,0,0,0,0,0,0,0,0,0
2,1,0,0,99,0,0,0,48,49,1,3,19,0,0,0,0,0,0,0,0,0`,
		},
	}

	for _, item := range testdata {
		sc := NewScanner(strings.NewReader(item.text), formats)

		var logs []*Log
		for {
			log := sc.Scan()
			if log == nil {
				break
			}

			logs = append(logs, log)
		}
		if sc.Err() != nil {
			t.Fatal("unexpected error after Scan:", sc.Err())
		}

		lines := strings.Split(item.expected, "\n")
		for _, log := range logs {
			if item.formatName != log.Format.Name {
				t.Errorf("expected %s, got %s", item.formatName, log.Format.Name)
			}

			expected := lines[0]
			lines = lines[1:]
			actual := strings.Join(log.Fields, ",")
			if expected != actual {
				t.Errorf("scan %s format, expected\n%s, got\n%s", item.formatName, expected, actual)
			}
		}
	}
}

var errReadTest = errors.New("Read Test")

type errorReader struct{}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, errReadTest
}

func TestScanError(t *testing.T) {
	sc := NewScanner(&errorReader{}, formats)
	if log := sc.Scan(); log != nil {
		t.Fatal("expected nil from Scan")
	}

	if sc.Err() != errReadTest {
		t.Fatalf("expected errReadTest, got %v", sc.Err())
	}
}
