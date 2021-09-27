package parsers

import (
	"strconv"
	"testing"
)

func TestAtomicFloat64(t *testing.T) {

	tests := []struct {
		name  string
		input string
	}{
		// TODO: Add test cases.
		{"#1", "0.0"},
		{"#2", "0.1"},
		{"#3", "1.2"},
		{"#4", "10.3"},
		{"#5", "1000.4"},
		{"#6", "10000000.5"},
		{"#7", "100000000.6"},
		{"#8", "9999999999.7"},
		{"#9", "101010101010.8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := 0.0

			if err := AtomicFloat64(&out)(tt.input); err != nil {
				t.Errorf("AtomicFloat64() : %s : %v", tt.name, err)
			}

			got := out
			want, err := strconv.ParseFloat(tt.input, 64)
			if err != nil {
				t.Errorf("AtomicFloat64(): %s : %v", tt.name, err)
			}
			if got != want {
				t.Errorf("got = %v want = %v", got, want)
			}
		})
	}
}

func TestAtomicFloat32(t *testing.T) {

	tests := []struct {
		name  string
		input string
	}{
		// TODO: Add test cases.
		{"#1", "0.0"},
		{"#2", "0.1"},
		{"#3", "1.2"},
		{"#4", "10.3"},
		{"#5", "1000.4"},
		{"#6", "10000000.5"},
		{"#7", "100000000.6"},
		{"#8", "9999999999.7"},
		{"#9", "101010101010.8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := float32(0.0)

			if err := AtomicFloat32(&out)(tt.input); err != nil {
				t.Errorf("AtomicFloat32() : %s : %v", tt.name, err)
			}

			got := float64(out)
			want, err := strconv.ParseFloat(tt.input, 32)
			if err != nil {
				t.Errorf("AtomicFloat32(): %s : %v", tt.name, err)
			}
			if got != want {
				t.Errorf("got = %v want = %v", got, want)
			}
		})
	}
}
