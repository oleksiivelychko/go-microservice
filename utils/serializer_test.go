package utils

import (
	"bytes"
	"testing"
)

type Test struct {
	Field string
}

func TestUtils_ToJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	err := ToJSON(&Test{Field: "test"}, buf)

	if err != nil {
		t.Error(err)
	}

	if buf.String() != "{\"Field\":\"test\"}\n" {
		t.Errorf("unable to compare JSON string %s", buf.String())
	}
}

func TestUtils_FromJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write([]byte("{\"Field\":\"test\"}\n"))

	test := &Test{}
	err := FromJSON(test, buf)

	if err != nil {
		t.Error(err)
	}

	if test.Field != "test" {
		t.Errorf("unable to compare JSON field %s to string", test.Field)
	}
}
