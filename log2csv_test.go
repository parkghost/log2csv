package main

import (
	"bytes"
	"io/ioutil"
	"os"
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
		GO11,
	},
	&TestData{
		"gc14(2): 1+1+0 ms 10 -> 5 MB 58439 -> 8912 (573381-564469) objects 184 handoff",
		"14,2,1,1,0,10,5,58439,8912,573381,564469,184",
		GO10,
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
		if output, err := convert(data.text, data.version); err != nil {
			t.Fatal(err)
		} else {
			if output != data.expected {
				t.Fatalf("expected %s, got %s", data.expected, output)
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
	}

	for _, item := range testGcLogs {
		in, err := os.Open(item + ".log")
		if err != nil {
			t.Fatalf("cannot open %s.log", item)
		}

		out := &bytes.Buffer{}

		process(in, out)

		csvData, err := ioutil.ReadFile(item + ".csv")

		if err != nil {
			t.Fatalf("cannot read %s.csv", item)
		}

		actual := string(out.Bytes())
		expected := string(csvData)

		if expected != actual {
			t.Fatalf("expected\n %s, got\n %s", expected, actual)
		}

	}
}

func TestGetReader(t *testing.T) {

	reader1, _ := getReader("testdata/testdata_go_1_0_3.log")
	if _, ok := reader1.(*os.File); !ok {
		t.Fatalf("expected get os.File, got %#+v", reader1)
	}

	if isStdin {
		t.Fatalf("expected isStdin was false, got true")
	}
	defer reader1.(*os.File).Close()

	reader2, _ := getReader("")
	if reader2 != os.Stdin {
		t.Fatalf("expected get os.Stdin, got %#+v", reader2)
	}
	if !isStdin {
		t.Fatalf("expected isStdin was true, got false")
	}

}

func TestGetWriter(t *testing.T) {
	_, err := getWriter("")
	if err == nil {
		t.Fatal("expected err, got nil")
	}
}
