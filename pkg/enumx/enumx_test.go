package enumx_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/enumx/pkg/enumx"
	"github.com/xoctopus/enumx/testdata"
)

func TestScan(t *testing.T) {
	cases := []struct {
		name   string
		src    any
		offset int
		expect int
		failed bool
	}{
		{"ParseBytes", []byte("100"), 10, 90, false},
		{"ParseBytesFailed", []byte("xxx"), 11, 11, true},
		{"ParseEmptyBytes", []byte{}, 18, 0, false},
		{"ParseString", "101", 0, 101, false},
		{"ParseStringFailed", "xxx", 0, 0, true},
		{"ParseEmptyString", "", 19, 0, false},
		{"Int", 10, 0, 10, false},
		{"Int8", int8(10), 0, 10, false},
		{"Int16", int16(10), 0, 10, false},
		{"Int32", int32(10), 0, 10, false},
		{"Int64", int64(10), 0, 10, false},
		{"Uint", uint(10), 1, 9, false},
		{"Uint8", uint8(10), 2, 8, false},
		{"Uint16", uint16(10), 3, 7, false},
		{"Uint32", uint32(10), 4, 6, false},
		{"Uint64", uint64(10), 5, 5, false},
		{"Nil", nil, 10, 0, false},
		{"OtherType", reflect.ValueOf(10), 0, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := enumx.Scan(c.src, c.offset)
			if err == nil {
				Expect(t, got, Equal(c.expect))
				Expect(t, c.failed, BeFalse())
			} else {
				Expect(t, got, Equal(c.offset))
				Expect(t, c.failed, BeTrue())
			}
		})
	}
}

func TestParseErrorFor(t *testing.T) {
	err := enumx.ParseErrorFor[testdata.Gender]("any")

	target := enumx.ParseErrorFor[int]("some other")
	Expect(t, errors.Is(target, err), BeFalse())

	target = enumx.ParseErrorFor[testdata.Gender]("some other")
	Expect(t, errors.Is(target, err), BeTrue())

	Expect(t, target.Error(), NotEqual(err.Error()))
}
