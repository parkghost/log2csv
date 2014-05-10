package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TestData struct {
	text     string
	expected string
	version  int
}

var testDatas = []*TestData{
	&TestData{
		"gc14(2): 1+1+0 ms 10 -> 5 MB 58439 -> 8912 (573381-564469) objects 184 handoff",
		"14,2,1,1,0,10,5,58439,8912,573381,564469,184",
		GO_1_0,
	},
	&TestData{ // 1.1
		"gc13(2): 48+24+2 ms, 263 -> 124 MB 1891444 -> 938285 (6426929-5488644) objects, 1(8) handoff, 3(11016) steal, 21/2/0 yields",
		"13,2,48,24,2,263,124,1891444,938285,6426929,5488644,1,8,3,11016,21,2,0",
		GO_1_1_AND_1_2,
	},
	&TestData{ // 1.2
		"gc63(2): 3+1+0 ms, 15 -> 7 MB 167805 -> 12894 (9983900-9971006) objects, 0(0) handoff, 4(350) steal, 16/2/0 yields",
		"63,2,3,1,0,15,7,167805,12894,9983900,9971006,0,0,4,350,16,2,0",
		GO_1_1_AND_1_2,
	},
}

func TestDetectLogVersion(t *testing.T) {
	for _, data := range testDatas {
		version, err := detectLogVersion(data.text)
		if err != nil || version != data.version {
			t.Fatalf("expected %d, got %d :%s", data.version, version, data.text)
		}
	}
}

func TestConvert(t *testing.T) {
	for _, data := range testDatas {
		if record, err := convert(data.text, data.version); err != nil {
			t.Fatal(err)
		} else {

			expected := strings.Split(data.expected, ",")
			if !reflect.DeepEqual(expected, record) {
				t.Fatalf("expected %s, got %s", expected, record)
			}
		}
	}
}

type TestTimeData struct {
	time     time.Time
	prec     int
	expected string
}

func TestFmtFrac(t *testing.T) {
	testTimeDatas := []*TestTimeData{
		&TestTimeData{
			time.Unix(1366274236, 919928546),
			6,
			"1366274236.919929",
		},
		&TestTimeData{
			time.Unix(1366274236, 919928546),
			3,
			"1366274236.920",
		},
	}

	for _, data := range testTimeDatas {
		if output := fmtFrac(data.time, data.prec); output != data.expected {
			t.Fatalf("expected %s, got %s", data.expected, output)
		}
	}
}

func TestConvertGcLog(t *testing.T) {
	testGcLogs := []string{
		"testdata/testdata_go_1_0_3",
		"testdata/testdata_go_1_1",
		"testdata/testdata_go_1_2",
	}

	for _, item := range testGcLogs {
		in, err := os.Open(item + ".log")
		if err != nil {
			t.Fatalf("cannot open %s.log", item)
		}

		out := &bytes.Buffer{}

		err = process(in, out)
		if err != nil {
			t.Fatalf("cannot process %s.csv: %s", item, err)
		}

		csvData, err := ioutil.ReadFile(item + ".csv")
		if err != nil {
			t.Fatalf("cannot read %s.csv: %s", item, err)
		}

		actual := string(out.Bytes())
		expected := string(csvData)

		if expected != actual {
			t.Fatalf("expected\n %s, got\n %s", expected, actual)
		}

	}
}
