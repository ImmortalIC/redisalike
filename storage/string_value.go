package storage

import (
	"fmt"
	"time"
)

type stringValue struct {
	value      string
	expireTime time.Time
}

func (v *stringValue) Value() interface{} {
	return v.value
}

func (v *stringValue) Expired() bool {
	return time.Now().After(v.expireTime)
}

func (v *stringValue) ByKey(key string) (string, error) {
	return "", fmt.Errorf("This value are not list or dictionary")
}
