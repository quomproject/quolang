package quoty

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		in      cty.Value
		wantTy  cty.Type
		want    cty.Value
		wantErr string
	}{
		// To NumberType
		{
			cty.StringVal("1"),
			NumberType,
			MustParseNumberVal("1"),
			``,
		},
		{
			cty.StringVal("1.2"),
			NumberType,
			MustParseNumberVal("1.2"),
			``,
		},
		{
			cty.StringVal("-1.2"),
			NumberType,
			MustParseNumberVal("-1.2"),
			``,
		},
		{
			StellarAssetAmountVal(StellarAssetAmount(25123456)),
			NumberType,
			MustParseNumberVal("2.5123456"),
			``,
		},
		{
			StellarAssetAmountVal(StellarAssetAmount(-25123456)),
			NumberType,
			MustParseNumberVal("-2.5123456"),
			``,
		},
		{
			cty.NumberIntVal(1),
			NumberType,
			MustParseNumberVal("1"),
			``,
		},
		{
			cty.NumberFloatVal(1.5),
			NumberType,
			MustParseNumberVal("1.5"),
			``,
		},
		{
			cty.UnknownVal(cty.String),
			NumberType,
			cty.UnknownVal(NumberType),
			``,
		},
		{
			cty.NullVal(cty.String),
			NumberType,
			cty.NullVal(NumberType),
			``,
		},
		{
			cty.StringVal("hi"),
			NumberType,
			cty.NilVal,
			`a number is required`,
		},
		{
			cty.StringVal("1/2"),
			NumberType,
			cty.NilVal,
			`a number is required`,
		},
		{
			cty.True,
			NumberType,
			cty.NilVal,
			`number required`,
		},

		// Normal cty conversions should still be working
		{
			cty.StringVal("hi"),
			cty.String,
			cty.StringVal("hi"),
			``,
		},
		{
			cty.StringVal("hi"),
			cty.Number,
			cty.NilVal,
			`a number is required`,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v to %#v", test.in, test.wantTy), func(t *testing.T) {
			got, err := Convert(test.in, test.wantTy)

			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if !RawEqual(got, test.want) {
					t.Errorf("wrong result\ngot:  %s\nwant: %s", ValGoString(got), ValGoString(test.want))
				}
			} else {
				if err == nil {
					t.Fatalf("wrong error:\ngot:  <no error>\nwant: %s", test.wantErr)
				}
				if got, want := err.Error(), test.wantErr; got != want {
					t.Fatalf("wrong error:\ngot:  %s\nwant: %s", got, want)
				}
			}
		})
	}
}
