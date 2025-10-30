package uuid

import (
	"testing"

	"github.com/google/uuid"
)

func TestDeterministicUUID(t *testing.T) {
	tests := []struct {
		name    string
		seed    []byte
		wantErr bool
	}{
		{
			name:    "valid seed",
			seed:    []byte("test-seed"),
			wantErr: false,
		},
		{
			name:    "empty seed",
			seed:    []byte(""),
			wantErr: true,
		},
		{
			name:    "nil seed",
			seed:    nil,
			wantErr: true,
		},
		{
			name:    "long seed",
			seed:    []byte("this is a very long seed string that contains many characters"),
			wantErr: false,
		},
		{
			name:    "binary seed",
			seed:    []byte{0x00, 0x01, 0x02, 0x03, 0xff, 0xfe, 0xfd},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeterministicUUID(tt.seed)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeterministicUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify it's a valid UUID
				if got == uuid.Nil {
					t.Error("DeterministicUUID() returned nil UUID")
				}
				// Verify determinism - same seed should produce same UUID
				got2, err2 := DeterministicUUID(tt.seed)
				if err2 != nil {
					t.Errorf("DeterministicUUID() second call error = %v", err2)
				}
				if got != got2 {
					t.Errorf("DeterministicUUID() not deterministic: first call = %v, second call = %v", got, got2)
				}
			}
		})
	}
}

func TestDeterministicUUID_Determinism(t *testing.T) {
	// Test that same input always produces same output
	seed := []byte("consistent-seed")

	results := make([]uuid.UUID, 100)
	for i := 0; i < 100; i++ {
		result, err := DeterministicUUID(seed)
		if err != nil {
			t.Fatalf("DeterministicUUID() iteration %d error = %v", i, err)
		}
		results[i] = result
	}

	// All results should be identical
	first := results[0]
	for i, result := range results {
		if result != first {
			t.Errorf("DeterministicUUID() iteration %d = %v, want %v", i, result, first)
		}
	}
}

func TestDeterministicUUID_Uniqueness(t *testing.T) {
	// Test that different inputs produce different outputs
	seeds := [][]byte{
		[]byte("seed1"),
		[]byte("seed2"),
		[]byte("seed3"),
		[]byte("seed4"),
		[]byte("seed5"),
	}

	results := make(map[uuid.UUID]bool)
	for _, seed := range seeds {
		result, err := DeterministicUUID(seed)
		if err != nil {
			t.Fatalf("DeterministicUUID() seed %s error = %v", seed, err)
		}
		if results[result] {
			t.Errorf("DeterministicUUID() produced duplicate UUID %v for different seeds", result)
		}
		results[result] = true
	}
}

func TestDeterministicUUIDString(t *testing.T) {
	tests := []struct {
		name    string
		seed    string
		wantErr bool
	}{
		{
			name:    "valid seed",
			seed:    "test-seed",
			wantErr: false,
		},
		{
			name:    "empty seed",
			seed:    "",
			wantErr: true,
		},
		{
			name:    "long seed",
			seed:    "this is a very long seed string that contains many characters",
			wantErr: false,
		},
		{
			name:    "binary seed",
			seed:    string([]byte{0x00, 0x01, 0x02, 0x03, 0xff, 0xfe, 0xfd}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeterministicUUIDString(tt.seed)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeterministicUUIDString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify it's a valid UUID string
				if got == uuid.Nil {
					t.Error("DeterministicUUIDString() returned nil UUID string")
				}
				// Verify determinism - same seed should produce same UUID string
				got2, err2 := DeterministicUUIDString(tt.seed)
				if err2 != nil {
					t.Errorf("DeterministicUUIDString() second call error = %v", err2)
				}
				if got != got2 {
					t.Errorf("DeterministicUUIDString() not deterministic: first call = %v, second call = %v", got, got2)
				}
			}
		})
	}
}

func TestDeterministicUUIDString_Consistency(t *testing.T) {
	// Test that DeterministicUUIDString returns the string representation of DeterministicUUID
	seed := []byte("test-seed")

	uuidResult, err := DeterministicUUID(seed)
	if err != nil {
		t.Fatalf("DeterministicUUID() error = %v", err)
	}

	stringResult, err := DeterministicUUIDString(string(seed))
	if err != nil {
		t.Fatalf("DeterministicUUIDString() error = %v", err)
	}

	if uuidResult != stringResult {
		t.Errorf("DeterministicUUIDString() = %v, want %v", stringResult, uuidResult)
	}
}

func TestDeterministicUUIDString_Format(t *testing.T) {
	// Test that the string format is correct (8-4-4-4-12)
	seed := "format-test"

	uuidResult, err := DeterministicUUIDString(seed)
	if err != nil {
		t.Fatalf("DeterministicUUIDString() error = %v", err)
	}

	result := uuidResult.String()

	// UUID string format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(result) != 36 {
		t.Errorf("DeterministicUUIDString() length = %d, want 36", len(result))
	}

	// Check for dashes at correct positions
	if result[8] != '-' || result[13] != '-' || result[18] != '-' || result[23] != '-' {
		t.Errorf("DeterministicUUIDString() format incorrect: %s", result)
	}
}

func BenchmarkDeterministicUUID(b *testing.B) {
	seed := []byte("benchmark-seed")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DeterministicUUID(seed)
	}
}

func BenchmarkDeterministicUUIDString(b *testing.B) {
	seed := "benchmark-seed"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DeterministicUUIDString(seed)
	}
}

func BenchmarkDeterministicUUID_LargeSeed(b *testing.B) {
	seed := make([]byte, 10000)
	for i := range seed {
		seed[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DeterministicUUID(seed)
	}
}
