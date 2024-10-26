package entity

import (
	"github.com/google/uuid"
	"github.com/protomem/chatik/pkg/werrors"
)

type ID = uuid.UUID

func GenID() ID {
	return uuid.New()
}

func ParseID(idStr string) (ID, error) {
	id, err := uuid.Parse(idStr)
	return id, werrors.Error(err, "id/parse")
}
