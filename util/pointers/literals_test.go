package pointers

import "testing"

func TestLiteral(t *testing.T) {
	var b *bool
	var s *string
	var i *int
	var u *uint
	var f32 *float32
	var f64 *float64
	var i8 *int8
	var i16 *int16
	var i32 *int32
	var i64 *int64
	var u8 *uint8
	var u16 *uint16
	var u32 *uint32
	var u64 *uint64

	if b = Literal[bool](true); *b != bool(true) {
		t.Errorf("*Literal[bool](true) = %v, wanted: %v", *b, true)
	}
	if b = Literal[bool](false); *b != bool(false) {
		t.Errorf("*Literal[bool](false) = %v, wanted: %v", *b, false)
	}
	if s = Literal[string]("foo"); *s != string("foo") {
		t.Errorf("*Literal[string]('foo') = %v, wanted: %v", *s, "foo")
	}
	if i = Literal[int](10); *i != int(10) {
		t.Errorf("*Literal[int](10) = %v, wanted: %v", *i, 10)
	}
	if u = Literal[uint](10); *u != uint(10) {
		t.Errorf("*Literal[uint](10) = %v, wanted: %v", *u, 10)
	}
	if f32 = Literal[float32](4.5); *f32 != float32(4.5) {
		t.Errorf("*Literal[float32](4.5) = %v, wanted: %v", *f32, 4.5)
	}
	if f64 = Literal[float64](4.5); *f64 != float64(4.5) {
		t.Errorf("*Literal[float64](4.5) = %v, wanted: %v", *f64, 4.5)
	}
	if i8 = Literal[int8](4); *i8 != int8(4) {
		t.Errorf("*Literal[int8](4) = %v, wanted: %v", *i8, 4)
	}
	if i16 = Literal[int16](4); *i16 != int16(4) {
		t.Errorf("*Literal[int16](4) = %v, wanted: %v", *i16, 4)
	}
	if i32 = Literal[int32](4); *i32 != int32(4) {
		t.Errorf("*Literal[int32](4) = %v, wanted: %v", *i32, 4)
	}
	if i64 = Literal[int64](4); *i64 != int64(4) {
		t.Errorf("*Literal[int64](4) = %v, wanted: %v", *i64, 4)
	}
	if u8 = Literal[uint8](4); *u8 != uint8(4) {
		t.Errorf("*Literal[uint8](4) = %v, wanted: %v", *u8, 4)
	}
	if u16 = Literal[uint16](4); *u16 != uint16(4) {
		t.Errorf("*Literal[uint16](4) = %v, wanted: %v", *u16, 4)
	}
	if u32 = Literal[uint32](4); *u32 != uint32(4) {
		t.Errorf("*Literal[uint32](4) = %v, wanted: %v", *u32, 4)
	}
	if u64 = Literal[uint64](4); *u64 != uint64(4) {
		t.Errorf("*Literal[uint64](4) = %v, wanted: %v", *u64, 4)
	}
}
