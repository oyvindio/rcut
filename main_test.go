package main

import (
	"testing"
)

func TestParseFieldArgsHappyDay(t *testing.T) {
	param := "1,2-,3-4,42,1337-1340"
	oc, err := createOutputConfig(param, false, " ", false)
	if err != nil {
		t.Errorf("createOutputConfig(%s) returned an error: %v\n", param, err)
	}
	for _, field := range []uint{1, 2, 3, 4, 42, 1337, 1338, 1339, 1340} {
		if !oc.ShouldOutputField(field) {
			t.Errorf("-f %s did not include field %v\n", param, field)
		}
	}
}

func TestParseFieldArgsSingle(t *testing.T) {
	param := "1"
	oc, err := createOutputConfig(param, false, " ", false)
	if err != nil {
		t.Errorf("createOutputConfig(%s) returned an error: %v\n", param, err)
	}
	if !oc.ShouldOutputField(1) {
		t.Errorf("-f %s did not include field 1\n", param)
	}
}

func TestParseFieldArgsStart(t *testing.T) {
	param := "1-"
	oc, err := createOutputConfig(param, false, " ", false)
	if err != nil {
		t.Errorf("createOutputConfig returned an error: %v\n", err)
	}
	for i := uint(1); i <= 10; i++ {
		if !oc.ShouldOutputField(i) {
			t.Errorf("-f %s did not include field %d", param, i)
		}
	}
}

func TestParseFieldArgsEnd(t *testing.T) {
	param := "-3"
	oc, err := createOutputConfig(param, false, " ", false)
	if err != nil {
		t.Errorf("createOutputConfig(%s) returned an error: %v\n", param, err)
	}

	for i := uint(1); i <= 3; i++ {
		if !oc.ShouldOutputField(i) {
			t.Errorf("-f %s did not include field %d\n", param, i)
		}
	}
	for i := uint(4); i <= 10; i++ {
		if oc.ShouldOutputField(i) {
			t.Errorf("-f %s included field %d\n", param, i)
		}
	}
}

func TestParseFieldArgsStartEnd(t *testing.T) {
	param := "3-5"
	oc, err := createOutputConfig(param, false, " ", false)
	if err != nil {
		t.Errorf("createOutputConfig(%s) returned an error: %v\n", param, err)
	}
	for i := uint(1); i <= 2; i++ {
		if oc.ShouldOutputField(i) {
			t.Errorf("-f %s included field %d\n", param, i)
		}
	}
	for i := uint(3); i <= 5; i++ {
		if !oc.ShouldOutputField(i) {
			t.Errorf("-f %s did not include field %d\n", param, i)
		}
	}
	for i := uint(6); i <= 10; i++ {
		if oc.ShouldOutputField(i) {
			t.Errorf("-f %s included field %d\n", param, i)
		}
	}
}

func TestParseFieldArgsNotAnInt(t *testing.T) {
	bogus := "bogus args"
	_, err := createOutputConfig(bogus, false, " ", false)
	if err == nil {
		t.Errorf("createOutputConfig didn't return an error with parameter: %v\n", bogus)
	}
}

func TestParseFieldArgsEmpty(t *testing.T) {
	bogus := ""
	_, err := createOutputConfig(bogus, false, " ", false)
	if err == nil {
		t.Error("createOutputConfig didn't return an error with empty parameter\n")
	}
}
