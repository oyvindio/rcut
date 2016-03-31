package main

import (
	"reflect"
	"testing"
)

func TestParseFieldArgsHappyDay(t *testing.T) {
	actual, err := parseFieldArgs("1,2,3,42,1337")
	if err != nil {
		t.Errorf("parseFieldArgs returned an error: %v", err)
	}
	expected := []uint64{1, 2, 3, 42, 1337}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("%v did not equal %v", actual, expected)
	}
}

func TestParseFieldArgsSingle(t *testing.T) {
	actual, err := parseFieldArgs("1")
	if err != nil {
		t.Errorf("parseFieldArgs returned an error: %v", err)
	}
	expected := []uint64{1}
	if reflect.DeepEqual(actual, expected) {
		t.Errorf("%v did not equal %v", actual, expected)
	}
}

func TestParseFieldArgsNotAnInt(t *testing.T) {
	bogus := "bogus args"
	_, err := parseFieldArgs(bogus)
	if err == nil {
		t.Errorf("parseFieldArgs didn't return an error with parameter: %v", bogus)
	}
}

func TestParseFieldArgsEmpty(t *testing.T) {
	bogus := ""
	_, err := parseFieldArgs(bogus)
	if err == nil {
		t.Error("parseFieldArgs didn't return an error with empty parameter")
	}
}
