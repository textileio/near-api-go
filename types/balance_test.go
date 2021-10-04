package types_test

import (
	"testing"

	fuzz "github.com/google/gofuzz"

	. "github.com/eteu-technologies/near-api-go/pkg/types"
)

func TestNEARToYocto(t *testing.T) {
	var NEAR uint64 = 10

	yoctoValue := NEARToYocto(NEAR)
	orig := YoctoToNEAR(yoctoValue)

	if NEAR != orig {
		t.Errorf("expected: %d, got: %d", NEAR, orig)
	}
}

func TestNEARToYocto_Fuzz(t *testing.T) {
	f := fuzz.New()

	// TODO: ?
	var value uint16

	for i := 0; i < 1000; i++ {
		f.Fuzz(&value)
		newValue := YoctoToNEAR(NEARToYocto(uint64(value)))
		if uint64(value) != newValue {
			t.Errorf("expected: %d, got: %d", value, newValue)
		}
	}
}
