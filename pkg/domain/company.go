package domain

import (
	"github.com/arangodb/go-driver"
)

type Company struct {
	Id     driver.DocumentID `json:"id,omitempty"`
	Status bool              `json:"status,omitempty"`
	Name   string            `json:"name,omitempty"`
	User   driver.DocumentID `json:"user,omitempty"`
}

type CompanyCreateDto struct {
	Name string `json:"name" validate:"required,min=2,max=32"`
}
