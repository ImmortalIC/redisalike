package storage

import (
	"fmt"
	"time"
)

type dictValue struct {
	value      map[string]string
	expireTime time.Time
}

func (v *dictValue) Value() interface{} {
	return v.value
}

func (v *dictValue) Expired() bool {
	return time.Now().After(v.expireTime)
}

func (v *dictValue) ByKey(key string) (string, error) {
	val, ok := v.value[key]
	if !ok {
		return "", fmt.Errorf("No such key: %s", key)
	}
	return val, nil
}
