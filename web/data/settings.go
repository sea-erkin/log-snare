package data

import (
	"gorm.io/gorm"
	"sync"
)

var (
	mu sync.Mutex

	_EnableValidation bool
)

type SettingValue struct {
	gorm.Model
	Key   string `gorm:"index:idx_key,unique"`
	Value bool
}

func SetValidation(value bool) {
	mu.Lock()
	defer mu.Unlock()
	_EnableValidation = value
}

func ValidationEnabled() bool {
	mu.Lock()
	defer mu.Unlock()
	return _EnableValidation
}
