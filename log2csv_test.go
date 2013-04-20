package main

import (
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
		"gc13(2): 48+24+2 ms, 263 -> 124 MB 1891444 -> 938285 (6426929-5488644) objects, 1(8) handoff, 3(11016) steal, 21/2/0 yields",
		"13,2,48,24,2,263,124,1891444,938285,6426929,5488644,1,8,3,11016,21,2,0",
		GO_1_1,
	},
	&TestData{
		"gc14(2): 1+1+0 ms 10 -> 5 MB 58439 -> 8912 (573381-564469) objects 184 handoff",
		"14,2,1,1,0,10,5,58439,8912,573381,564469,184",
		GO_1_0,
	},
}

func TestDetectLogVersion(t *testing.T) {
	for _, data := range testDatas {
		version := detectLogVersion(data.text)
		if version != data.version {
			t.Fatalf("expected %d, get %d :%s", data.version, version, data.text)
		}
	}
}

func TestConvert(t *testing.T) {
	for _, data := range testDatas {
		if output, err := convert(data.text, data.version); err != nil {
			t.Fatal(err)
		} else {
			if output != data.expected {
				t.Fatalf("expected %s, get %s", data.expected, output)
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
			t.Fatalf("expected %s, get %s", data.expected, output)
		}
	}
}
