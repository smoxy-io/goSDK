package json

import "encoding/json"

func ToString[T any](data T) string {
	return string(ToBytes(data))
}

func ToBytes[T any](data T) []byte {
	b, err := json.Marshal(data)

	if err != nil {
		return nil
	}

	return b
}

func ToStringPretty[T any](data T) string {
	return string(ToBytesPretty(data))
}

func ToBytesPretty[T any](data T) []byte {
	b, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return nil
	}

	return b
}

func FromString[T any](jsonStr string) T {
	var data T

	_ = json.Unmarshal([]byte(jsonStr), &data)

	return data
}

func FromBytes[T any](jsonBytes []byte) T {
	var data T

	_ = json.Unmarshal(jsonBytes, &data)

	return data
}
