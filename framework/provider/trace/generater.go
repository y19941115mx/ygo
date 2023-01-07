package trace

import "github.com/google/uuid"

type DefaultIDGenerater func() string

func (f DefaultIDGenerater) NewID() string {
	return f()
}

func uuidGenerater() string {
	return uuid.New().String()
}
