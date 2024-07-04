package utils

import "github.com/google/uuid"

func ConvertStringToUUID(id string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}
