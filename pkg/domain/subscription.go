package domain

import (
	"errors"
	"time"
)

type ServiceType int

const (
	ServiceTypeSupplierProgram ServiceType = iota + 1
)

type Subscription struct {
	Status    bool        `json:"status" validate:"required"`
	AppType   ServiceType `json:"service_type" validate:"required,gte=1"`
	StartDate time.Time   `json:"start_date" validate:"required"`
	EndDate   time.Time   `json:"end_date" validate:"required"`
}

func (st ServiceType) IsValid() error {
	switch st {
	case ServiceTypeSupplierProgram:
		return nil
	}
	return errors.New("invalid service type")
}
