package testdata

import (
	"testing"
	"time"
)

func TestDefaultMapper(t *testing.T) {
	// Create a new mapper
	mapper := &defaultMapper{}

	// Create a source object with some values
	now := time.Now()
	source := &SimpleSource{
		ID:        123,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: now,
	}

	// Map the source to a target
	target := mapper.ToTarget(source)

	// Verify that all fields were correctly mapped
	if target.ID != source.ID {
		t.Errorf("ID not mapped correctly. Expected %d, got %d", source.ID, target.ID)
	}

	if target.Name != source.Name {
		t.Errorf("Name not mapped correctly. Expected %s, got %s", source.Name, target.Name)
	}

	if target.Email != source.Email {
		t.Errorf("Email not mapped correctly. Expected %s, got %s", source.Email, target.Email)
	}

	if !target.CreatedAt.Equal(source.CreatedAt) {
		t.Errorf("CreatedAt not mapped correctly. Expected %v, got %v", source.CreatedAt, target.CreatedAt)
	}
}
