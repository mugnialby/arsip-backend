package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type DateOnly struct {
	time.Time
}

// JSON → Struct
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove quotes

	// Accept empty
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

// Struct → JSON
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(d.Format("2006-01-02"))
}

// ===== Required for GORM =====

// Scan (DB → Struct)
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil

	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil

	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	}

	return fmt.Errorf("cannot scan type %T into DateOnly", value)
}

// Value (Struct → DB)
func (d DateOnly) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Format("2006-01-02"), nil
}
