package pointers

type Literals interface {
	bool | string | int | uint | float32 | float64 | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

func Literal[T Literals](v T) *T {
	val := v
	return &val
}
