package datetime

import (
	"database/sql/driver"
	"strings"
	"time"
)

/*
TimeJSON represents time in RFC3339 format and makes that appropriate for DB structure.
*/
type TimeJSON struct {
	time.Time
}

/*
MarshalJSON invoked when json.Marshal and returns (current) time in RFC3339.
*/
func (datetime *TimeJSON) MarshalJSON() ([]byte, error) {
	var instantTime time.Time
	if datetime.Time.UnixNano() == instantTime.UnixNano() {
		return []byte("\"" + time.Now().Format(time.RFC3339) + "\""), nil
	}

	return []byte("\"" + datetime.Time.Format(time.RFC3339) + "\""), nil
}

/*
UnmarshalJSON invoked when json.Unmarshal and assigns the given value to struct.
*/
func (datetime *TimeJSON) UnmarshalJSON(b []byte) (err error) {
	parsedTime, err := time.Parse(time.RFC3339, strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}

	datetime.Time = parsedTime
	return
}

/*
Value implements Valuer interface. Transforms struct into value, appropriate for DB.
*/
func (datetime *TimeJSON) Value() (driver.Value, error) {
	var instantTime time.Time
	if datetime.Time.UnixNano() == instantTime.UnixNano() {
		return nil, nil
	}

	return datetime.Time, nil
}

/*
Scan implements sql.Scanner interface. Transforms value from DB into appropriate for struct.
*/
func (datetime *TimeJSON) Scan(v interface{}) error {
	var bytesArr []byte
	for _, b := range v.([]uint8) {
		bytesArr = append(bytesArr, b)
	}

	parsedTime, err := time.Parse(time.DateTime, string(bytesArr))
	if err != nil {
		return err
	}

	*datetime = TimeJSON{parsedTime}
	return nil
}
