package errors

import "fmt"

type ExpiredError struct {
	Arg any
}

func (e ExpiredError) Error() string {
	return fmt.Sprintf("%v - key expired", e.Arg)
}

type NotFoundError struct {
	Arg any
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%v - key not found", e.Arg)
}

type TTLError struct{}

func (e TTLError) Error() string {
	return fmt.Sprint("ttl must be greater than zero")
}
