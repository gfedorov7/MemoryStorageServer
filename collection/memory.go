package collection

import (
	"MemoryStorageServer/errors"
	"time"
)

type MemoryCollection struct {
	Value     []byte
	CreatedAt time.Time
	TTL       time.Duration
}

func Create(value string, ttl time.Duration) (MemoryCollection, error) {
	if ttl <= 0 {
		return MemoryCollection{}, errors.TTLError{}
	}

	return MemoryCollection{
		Value:     []byte(value),
		TTL:       ttl,
		CreatedAt: getCurrentTime(),
	}, nil
}

func getCurrentTime() time.Time {
	return time.Now()
}

func (mc *MemoryCollection) IsExpired() bool {
	expiredTime := mc.CreatedAt.Add(mc.TTL)
	return expiredTime.Before(time.Now())
}
