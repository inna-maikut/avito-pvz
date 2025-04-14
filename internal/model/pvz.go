package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID           PVZID
	City         string
	RegisteredAt time.Time
}

type PVZID uuid.UUID

func NewPVZID() PVZID {
	return PVZID(uuid.New())
}

func (id PVZID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func ParsePVZID(s string) (PVZID, error) {
	ID, err := uuid.Parse(s)
	if err != nil {
		return PVZID{}, fmt.Errorf("uuid.parse: %w", err)
	}

	return PVZID(ID), nil
}
