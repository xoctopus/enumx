package enumx

import (
	"go/constant"
	"go/types"
	"strings"

	"github.com/xoctopus/genx"
	"github.com/xoctopus/pkgx"
	"github.com/xoctopus/x/stringsx"
)

type Option struct {
	name  string
	desc  string
	value int64
}

type Enum struct {
	typ     types.Type
	key     string
	unknown *pkgx.Constant
	values  []*pkgx.Constant
}

func (e *Enum) IsValid() bool {
	return e.unknown != nil || len(e.values) > 0
}

func (e *Enum) add(c *pkgx.Constant) {
	name := c.Name()
	if name[0] == '_' {
		return
	}

	prefix := stringsx.UpperSnakeCase(e.key)
	if name == prefix+"_UNKNOWN" {
		e.unknown = c
		return
	}
	parts := strings.SplitN(name, "__", 2)
	if len(parts) == 2 && parts[0] == prefix {
		e.values = append(e.values, c)
	}
}

func NewEnums(g genx.Context) *Enums {
	es := &Enums{
		e: make(map[types.Type]*Enum),
		p: g.Package(),
	}

	// Elements has been ordered by node(position)
	for elem := range es.p.Constants().Elements() {
		typ := elem.Type()
		if _, ok := typ.(*types.Named); !ok {
			continue
		}
		if elem.Value().Kind() != constant.Int {
			continue
		}

		if _, ok := es.e[typ]; !ok {
			es.e[typ] = &Enum{
				typ:    typ,
				key:    elem.TypeName(),
				values: make([]*pkgx.Constant, 0),
			}
		}
		es.e[typ].add(elem)
	}
	return es
}

type Enums struct {
	p pkgx.Package
	e map[types.Type]*Enum
}

func (es Enums) Resolve(t types.Type) (*Enum, bool) {
	if _, ok := t.(*types.Named); !ok {
		return nil, false
	}
	e, ok := es.e[t]
	return e, ok
}
