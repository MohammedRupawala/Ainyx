package models

import "time"

type CreateUserInput struct {
	Name string
	DOB  time.Time
}

type UpdateUserInput struct {
	ID   int32
	Name string
	DOB  time.Time
}

type User struct {
	ID   int32
	Name string
	DOB  time.Time
	Age  int
}