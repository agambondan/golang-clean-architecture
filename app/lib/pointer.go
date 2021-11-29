package lib

import "github.com/google/uuid"

func Strptr(s string) *string {
	return &s
}

func Intptr(i int) *int {
	return &i
}

func Int64ptr(i int64) *int64 {
	return &i
}

func Float64ptr(f float64) *float64 {
	return &f
}

func UUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func BoolPtr(b bool) *bool {
	return &b
}
