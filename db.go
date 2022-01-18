package main

import (
	"github.com/peterbourgon/diskv/v3"
	"fmt"
)


func GetDB() *diskv.Diskv {
	return diskv.New(diskv.Options{})
}

func Set(db *diskv.Diskv, key, value string) error {
	if err := db.Write(key, []byte(value)); err != nil {
		return fmt.Errorf("Error setting key %s: %w", err)
	}
	return nil
}

func Get(db *diskv.Diskv, key string) (string, error) {
	val, err := db.Read(key)
	if err != nil {
		return "", fmt.Errorf("Error reading key %s: %w", key, err)
	}
	return string(val), nil
}
func Exists(db *diskv.Diskv, key string) bool {
	_, err := Get(db, key)
	return err == nil
}
func Delete(db *diskv.Diskv, key string) error {
	if err := db.Erase(key); err != nil {
		return fmt.Errorf("Error during deleting key %s: %w", key, err)
	}
	return nil
}
