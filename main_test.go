package main

import (
	"testing"
)

func TestParseFieldArgsHappyDay(t *testing.T) {
	actual, err := parseFieldArgs("1,2-,3-4,42,1337-1340")
	if err != nil {
		t.Errorf("parseFieldArgs returned an error: %v", err)
	}
	for _, field := range []uint{1, 2, 3, 4, 42, 1337, 1338, 1339, 1340} {
		if !actual[field] {
			t.Errorf("%v was not in the range as expected", field)
		}
	}
}

// func TestParseFieldArgsSingle(t *testing.T) {
// 	actual, err := parseFieldArgs("1")
// 	if err != nil {
// 		t.Errorf("parseFieldArgs returned an error: %v", err)
// 	}
// 	expected := []uint64{1}
// 	if reflect.DeepEqual(actual, expected) {
// 		t.Errorf("%v did not equal %v", actual, expected)
// 	}
// }

// func TestParseFieldArgsNotAnInt(t *testing.T) {
// 	bogus := "bogus args"
// 	_, err := parseFieldArgs(bogus)
// 	if err == nil {
// 		t.Errorf("parseFieldArgs didn't return an error with parameter: %v", bogus)
// 	}
// }

// func TestParseFieldArgsEmpty(t *testing.T) {
// 	bogus := ""
// 	_, err := parseFieldArgs(bogus)
// 	if err == nil {
// 		t.Error("parseFieldArgs didn't return an error with empty parameter")
// 	}
// }
