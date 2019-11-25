package quoty

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// Number is the specialized number type used for numeric expressions in
// the Quo language, instead of cty's own cty.Number.
//
// Quo numbers are rational numbers, capable of representing any possible
// Stellar price but not restricted to the range of Stellar prices.
var Number cty.Type

var numberOps = &cty.CapsuleOps{
	GoString: func(v interface{}) string {
		br := v.(*big.Rat)
		switch {
		case br.IsInt() && br.Num().IsInt64():
			iv := br.Num().Int64()
			switch iv {
			case 0:
				return "quoty.Zero"
			default:
				return fmt.Sprintf("quoty.NumberIntVal(%d)", iv)
			}
		default:
			// FIXME: Try to infer automatically how many digits we need, rather
			// than just always using 14 decimal places.
			return fmt.Sprintf("quoty.MustParseNumberVal(%q)", br.FloatString(14))
		}
	},
	TypeGoString: func(ty reflect.Type) string {
		return "quoty.Number"
	},
	RawEquals: func(a, b interface{}) bool {
		brA := a.(*big.Rat)
		brB := b.(*big.Rat)
		return brA.Cmp(brB) == 0
	},
	ConversionTo: func(srcTy cty.Type) func(cty.Value, cty.Path) (interface{}, error) {
		switch {
		case srcTy.Equals(cty.Number):
			// We don't use cty.Number directly in Quo, but we'll convert from
			// it just in case it shows up from an integration with some other
			// HCL-based system.
			return func(in cty.Value, path cty.Path) (interface{}, error) {
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
				return br, nil
			}
		case srcTy.Equals(cty.String):
			return func(in cty.Value, path cty.Path) (interface{}, error) {
				v, err := ParseNumberVal(in.AsString())
				if err != nil {
					return cty.NilVal, errors.New("a number is required")
				}
				return v.EncapsulatedValue(), nil
			}
		case srcTy.Equals(StellarAssetAmountType):
			return func(in cty.Value, path cty.Path) (interface{}, error) {
				inAmtPtr := in.EncapsulatedValue().(*StellarAssetAmount)
				return numberValFromStellarAssetVal(*inAmtPtr).EncapsulatedValue(), nil
			}
		default:
			return nil
		}
	},
}

// Zero is a Number value representing zero.
var Zero cty.Value

// NumberVal wraps the given rational in a NumberType value.
//
// big.Rat is the native internal type of NumberVal. Callers must not mutate
// the given value after passing it to this function.
func NumberVal(v *big.Rat) cty.Value {
	return cty.CapsuleVal(Number, v)
}

// NumberIntVal converts the given integer to a NumberType value.
func NumberIntVal(v int64) cty.Value {
	br := big.NewRat(v, 1)
	return cty.CapsuleVal(Number, br)
}

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
	return cty.CapsuleVal(Number, br), nil
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
	return cty.CapsuleVal(Number, br)
}

func init() {
	Number = cty.CapsuleWithOps(
		"number",
		reflect.TypeOf(big.Rat{}),
		numberOps,
	)
	Zero = NumberIntVal(0)
}
