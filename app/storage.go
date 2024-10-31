package main

import (
	"sync"
	"time"
)

type StorageValue struct {
	Value    string
	ExpireAt time.Time
}

var Storage = map[string]StorageValue{}
var StorageMutex = sync.RWMutex{}
