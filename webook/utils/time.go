package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TimeTemplate1 = "2006-01-02 15:04:05" //常规类型
	TimeTemplate2 = "2006/01/02 15:04:05" //其他类型
	TimeTemplate3 = "2006-01-02"          //其他类型
	TimeTemplate4 = "15:04:05"            //其他类型
)

type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
