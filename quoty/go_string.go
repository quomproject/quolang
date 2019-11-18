package quoty

import (
	"fmt"
	"math/big"

	"github.com/zclconf/go-cty/cty"
)

// ValGoString is an extension of cty.Value.GoString that understands the
// quoty types.
func ValGoString(v cty.Value) string {
	switch {
	case v == cty.NilVal:
		return "cty.NilVal"
	case v.IsNull():
		return fmt.Sprintf("cty.NullVal(%s)", TypeGoString(v.Type()))
	case !v.IsKnown():
		return fmt.Sprintf("cty.UnknownVal(%s)", TypeGoString(v.Type()))
	case v.Type().Equals(NumberType):
		br := v.EncapsulatedValue().(*big.Rat)
		// FIXME: Try to infer automatically how many digits we need, rather
		// than just always using 14 decimal places.
		return fmt.Sprintf("quoty.MustParseNumberVal(%s)", br.FloatString(14))
	default:
		return v.GoString()
	}
}

// TypeGoString is an extension of cty.Type.GoString that understands the
// quoty types.
func TypeGoString(ty cty.Type) string {
	switch {
	default:
		return ty.GoString()
	}
}
