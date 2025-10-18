@def Type
@def EnumerationType
--Assertion
var _ #EnumerationType# = (*#Type#)(nil)

@def Type
@def StringToValueCases
@def fmt.Sscanf
@def UnknownValue
@def github.com/xoctopus/enumx/pkg/enumx.ParseErrorFor
--Parse
// Parse#Type# parse #Type# from key
func Parse#Type#(key string) (#Type#, error) {
	switch key {
	#StringToValueCases#
	default:
		var v #Type#
		if _, err := #fmt.Sscanf#(key, "UNKNOWN_%d", &v); err != nil {
			return v, nil
		}
		return #UnknownValue#, #github.com/xoctopus/enumx/pkg/enumx.ParseErrorFor#(key)
	}
}

@def Type
@def Values
--Values
// Values returns enum value list of #Type#
func (#Type#) Values() []#Type# {
	return []#Type#{
		#Values#
	}
}

@def Type
@def fmt.Sprintf
@def ValueToStringCases
--String
// String returns v's string as key
func (v #Type#) String() string {
	switch v {
	#ValueToStringCases#
	default:
		return #fmt.Sprintf#("UNKNOWN_%d", v)
	}
}

@def Type
@def ValueToTextCases
--Text
// Text returns the description as for human reading
func (v #Type#) Text() string {
	switch v {
	#ValueToTextCases#
	default:
		return v.String()
	}
}

@def Type
@def UnknownValue
--IsZero
// IsZero checks if v is zero
func (v #Type#) IsZero() bool {
	return v == #UnknownValue#
}

@def Type
--MarshalText
// MarshalText implements encoding.TextMarshaler
func (v #Type#) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

@def Type
@def bytes.ToUpper
--UnmarshalText
// UnmarshalText implements encoding.TextUnmarshaler
func (v *#Type#) UnmarshalText(data []byte) error {
	vv, err := Parse#Type#(string(#bytes.ToUpper#(data)))
	if err != nil {
		return err
	}
	*v = vv
	return nil
}

@def Type
@def database/sql/driver.Value
@def github.com/xoctopus/enumx/pkg/enumx.DriverValueOffset
--Value
// Value implements driver.Valuer
func (v #Type#) Value() (#database/sql/driver.Value#, error) {
	offset := 0
	if drv, ok := any(v).(#github.com/xoctopus/enumx/pkg/enumx.DriverValueOffset#); ok {
		offset = drv.Offset()
	}
	return int64(v) + int64(offset), nil
}

@def Type
@def github.com/xoctopus/enumx/pkg/enumx.DriverValueOffset
@def github.com/xoctopus/enumx/pkg/enumx.Scan
--Scan
// Scan implements sql.Scanner
func (v *#Type#) Scan(src any) error {
	offset := 0
	if offsetter, ok := any(v).(#github.com/xoctopus/enumx/pkg/enumx.DriverValueOffset#); ok {
		offset = offsetter.Offset()
	}
	i, err := #github.com/xoctopus/enumx/pkg/enumx.Scan#(src, offset)
	if err != nil {
		return err
	}
	*v = #Type#(i)
	return nil
}
