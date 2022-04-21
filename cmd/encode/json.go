package encode

import (
	"encoding/json"
	"fmt"
)

func JSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal: %v", err)
	}

	return string(b), nil
}
