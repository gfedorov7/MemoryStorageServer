package collection

import (
	"MemoryStorageServer/internal/errors"
	"strconv"
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
	expiredTime := mc.expiredAt()
	return expiredTime.Before(now)
}

func (mc *MemoryCollection) String() string {
	return "value=" + string(mc.Value) + "; ttl=" + strconv.FormatInt(int64(mc.TTL), 10) +
		"; createdAt=" + mc.CreatedAt.String() + " expiredAt=" + mc.expiredAt().String() + ";"
}

func (mc *MemoryCollection) expiredAt() time.Time {
	return mc.CreatedAt.Add(mc.TTL)
}
