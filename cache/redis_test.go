package cache

import "testing"

func TestNewConnection(t *testing.T) {
	_, err := NewCache("localhost", 6379)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
}
