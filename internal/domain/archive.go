package domain

import (
	"time"

	"github.com/google/uuid"
)

type Region string

const (
	MOSCOW           Region = "MOSCOW"
	SAINT_PETERSBURG Region = "SAINT_PETERSBURG"
	NOVOSIBIRSK      Region = "NOVOSIBIRSK"
	EKATERINBURG     Region = "EKATERINBURG"
	KAZAN            Region = "KAZAN"
	NIZHNY_NOVGOROD  Region = "NIZHNY_NOVGOROD"
	CHELYABINSK      Region = "CHELYABINSK"
	SAMARA           Region = "SAMARA"
	OMSK             Region = "OMSK"
	ROSTOV_ON_DON    Region = "ROSTOV_ON_DON"
)

type Document struct {
	ID        uuid.UUID
	Document  any
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Request struct {
	ID         uuid.UUID
	DocumentID uuid.UUID
	Sender     string
	Region     Region
	CreatedAt  time.Time
}
