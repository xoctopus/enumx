package enumx

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

type _Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Enum defines an enumeration interface for an int<->string type
// implements encoding.TextMarshaler/TextUnmarshaler for interaction in networking
// implements driver.Valuer/sql.Scanner for interaction with database
type Enum[E _Int] interface {
	// Values presents enum value list
	Values() []E
	// String returns enum key as string identifier
	String() string
	// Text returns enum description for presents
	Text() string
	// IsZero check if v is valid
	IsZero() bool

	encoding.TextMarshaler
	encoding.TextUnmarshaler

	driver.Valuer
	sql.Scanner
}

// DriverValueOffset as an adaptor between code and database
type DriverValueOffset interface {
	Offset() int
}

func Scan(src any, offset int) (int, error) {
	switch v := src.(type) {
	case []byte:
		if len(v) == 0 {
			return offset - offset, nil
		}
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return offset, err
		}
		return int(i) - offset, nil
	case string:
		if len(v) == 0 {
			return offset - offset, nil
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return offset, err
		}
		return int(i) - offset, nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(v).Int()) - offset, nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(v).Uint()) - offset, nil
	default:
		return offset - offset, nil
	}
}

func ParseErrorFor[E _Int](from string) error {
	return &ParseError[E]{from: from}
}

type ParseError[E _Int] struct {
	from string
}

func (e *ParseError[E]) Error() string {
	return fmt.Sprintf("failed to parse `%s` to %s", e.from, reflect.TypeFor[E]())
}

func (e *ParseError[E]) Is(err error) bool {
	var target *ParseError[E]
	return errors.As(err, &target)
}
