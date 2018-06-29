package storage_test

import (
	"context"
	"github.com/ImmortalIC/redisalike/storage"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	storage.Init(ctx, wg)
	m.Run()
	cancel()
	wg.Wait()

}

func TestStorage(t *testing.T) {
	wg := new(sync.WaitGroup)
	testValues := map[string]interface{}{
		"key1":  "abyrvalg",
		"key2":  []string{"aaaaa", "bbbb", "cccc"},
		"key3":  map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key4":  "abyrvalg",
		"key5":  []string{"aaaaa", "bbbb", "cccc"},
		"key6":  map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key7":  "abyrvalg",
		"key8":  []string{"aaaaa", "bbbb", "cccc"},
		"key9":  map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key10": "abyrvalg",
		"key11": []string{"aaaaa", "bbbb", "cccc"},
		"key12": map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key13": "abyrvalg",
		"key14": []string{"aaaaa", "bbbb", "cccc"},
		"key15": map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key16": "abyrvalg",
		"key17": []string{"aaaaa", "bbbb", "cccc"},
		"key18": map[string]string{"asdas": "ssss", "ssss": "qqqq"},
	}
	var count int64
	for key, value := range testValues {
		wg.Add(1)
		go func(key string, value interface{}) {
			defer wg.Done()
			err := storage.Set(key, value, nil)
			if err == nil {
				atomic.AddInt64(&count, 1)
			}
		}(key, value)
	}
	wg.Wait()
	if count != 18 {
		t.Errorf("Only %d of inserts completed successfully", count)
		t.FailNow()
	}
	testValues2 := map[string]interface{}{
		"key4231":  "abyrvalg",
		"key2":     []string{"aaaaa", "bbbb", "cccc"},
		"key33":    map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key4":     "abyrvalg",
		"key435":   []string{"aaaaa", "bbbb", "cccc"},
		"key346":   map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key7":     "abyrvalg",
		"key843":   []string{"aaaaa", "bbbb", "cccc"},
		"key954":   map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key1540":  "abyrvalg",
		"key1541":  []string{"aaaaa", "bbbb", "cccc"},
		"key142":   map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key143":   "abyrvalg",
		"key134":   []string{"aaaaa", "bbbb", "cccc"},
		"key15":    map[string]string{"asdas": "ssss", "ssss": "qqqq"},
		"key1336":  "abyrvalg",
		"key12327": []string{"aaaaa", "bbbb", "cccc"},
		"key1338":  map[string]string{"asdas": "ssss", "ssss": "qqqq"},
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for key, value := range testValues2 {
			err := storage.Set(key, value, nil)
			if err != nil {
				t.Errorf("Cant insert key %s, got error %v", key, err)
				t.FailNow()
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for key, value := range testValues {
			val, err := storage.Get(key)
			if err != nil {
				t.Errorf("Cant read key %s, got error %v", key, err)
				t.FailNow()
			}
			if !reflect.DeepEqual(val.Value(), value) {
				t.Errorf("%v not equals %v", val.Value(), value)
				t.FailNow()
			}
		}
	}()
	wg.Wait()

}
