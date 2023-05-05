package bits

type Bits interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint
}

func Set[T Bits](b, flag T) T    { return b | flag }
func Clear[T Bits](b, flag T) T  { return b &^ flag }
func Toggle[T Bits](b, flag T) T { return b ^ flag }
func Has[T Bits](b, flag T) bool { return b&flag != 0 }
