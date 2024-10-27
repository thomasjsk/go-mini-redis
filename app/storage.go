package main

import (
	"sync"
	"time"
)

type StorageValue struct {
	Value    string
	ExpireAt time.Time
}

var SetStorage = map[string]StorageValue{}
var SetMutex = sync.RWMutex{}
