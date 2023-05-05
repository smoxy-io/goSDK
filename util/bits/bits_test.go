package bits

import "testing"

const (
	TestBitsA uint = 1 << iota
	TestBitsB
	TestBitsC
	TestBitsD
)

func TestSet(t *testing.T) {
	var mask uint

	mask = Set(mask, TestBitsB)

	if mask != TestBitsB {
		t.Errorf("Set(TestBitsB) = %v, wanted: %v", mask, TestBitsB)
	}

	mask = Set(mask, TestBitsD)

	if mask != (TestBitsB | TestBitsD) {
		t.Errorf("Set(TestBitsD) = %v, wanted: %v", mask, TestBitsB|TestBitsD)
	}
}

func TestClear(t *testing.T) {
	var mask uint

	mask = Set(mask, TestBitsA)
	mask = Set(mask, TestBitsC)

	mask = Clear(mask, TestBitsA)

	hasA := mask&TestBitsA != 0
	hasC := mask&TestBitsC != 0

	if hasA || !hasC {
		t.Errorf("Clear(TestBitsA) = %v, wanted: %v", mask, TestBitsC)
	}

	mask = Clear(mask, TestBitsC)

	if mask != 0 {
		t.Errorf("Clear(TestBitsC) = %v, wanted: %v", mask, 0)
	}

}

func TestHas(t *testing.T) {
	var mask uint

	mask = Set(mask, TestBitsA)
	mask = Set(mask, TestBitsB)
	mask = Set(mask, TestBitsD)

	if ok := Has(mask, TestBitsA); !ok {
		t.Errorf("Has(TestBitsA) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsB); !ok {
		t.Errorf("Has(TestBitsB) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsC); ok {
		t.Errorf("Has(TestBitsC) = %v, wanted: %v", ok, false)
	}

	if ok := Has(mask, TestBitsD); !ok {
		t.Errorf("Has(TestBitsD) = %v, wanted: %v", ok, true)
	}
}

func TestToggle(t *testing.T) {
	var mask uint

	mask = Set(mask, TestBitsA)
	mask = Set(mask, TestBitsB)
	mask = Set(mask, TestBitsD)

	if ok := Has(mask, TestBitsA); !ok {
		t.Errorf("Has(TestBitsA) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsB); !ok {
		t.Errorf("Has(TestBitsB) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsC); ok {
		t.Errorf("Has(TestBitsC) = %v, wanted: %v", ok, false)
	}

	if ok := Has(mask, TestBitsD); !ok {
		t.Errorf("Has(TestBitsD) = %v, wanted: %v", ok, true)
	}

	mask = Toggle(mask, TestBitsC)

	if ok := Has(mask, TestBitsA); !ok {
		t.Errorf("Has(TestBitsA) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsB); !ok {
		t.Errorf("Has(TestBitsB) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsC); !ok {
		t.Errorf("Has(TestBitsC) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsD); !ok {
		t.Errorf("Has(TestBitsD) = %v, wanted: %v", ok, true)
	}

	mask = Toggle(mask, TestBitsC)

	if ok := Has(mask, TestBitsA); !ok {
		t.Errorf("Has(TestBitsA) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsB); !ok {
		t.Errorf("Has(TestBitsB) = %v, wanted: %v", ok, true)
	}

	if ok := Has(mask, TestBitsC); ok {
		t.Errorf("Has(TestBitsC) = %v, wanted: %v", ok, false)
	}

	if ok := Has(mask, TestBitsD); !ok {
		t.Errorf("Has(TestBitsD) = %v, wanted: %v", ok, true)
	}
}
