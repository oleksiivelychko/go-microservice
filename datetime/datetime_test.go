package datetime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const sampleRFC3339 = "2002-10-02T10:00:00-05:00"

type TestDateTimeJSON struct {
	DateTime JSON `json:"datetime"`
}

func TestDateTime_MarshalNowTimeJSON(t *testing.T) {
	toMarshal := &TestDateTimeJSON{DateTime: JSON{}}
	marshaledJSON, err := json.Marshal(toMarshal)
	if err != nil {
		t.Fatal(err)
	}

	unmarshalTo := &TestDateTimeJSON{}
	err = json.Unmarshal(marshaledJSON, &unmarshalTo)
	if err != nil {
		t.Fatal(err)
	}

	datetime := unmarshalTo.DateTime.Format(time.RFC3339)
	datetimeNow := time.Now().Format(time.RFC3339)
	if datetime != datetimeNow {
		t.Errorf("time mismatch: %s != %s", datetime, datetimeNow)
	}
}

func TestDateTime_MarshalTimeJSON(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339, sampleRFC3339)
	if err != nil {
		t.Fatal(err)
	}

	toMarshal := &TestDateTimeJSON{DateTime: JSON{parsedTime}}
	marshaledJSON, err := json.Marshal(toMarshal)
	if err != nil {
		t.Fatal(err)
	}

	unmarshalTo := &TestDateTimeJSON{}
	err = json.Unmarshal(marshaledJSON, &unmarshalTo)
	if err != nil {
		t.Fatal(err)
	}

	datetime := unmarshalTo.DateTime.Format(time.RFC3339)
	if datetime != sampleRFC3339 {
		t.Errorf("time mismatch: %s != %s", datetime, sampleRFC3339)
	}
}

func TestDateTime_UnmarshalTimeJSON(t *testing.T) {
	unmarshalTo := &TestDateTimeJSON{}
	stringJSON := []byte(fmt.Sprintf(`{"datetime":"%s"}`, sampleRFC3339))

	err := json.Unmarshal(stringJSON, &unmarshalTo)
	if err != nil {
		t.Fatal(err)
	}

	parsedTime, err := time.Parse(time.RFC3339, sampleRFC3339)
	if err != nil {
		t.Fatal(err)
	}

	datetimeJSON := &TestDateTimeJSON{DateTime: JSON{parsedTime}}
	datetime := datetimeJSON.DateTime.Format(time.RFC3339)
	if datetime != sampleRFC3339 {
		t.Errorf("time mismatch: %s != %s", datetime, sampleRFC3339)
	}
}

func TestDateTime_EncodeTimeJSON(t *testing.T) {
	buf := new(bytes.Buffer)

	parsedTime, err := time.Parse(time.RFC3339, sampleRFC3339)
	if err != nil {
		t.Fatal(err)
	}

	toMarshal := &TestDateTimeJSON{DateTime: JSON{parsedTime}}
	err = json.NewEncoder(buf).Encode(toMarshal)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != fmt.Sprintf("{\"datetime\":\"%s\"}\n", sampleRFC3339) {
		t.Errorf("JSON string %s and structure field are not equal", buf.String())
	}
}

func TestDateTime_DecodeTimeJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(fmt.Sprintf("{\"datetime\":\"%s\"}\n", sampleRFC3339)))

	unmarshalTo := &TestDateTimeJSON{}
	err := json.NewDecoder(buf).Decode(unmarshalTo)
	if err != nil {
		t.Fatal(err)
	}

	datetime := unmarshalTo.DateTime.Format(time.RFC3339)
	if datetime != sampleRFC3339 {
		t.Errorf("structure field %s and JSON string %s are not equal", datetime, buf.String())
	}
}
