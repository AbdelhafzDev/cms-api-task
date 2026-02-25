package uuidutil

import (
	"errors"

	"github.com/google/uuid"
)

func StringToBytes(uuidStr string) ([]byte, error) {
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	return u.MarshalBinary()
}

func BytesToString(b []byte) (string, error) {
	if len(b) != 16 {
		return "", errors.New("invalid UUID bytes length")
	}
	u, err := uuid.FromBytes(b)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func NewV7Bytes() ([]byte, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return u.MarshalBinary()
}

func NewV7String() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
