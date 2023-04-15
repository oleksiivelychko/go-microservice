package data

import (
	"encoding/json"
	"os"
	"testing"
)

func TestData_ReadFile(t *testing.T) {
	data := map[string]int{
		"id": 1,
	}

	file, _ := json.MarshalIndent(data, "", " ")
	err := os.WriteFile("test.json", file, 0644)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ReadFile("test.json")
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) == 0 {
		t.Fatal("unable to read data from file")
	}

	err = os.Remove("test.json")
	if err != nil {
		t.Fatal(err)
	}
}
