package enumx

import "strconv"

type IntStringerEnum interface {
	Typename() string
	Int() int
	String() string
	Label() string
	ConstValues() []IntStringerEnum
}

type Enum = IntStringerEnum

type ValueOffset interface {
	Offset() int
}

func Scan(src any, offset int) (int, error) {
	v, err := AsInteger(src, offset)
	if err != nil {
		return offset, err
	}
	return v - offset, nil
}

func AsInteger(src any, defaults int) (int, error) {
	switch v := src.(type) {
	case []byte:
		if len(v) == 0 {
			return defaults, nil
		}
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return defaults, err
		}
		return int(i), nil
	case string:
		if len(v) == 0 {
			return defaults, nil
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaults, err
		}
		return int(i), nil
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	default:
		return defaults, nil
	}
}
