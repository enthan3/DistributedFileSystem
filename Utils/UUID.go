package Utils

import (
	"github.com/google/uuid"
)

func GenerateUniqueID() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "Cannot Generate UUID", err
	}
	return newUUID.String(), nil
}
