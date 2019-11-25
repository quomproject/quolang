package quofn

import (
	"math/big"

	"github.com/quomproject/quolang/quoty"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var AddFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "a",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
		{
			Name:             "b",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(quoty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		abr := args[0].EncapsulatedValue().(*big.Rat)
		bbr := args[1].EncapsulatedValue().(*big.Rat)
		var retbr big.Rat
		retbr.Add(abr, bbr)
		return quoty.NumberVal(&retbr), nil
	},
})

var SubtractFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "a",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
		{
			Name:             "b",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(quoty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		abr := args[0].EncapsulatedValue().(*big.Rat)
		bbr := args[1].EncapsulatedValue().(*big.Rat)
		var retbr big.Rat
		retbr.Sub(abr, bbr)
		return quoty.NumberVal(&retbr), nil
	},
})

var MultiplyFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "a",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
		{
			Name:             "b",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(quoty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		abr := args[0].EncapsulatedValue().(*big.Rat)
		bbr := args[1].EncapsulatedValue().(*big.Rat)
		var retbr big.Rat
		retbr.Mul(abr, bbr)
		return quoty.NumberVal(&retbr), nil
	},
})

var DivideFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "a",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
		{
			Name:             "b",
			Type:             quoty.Number,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(quoty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		abr := args[0].EncapsulatedValue().(*big.Rat)
		bbr := args[1].EncapsulatedValue().(*big.Rat)
		var retbr big.Rat
		retbr.Quo(abr, bbr)
		return quoty.NumberVal(&retbr), nil
	},
})

// Add returns the sum of the two given numbers.
func Add(a cty.Value, b cty.Value) (cty.Value, error) {
	return AddFunc.Call([]cty.Value{a, b})
}

// Subtract returns the difference between the two given numbers.
func Subtract(a cty.Value, b cty.Value) (cty.Value, error) {
	return SubtractFunc.Call([]cty.Value{a, b})
}

// Multiply returns the product of the two given numbers.
func Multiply(a cty.Value, b cty.Value) (cty.Value, error) {
	return MultiplyFunc.Call([]cty.Value{a, b})
}

// Divide returns a divided by b, where both a and b are numbers.
func Divide(a cty.Value, b cty.Value) (cty.Value, error) {
	return DivideFunc.Call([]cty.Value{a, b})
}
