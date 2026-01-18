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

func Create(value string, ttl time.Duration, now time.Time) (MemoryCollection, error) {
	if ttl <= 0 {
		return MemoryCollection{}, errors.TTLError{}
	}

	return MemoryCollection{
		Value:     []byte(value),
		TTL:       ttl,
		CreatedAt: now,
	}, nil
}

func (mc *MemoryCollection) IsExpired(now time.Time) bool {
	expiredTime := mc.CreatedAt.Add(mc.TTL)
	return expiredTime.Before(now)
}
