package quoty

import (
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// StellarAssetAmountType is a more constrained numeric type than NumberType that
// represents Stellar amounts specifically. This type can be safely and
// losslessly converted to NumberType. Not all NumberType values can be
// converted to StellarAssetAmountType.
var StellarAssetAmountType cty.Type

// StellarAssetAmountVal wraps a StallarAssetAmount in a cty.Value of type
// StellarAssetAmountType.
func StellarAssetAmountVal(n StellarAssetAmount) cty.Value {
	return cty.CapsuleVal(StellarAssetAmountType, &n)
}

// StellarAssetAmount is a specialzed int64 that marks the intent for the
// number to represent ten-millionths of the asset unit.
//
// This type exists just to help us make sure we're always intentionally
// marking a specific int64 as being an amount, to avoid inadvertently
// ending up with a value ten million times larger or smaller than what
// we intended.
type StellarAssetAmount int64

func (a StellarAssetAmount) String() string {
	// FIXME: This is a pretty wasteful way to implement this, allocating a
	// temporary object we just throw away.
	br := ratFromStellarAssetVal(a)
	return br.FloatString(7)
}

// StellarAssetType is a cty object type used to represent Stellar assets
// in the Quo language.
//
// Quo has a structural type system, so an asset value is just a normal object
// value that happens to have the "code" and "issuer" attributes as defined
// here.
var StellarAssetType = cty.Object(map[string]cty.Type{
	"code":   cty.String,
	"issuer": cty.String,
})

// StellarNativeAssetVal is a cty value of type StellarAssetType that
// represents XLM (Lumens), the native asset of Stellar.
//
// Note that this is a somewhat artificial representation of XLM, because
// the Stellar network treats this asset as a special case. Any code that
// maps from StellarAssetType values to real Stellar asset types must
// recognize this particular value as special and translate it to the
// specialized native asset representation in the Stellar protocol.
var StellarNativeAssetVal = cty.ObjectVal(map[string]cty.Value{
	"code":   cty.StringVal("XLM"),
	"issuer": cty.NullVal(cty.String),
})

func init() {
	StellarAssetAmountType = cty.Capsule("Stellar asset amount", reflect.TypeOf(StellarAssetAmount(0)))
}
