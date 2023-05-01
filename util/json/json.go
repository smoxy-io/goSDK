package json

import "encoding/json"

func ToString[T any](data T) string {
	b, err := json.Marshal(data)

	if err != nil {
		return ""
	}

	return string(b)
}

func FromString[T any](jsonStr string) T {
	var data T

	_ = json.Unmarshal([]byte(jsonStr), &data)

	return data
}
