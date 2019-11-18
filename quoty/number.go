package quoty

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// NumberType is the specialized number type used for numeric expressions in
// the Quo language, instead of cty's own cty.Number.
//
// Quo numbers are rational numbers, capable of representing any possible
// Stellar price but not restricted to the range of Stellar prices.
var NumberType cty.Type

// ParseNumberVal attempts to parse the given string as a representation of
// a decimal number using ASCII decimal digits, returning a number value
// if successful. If not successful, an error is returned instead.
//
// If the returned error is nil then the returned value is guaranteed to
// be of type quoty.NumberType.
func ParseNumberVal(s string) (cty.Value, error) {
	// The parsing done by big.Rat.SetString is more liberal than we want in
	// that it permits giving two numbers separated by a slash to give the
	// value as a fraction. We'll pre-screen the given string to make sure
	// it's using the subset of Rat number syntax we're expecting, and then
	// delegate to big.Rat.SetString only if it is.

	for _, c := range s {
		if !((c >= '0' && c <= '9') || c == '.' || c == '-' || c == 'e' || c == 'E') {
			return cty.NilVal, errors.New("invalid number syntax")
		}
	}

	br := new(big.Rat)
	_, ok := br.SetString(s)
	if !ok {
		return cty.NilVal, errors.New("invalid number syntax")
	}
	return cty.CapsuleVal(NumberType, br), nil
}

// MustParseNumberVal is a variant of ParseNumberVal that panics if it fails
// to parse the given string, rather than returning an error.
func MustParseNumberVal(s string) cty.Value {
	v, err := ParseNumberVal(s)
	if err != nil {
		panic(err)
	}
	return v
}

func ratFromStellarAssetVal(v StellarAssetAmount) *big.Rat {
	var br big.Rat
	br.SetFrac64(int64(v), 10000000)
	return &br
}

func numberValFromStellarAssetVal(v StellarAssetAmount) cty.Value {
	br := ratFromStellarAssetVal(v)
	return cty.CapsuleVal(NumberType, br)
}

func numberRawEquals(a, b cty.Value) bool {
	if a.IsKnown() != b.IsKnown() {
		return false
	}
	if a.IsNull() != b.IsNull() {
		return false
	}
	if a.IsNull() || !a.IsKnown() {
		return true
	}

	brA := a.EncapsulatedValue().(*big.Rat)
	brB := b.EncapsulatedValue().(*big.Rat)
	return brA.Cmp(brB) == 0
}

func conversionToNumber(in cty.Type) convert.Conversion {
	switch {
	case in.Equals(cty.Number):
		// We don't use cty.Number directly in Quo, but we'll convert from
		// it just in case it shows up from an integration with some other
		// HCL-based system.
		return func(in cty.Value) (cty.Value, error) {
			// big.Float is a base-2 float, so unless it is an infinity we
			// will be able to convert it, but if the float was originally
			// derived from a string of decimal digits then we might not
			// get exactly what that original string specified.
			bf := in.AsBigFloat()
			br, acc := bf.Rat(nil)
			if acc != big.Exact {
				// Rat is defined to return non-exact only if the input is
				// an infinity.
				return cty.NilVal, errors.New("infinity is not allowed")
			}
			return cty.CapsuleVal(NumberType, br), nil
		}
	case in.Equals(cty.String):
		return func(in cty.Value) (cty.Value, error) {
			v, err := ParseNumberVal(in.AsString())
			if err != nil {
				return cty.NilVal, errors.New("a number is required")
			}
			return v, nil
		}
	case in.Equals(StellarAssetAmountType):
		return func(in cty.Value) (cty.Value, error) {
			inAmtPtr := in.EncapsulatedValue().(*StellarAssetAmount)
			return numberValFromStellarAssetVal(*inAmtPtr), nil
		}
	default:
		return nil
	}
}

func init() {
	NumberType = cty.Capsule("number", reflect.TypeOf(big.Rat{}))
}
