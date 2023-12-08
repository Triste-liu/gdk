package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Bool bool

type UnixTime int64

type Map map[string]interface{}

func (t *UnixTime) Scan(value interface{}) error {
	v := value.(time.Time)
	*t = UnixTime(v.UnixMilli())
	return nil
}

func (t UnixTime) Value() (driver.Value, error) {
	return time.Unix(int64(t), 0), nil
}

func (b *Bool) Scan(value interface{}) error {
	v := value.(int64)
	if v == 1 {
		*b = true
	}
	return nil
}

func (b Bool) Value() (driver.Value, error) {
	var v int64
	if b {
		v = 1
	}
	return v, nil
}

func (m *Map) Scan(value interface{}) error {
	v := value.([]uint8)
	err := json.Unmarshal(v, &m)
	if err != nil {
		return err
	}
	return nil
}

func (m Map) Value() (driver.Value, error) {
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	v := string(marshal)
	return v, nil
}
