//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type Orders struct {
	ID              uuid.UUID `sql:"primary_key"`
	Amount          float64
	Currency        string
	PaymentAddress  string
	Status          OrderStatus
	TxHash          *string
	CreatedAt       time.Time
	OrderExpiration time.Time
}
