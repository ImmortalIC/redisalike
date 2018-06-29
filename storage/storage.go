package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const defaultTTL time.Duration = 5 * time.Minute

type ValueKeeper interface {
	Value() interface{}
	Expired() bool
	ByKey(string) (string, error)
}

type response struct {
	Value ValueKeeper
	Err   error
}
type taskPackage struct {
	key   string
	value ValueKeeper
	resp  chan response
}

type keyListTask struct {
	resp chan []string
}

var storageObject map[string]ValueKeeper
var readChannel chan taskPackage
var writeChannel chan taskPackage
var deleteChannel chan taskPackage
var keyListChannel chan keyListTask

func storekeeper(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case writeTask := <-writeChannel:
			storageObject[writeTask.key] = writeTask.value
			writeTask.resp <- response{Value: nil, Err: nil}
		case readTask := <-readChannel:
			val, ok := storageObject[readTask.key]
			resp := response{Value: nil, Err: nil}
			if !ok || val.Expired() {
				resp.Err = fmt.Errorf("No such key %s in storage", readTask.key)
				delete(storageObject, readTask.key)
			} else {
				resp.Value = val
			}
			readTask.resp <- resp
		case delTask := <-deleteChannel:
			val, ok := storageObject[delTask.key]
			resp := response{Value: nil, Err: nil}
			if !ok || val.Expired() {
				resp.Err = fmt.Errorf("No such key %s in storage", delTask.key)
			}
			delete(storageObject, delTask.key)
			delTask.resp <- resp
		case keysTask := <-keyListChannel:
			keys := make([]string, 0, len(storageObject))
			for key, value := range storageObject {
				if value.Expired() {
					delete(storageObject, key)
					continue
				}
				keys = append(keys, key)
			}
			keysTask.resp <- keys
		}
	}
}

func Init(ctx context.Context, wg *sync.WaitGroup) {
	storageObject = make(map[string]ValueKeeper)
	writeChannel = make(chan taskPackage)
	readChannel = make(chan taskPackage)
	deleteChannel = make(chan taskPackage)
	keyListChannel = make(chan keyListTask)
	wg.Add(1)
	go func() {
		defer wg.Done()
		storekeeper(ctx)
	}()
}

func Set(key string, value interface{}, ttl *time.Duration) error {
	timeToLife := defaultTTL
	if ttl != nil {
		timeToLife = *ttl
	}
	var keeper ValueKeeper
	switch value.(type) {
	case string:
		keeper = &stringValue{
			value:      value.(string),
			expireTime: time.Now().Add(timeToLife),
		}
	case []string:
		keeper = &listValue{
			value:      value.([]string),
			expireTime: time.Now().Add(timeToLife),
		}
	case map[string]string:
		keeper = &dictValue{
			value:      value.(map[string]string),
			expireTime: time.Now().Add(timeToLife),
		}
	default:
		return fmt.Errorf("Unexpected value type")
	}
	task := taskPackage{
		key:   key,
		value: keeper,
		resp:  make(chan response),
	}
	writeChannel <- task
	response := <-task.resp
	return response.Err
}

func Get(key string) (ValueKeeper, error) {
	task := taskPackage{
		key:  key,
		resp: make(chan response),
	}
	readChannel <- task
	response := <-task.resp
	if response.Err != nil {
		return nil, response.Err
	}
	return response.Value, nil
}

func Remove(key string) error {
	task := taskPackage{
		key:  key,
		resp: make(chan response),
	}
	deleteChannel <- task
	response := <-task.resp
	return response.Err
}

func Keys() []string {
	task := keyListTask{
		resp: make(chan []string),
	}
	keyListChannel <- task
	return <-task.resp
}
