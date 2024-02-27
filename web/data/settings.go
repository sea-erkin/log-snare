package data

import (
	"gorm.io/gorm"
	"sync"
)

// Globals saved in memory
var (
	mu sync.Mutex

	_EnableLogging    bool
	_EnableValidation bool
)

type SettingValue struct {
	gorm.Model
	Key   string `gorm:"index:idx_key,unique"`
	Value bool
}

func SetLogging(value bool) {
	mu.Lock()
	defer mu.Unlock()
	_EnableLogging = value
}

func SetValidation(value bool) {
	mu.Lock()
	defer mu.Unlock()
	_EnableValidation = value
}

func LoggingEnabled() bool {
	mu.Lock()
	defer mu.Unlock()
	return _EnableLogging
}

func ValidationEnabled() bool {
	mu.Lock()
	defer mu.Unlock()
	return _EnableValidation
}
