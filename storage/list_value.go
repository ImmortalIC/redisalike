package storage

import (
	"fmt"
	"strconv"
	"time"
)

type listValue struct {
	value      []string
	expireTime time.Time
}

func (v *listValue) Value() interface{} {
	return v.value
}

func (v *listValue) Expired() bool {
	return time.Now().After(v.expireTime)
}

func (v *listValue) ByKey(key string) (string, error) {
	index, err := strconv.Atoi(key)
	if err != nil {
		return "", err
	}
	if len(v.value) >= index {
		return "", fmt.Errorf("Index %d is out of bounds", index)
	}
	return v.value[index], nil
}
