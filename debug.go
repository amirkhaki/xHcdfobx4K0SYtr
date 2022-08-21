package main

import (
	"log"
	"os"
)

func debug(messages ...interface{}) {
	if os.Getenv("X_DEBUG") == "" {
		return
	}
	log.Println(messages...)
}
