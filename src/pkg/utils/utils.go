package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func StringToTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func Unmarshalinterface[T any](unparsedPayload any) (T, error) {
	var payload T
	payloadBytes, err := json.Marshal(unparsedPayload)
	if err != nil {
		return payload, fmt.Errorf("marshal error: %w", err)
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return payload, fmt.Errorf("unmarshal error: %w", err)
	}

	return payload, nil
}

func LogError(msg string, args ...any) {
	log.Printf("ERROR: %s", fmt.Sprintf(msg, args...))
}

func GetTechNameById(id int) string {
	switch id {
	case 1:
		return "Banking"
	case 2:
		return "Manufacturing"
	case 3:
		return "Hyperspace Range"
	case 4:
		return "Experimentation"
	case 5:
		return "Scanning"
	case 6:
		return "Terraforming"
	case 7:
		return "Weapons"
	default:
		return "Unknown"
	}
}
