package main

import (
	"log"
)

func debug(messages ...interface{}) {
	if os.Getenv("X_DEBUG") == "" {
		return
	}
	log.Println(messages...)
}
