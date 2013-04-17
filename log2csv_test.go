package main

import (
	"testing"
)

type TestData struct {
	text     string
	expected string
}

func TestConvert(t *testing.T) {
	testData := &TestData{
		`gc13(2): 48+24+2 ms, 263 -> 124 MB 1891444 -> 938285 (6426929-5488644) objects, 1(8) handoff, 3(11016) steal, 21/2/0 yields`,
		"13,2,48,24,2,263,124,1891444,938285,6426929,5488644,1,8,3,11016,21,2,0",
	}

	if output, err := convert(testData.text); err != nil {
		t.Fatal(err)
	} else {
		if output != testData.expected {
			t.Fatalf("expected %s, get %s", testData.expected, output)
		}
	}

}
