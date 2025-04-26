package utils

import (
	"encoding/json"
)

func MustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
