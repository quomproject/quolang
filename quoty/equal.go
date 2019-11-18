package quoty

import (
	"github.com/zclconf/go-cty/cty"
)

// RawEqual returns true if the two values are "raw equal" in the same sense
// that cty.Value.RawEqual means it, extended to include equality for
// NumberType and StellarAssetAmountType.
func RawEqual(a, b cty.Value) bool {
	if !a.Type().Equals(b.Type()) {
		return false
	}
	if a.IsKnown() != b.IsKnown() {
		return false
	}
	if a.IsNull() != b.IsNull() {
		return false
	}

	switch {
	case a.Type().Equals(NumberType):
		return numberRawEquals(a, b)
	default:
		return a.RawEquals(b)
	}
}
