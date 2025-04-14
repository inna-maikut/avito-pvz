package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ReceptionStatus int16

const (
	ReceptionStatusInProgress ReceptionStatus = 1
	ReceptionStatusClose      ReceptionStatus = 2
)

type Reception struct {
	ID              ReceptionID
	PVZID           PVZID
	ReceptionStatus ReceptionStatus
	ReceptedAt      time.Time
}

type ReceptionID uuid.UUID

func (s ReceptionStatus) String() string {
	switch s {
	case ReceptionStatusInProgress:
		return "in_progress"
	case ReceptionStatusClose:
		return "close"
	}

	return ""
}

func NewReceptionID() ReceptionID {
	return ReceptionID(uuid.New())
}

func (id ReceptionID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func ParseReceptionID(s string) (ReceptionID, error) {
	ID, err := uuid.Parse(s)
	if err != nil {
		return ReceptionID{}, fmt.Errorf("uuid.parse: %w", err)
	}

	return ReceptionID(ID), nil
}
