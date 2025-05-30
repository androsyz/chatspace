package helper

import (
	"github.com/google/uuid"
)

func StrToUUID(str string) (*uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
