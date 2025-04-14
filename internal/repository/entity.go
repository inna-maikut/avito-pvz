package repository

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID           uuid.UUID `db:"id"`
	City         string    `db:"city"`
	RegisteredAt time.Time `db:"registered_at"`
}

type Reception struct {
	ID         uuid.UUID `db:"id"`
	PVZID      uuid.UUID `db:"pvz_id"`
	Status     int16     `db:"status"`
	ReceptedAt time.Time `db:"recepted_at"`
}

type Product struct {
	ID          uuid.UUID `db:"id"`
	ReceptionID uuid.UUID `db:"reception_id"`
	Category    int16     `db:"category"`
	AddedAt     time.Time `db:"added_at"`
}
