package testdata

import (
	"time"
)

// SimpleSource is a simple source type for testing default mapping
type SimpleSource struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// SimpleTarget is a simple target type for testing default mapping
type SimpleTarget struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// +mapgen:mapper impl:defaultMapper target:default_mapper_impl.go
type DefaultMapper interface {
	// No explicit mapping annotations - should use default field-to-field mapping
	ToTarget(*SimpleSource) *SimpleTarget
}
