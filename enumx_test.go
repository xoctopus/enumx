package enumx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/enumx"
)

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		offset  int
		want    int
		wantErr bool
	}{
		{"ParseBytes", []byte("100"), 10, 90, false},
		{"ParseBytesFailed", []byte("xxx"), 11, 11, true},
		{"ParseEmptyBytes", []byte{}, 18, 0, false},
		{"ParseString", "101", 0, 101, false},
		{"ParseStringFailed", "xxx", 0, 0, true},
		{"ParseEmptyString", "", 19, 0, false},
		{"Int", 10, 0, 10, false},
		{"Int8", 10, 0, 10, false},
		{"Int16", 10, 0, 10, false},
		{"Int32", 10, 0, 10, false},
		{"Int64", 10, 0, 10, false},
		{"Uint", 10, 1, 9, false},
		{"Uint8", 10, 2, 8, false},
		{"Uint16", 10, 3, 7, false},
		{"Uint32", 10, 4, 6, false},
		{"Uint64", 10, 5, 5, false},
		{"Nil", nil, 10, 0, false},
		{"OtherType", reflect.ValueOf(10), 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := enumx.Scan(tt.src, tt.offset)
			if err == nil {
				NewWithT(t).Expect(got).To(Equal(tt.want))
				NewWithT(t).Expect(tt.wantErr).To(BeFalse())
			} else {
				NewWithT(t).Expect(got).To(Equal(tt.offset))
				NewWithT(t).Expect(tt.wantErr).To(BeTrue())
			}
		})
	}
}
