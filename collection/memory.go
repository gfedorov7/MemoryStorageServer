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
	CreatedAt int64
	TTL       int
}

func Create(value string, valueType ValueType, ttl int) (*MemoryCollection, error) {
	if ttl <= 0 {
		return nil, errors.TTLError{}
	}

	return &MemoryCollection{
		Value:     []byte(value),
		TTL:       ttl,
		ValueType: valueType,
		CreatedAt: getCurrentTime(),
	}, nil
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}

func (mc *MemoryCollection) IsExpired() bool {
	expiredTime := mc.CreatedAt + int64(mc.TTL)
	return expiredTime < time.Now().Unix()
}
