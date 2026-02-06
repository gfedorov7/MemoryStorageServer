package collection

import (
	"testing"
	"time"
)

func TestMemoryCollection_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		Value     []byte
		CreatedAt time.Time
		TTL       time.Duration
		clock     Clock
		expected  bool
	}{
		{
			Value:     []byte("not expired"),
			CreatedAt: now.Add(-500 * time.Millisecond),
			TTL:       time.Second,
			expected:  false,
		},
		{
			Value:     []byte("just expired"),
			CreatedAt: now.Add(-2 * time.Second),
			TTL:       time.Second,
			expected:  true,
		},
		{
			Value:     []byte("expired long ago"),
			CreatedAt: now.Add(-1000 * time.Second),
			TTL:       time.Second,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Value), func(t *testing.T) {
			mc := MemoryCollection{
				Value:     tt.Value,
				CreatedAt: tt.CreatedAt,
				TTL:       tt.TTL,
			}

			result := mc.IsExpired(RealClock{}.Now())
			if result != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
