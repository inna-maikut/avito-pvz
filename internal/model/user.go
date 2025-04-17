package model

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	UserID   UserID
	Email    string
	Password string
	UserRole UserRole
}

type UserID uuid.UUID

var DefaultUserID = UserID(uuid.UUID{})

func NewUserID() UserID {
	return UserID(uuid.New())
}

func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func ParseUserID(s string) (UserID, error) {
	ID, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, fmt.Errorf("uuid.parse: %w", err)
	}

	return UserID(ID), nil
}
