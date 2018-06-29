package main

import (
	"flag"
	"fmt"
	"github.com/ImmortalIC/redisalike/client"
	"os"
	"strings"
)

var action, key, payload, index string

func init() {
	flag.StringVar(&action, "action", "", "Choose action for client. Possible actions: get, set, keys, remove")
	flag.StringVar(&key, "key", "", "Key in storage. Required for get,set and remove")
	flag.StringVar(&payload, "data", "", "JSON or string for saving in storage. Required for set")
	flag.StringVar(&index, "index", "", "Index in list or dictionary for fetching. Optional for get")
}

func main() {
	flag.Parse()
	switch action {
	case "get":
		if key == "" {
			flag.Usage()
			os.Exit(1)
		}
		value, err := client.Get(key, index)
		if err != nil {
			fmt.Printf("Something gone wrong. Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Value of key %s is %s\n", key, value)
	case "set":
		if key == "" || payload == "" {
			flag.Usage()
			os.Exit(1)
		}
		err := client.Set(key, payload)
		if err != nil {
			fmt.Printf("Something gone wrong. Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Your value successfully saved on key %s\n", key)
	case "remove":
		if key == "" {
			flag.Usage()
			os.Exit(1)
		}
		err := client.Remove(key)
		if err != nil {
			fmt.Printf("Something gone wrong. Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Key %s successfuly removed from storage.\n", key)
	case "keys":
		keys, err := client.Keys()
		if err != nil {
			fmt.Printf("Something gone wrong. Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Keys contained in storage %s", strings.Join(keys, "\n"))
	default:
		flag.Usage()
		os.Exit(1)
	}
}
