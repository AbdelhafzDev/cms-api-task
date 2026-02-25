package cursor

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

func Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func EncodeInt(value int) string {
	return Encode(strconv.Itoa(value))
}

func Decode(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func DecodeInt(encoded string) (int, error) {
    if encoded == "" {
        return 0, nil
    }

    decoded, err := Decode(encoded)
    if err != nil {
        return 0, err
    }

    return strconv.Atoi(decoded)
}

func EncodePair(t time.Time, id string) string {
	return Encode(t.UTC().Format(time.RFC3339Nano) + "|" + id)
}

func DecodePair(encoded string) (time.Time, string, error) {
	raw, err := Decode(encoded)
	if err != nil {
		return time.Time{}, "", err
	}

	parts := strings.SplitN(raw, "|", 2)
	if len(parts) != 2 {
		return time.Time{}, "", errors.New("invalid cursor format")
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", err
	}

	return t, parts[1], nil
}
