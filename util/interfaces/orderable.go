package interfaces

type Orderable interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint | ~int64 | ~int32 | ~int16 | ~int8 | ~int | ~float64 | ~float32 | ~string
}
