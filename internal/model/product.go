package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProductCategory int16

const (
	ProductCategoryElectronics ProductCategory = 1
	ProductCategoryClothes     ProductCategory = 2
	ProductCategoryShoes       ProductCategory = 3
)

type Product struct {
	ID          ProductID
	ReceptionID ReceptionID
	Category    ProductCategory
	AddedAt     time.Time
}

type ProductID uuid.UUID

func (s ProductCategory) String() string {
	switch s {
	case ProductCategoryElectronics:
		return "электроника"
	case ProductCategoryClothes:
		return "одежда"
	case ProductCategoryShoes:
		return "обувь"
	}
	return ""
}

func NewProductID() ProductID {
	return ProductID(uuid.New())
}

func (id ProductID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func ParseProductID(s string) (ProductID, error) {
	ID, err := uuid.Parse(s)
	if err != nil {
		return ProductID{}, fmt.Errorf("uuid.parse: %w", err)
	}

	return ProductID(ID), nil
}

func ParseProductCategory(s string) (ProductCategory, error) {
	switch s {
	case "электроника":
		return ProductCategoryElectronics, nil
	case "одежда":
		return ProductCategoryClothes, nil
	case "обувь":
		return ProductCategoryShoes, nil
	}
	return ProductCategory(0), errors.New("category not found")
}
