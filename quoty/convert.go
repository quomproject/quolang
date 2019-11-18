package quoty

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// Convert tries to convert the given value to the given type, returning an
// error if no conversion is possible.
func Convert(in cty.Value, out cty.Type) (cty.Value, error) {
	if in.Type().Equals(out) {
		return in, nil
	}

	conv := GetConversionUnsafe(in.Type(), out)
	if conv == nil {
		return cty.NilVal, errors.New(convert.MismatchMessage(in.Type(), out))
	}
	return conv(in)
}

// GetConversionUnsafe tries to find a conversion from the first given type
// to the second given type, returning nil if no such conversion is available.
func GetConversionUnsafe(in cty.Type, out cty.Type) convert.Conversion {
	if in.Equals(out) {
		return nil
	}

	switch {
	case out.Equals(NumberType):
		return prepConv(conversionToNumber(in), out)
	//case in.Equals(NumberType):
	//	return prepConv(conversionFromNumber(out), out)
	//case out.Equals(StellarAssetAmountType):
	//	return prepConv(conversionToStellarAssetAmount(in), out)
	//case in.Equals(StellarAssetAmountType):
	//	return prepConv(conversionFromStellarAssetAmount(out), out)
	default:
		// Fall back on cty's default conversions for any other type
		// combinations.
		return convert.GetConversionUnsafe(in, out)
	}
}

func prepConv(conv convert.Conversion, out cty.Type) convert.Conversion {
	if conv == nil {
		return nil
	}

	// We're mainly just wrapping cty's convert here, but adding some
	// additional rules for our custom types.
	// Unfortunately that means we end up having to re-implement some of
	// the type-agnostic rules that cty convert gives "for free" around
	// dealing with unknown values and the dynamic pseudo-type.
	return func(in cty.Value) (cty.Value, error) {
		if out == cty.DynamicPseudoType {
			return in, nil
		}
		if !in.IsKnown() {
			return cty.UnknownVal(out), nil
		}
		if in.IsNull() {
			return cty.NullVal(out), nil
		}
		return conv(in)
	}
}
