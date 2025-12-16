package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type NullableDateOnly struct {
	Time  *time.Time
	Valid bool
}

// ========= JSON → Struct =========
func (d *NullableDateOnly) UnmarshalJSON(b []byte) error {
	// handle "null"
	if string(b) == "null" {
		d.Time = nil
		d.Valid = false
		return nil
	}

	// handle ""
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		d.Time = nil
		d.Valid = false
		return nil
	}

	// parse date
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	d.Time = &t
	d.Valid = true
	return nil
}

// ========= Struct → JSON =========
func (d NullableDateOnly) MarshalJSON() ([]byte, error) {
	if !d.Valid || d.Time == nil {
		return []byte(`null`), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// ========= DB → Struct (Scan) =========
func (d *NullableDateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = nil
		d.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = &v
		d.Valid = true
		return nil

	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = &t
		d.Valid = true
		return nil

	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = &t
		d.Valid = true
		return nil
	}

	return fmt.Errorf("cannot scan type %T into NullableDateOnly", value)
}

// ========= Struct → DB (Value) =========
func (d NullableDateOnly) Value() (driver.Value, error) {
	if !d.Valid || d.Time == nil {
		return nil, nil // Insert NULL into DB
	}
	return d.Time.Format("2006-01-02"), nil
}
