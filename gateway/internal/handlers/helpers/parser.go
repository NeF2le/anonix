package helpers

import "github.com/google/uuid"

func ParseUUID(value string) (uuid.UUID, error) {
	uuidField, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidField, nil
}
