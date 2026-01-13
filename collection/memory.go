package collection

import (
	"MemoryStorageServer/errors"
	"time"
)

type ValueType uint8

const (
	TypeString ValueType = iota + 1
	TypeInt
	TypeFloat
)

type MemoryCollection struct {
	Value     []byte
	ValueType ValueType
	CreatedAt time.Time
	TTL       time.Duration
}

func Create(value string, valueType ValueType, ttl time.Duration) (MemoryCollection, error) {
	if ttl <= 0 {
		return MemoryCollection{}, errors.TTLError{}
	}

	return MemoryCollection{
		Value:     []byte(value),
		TTL:       ttl,
		ValueType: valueType,
		CreatedAt: getCurrentTime(),
	}, nil
}

func getCurrentTime() time.Time {
	return time.Now()
}

func (mc *MemoryCollection) IsExpired() bool {
	expiredTime := mc.CreatedAt.Unix() + int64(mc.TTL.Seconds())
	return expiredTime < time.Now().Unix()
}
