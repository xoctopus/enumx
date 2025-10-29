package enumx

import (
	"context"
	"fmt"
	"go/constant"
	"go/types"
	"sort"
	"strings"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
	"github.com/xoctopus/pkgx"
	"github.com/xoctopus/x/stringsx"
	"golang.org/x/exp/maps"
)

type option struct {
	name  string
	text  string
	attrs map[string]string
	value *pkgx.Constant
}

type Enum struct {
	typ     types.Type
	key     string
	unknown *pkgx.Constant
	values  []*option
	attrs   map[string]struct{}
}

func (e *Enum) IsValid() bool {
	return e.unknown != nil || len(e.values) > 0
}

// add adds option
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
	if len(parts) != 2 || parts[0] != prefix {
		return
	}

	o := &option{
		value: c,
		name:  parts[1],
		text:  "",
		attrs: map[string]string{},
	}

	lines := make([]string, 0)
	for _, desc := range c.Doc().Desc() {
		if strings.HasPrefix(desc, "@attr ") {
			attr := strings.TrimPrefix(desc, "@attr ")
			if idx := strings.Index(attr, "="); idx != -1 && idx != len(attr)-1 {
				k, v := attr[:idx], attr[idx+1:]
				o.attrs[k] = v
				e.attrs[k] = struct{}{}
			}
			continue
		}
		// maybe more attributes, it can prefix with '@' for document parsing
		if desc == c.Name() || strings.HasPrefix(desc, "@") {
			continue
		}
		lines = append(lines, desc)
	}
	o.text = strings.Join(lines, " ")
	e.values = append(e.values, o)
}

// Values generates code snippet of const value list
func (e *Enum) Values(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.values {
		expose := s.ExposeObject(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(2), expose, s.Block(",")),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

// ValueToStringCases generates code snippet cases from enum value to string
func (e *Enum) ValueToStringCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.values {
		name := strings.TrimPrefix(
			v.value.Name(),
			stringsx.UpperSnakeCase(v.value.TypeName())+"__",
		)
		expose := s.ExposeObject(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), s.BlockF("return %q", name)),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

// StringToValueCases generates code snippet cases from string to const value
func (e *Enum) StringToValueCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.values {
		expose := s.ExposeObject(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.BlockF("case %q:", v.name)),
			s.Compose(s.Indent(2), s.Block("return "), expose, s.Block(", nil")),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

// ValueToTextCases generates code snippet cases from enum value to text
func (e *Enum) ValueToTextCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.values {
		text := v.text
		if len(text) == 0 {
			text = v.name
		}
		expose := s.ExposeObject(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), s.BlockF("return %q", text)),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (e *Enum) Attr(ctx context.Context, attr string) s.Snippet {
	ss := make([]s.Snippet, 0)

	name := stringsx.UpperCamelCase(attr)

	ss = append(
		ss,
		s.Comments(fmt.Sprintf("%s describes %s attribute", name, name)),
		s.Compose(s.Block("func (v "), s.IdentTT(ctx, e.typ), s.BlockF(") %s() string {", name)),
		s.Compose(s.Indent(1), s.Block("switch v {")),
	)

	for _, v := range e.values {
		expose := s.ExposeObject(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), s.BlockF("return %q", v.attrs[attr])),
		)
	}

	ss = append(
		ss,
		s.Compose(s.Indent(1), s.Block("default:")),
		s.Compose(s.Indent(2), s.BlockF("return %q", "")),
		s.Compose(s.Indent(1), s.Block("}")),
		s.Compose(s.Block("}\n")),
	)

	return s.Snippets(s.NewLine(1), ss...)
}

func (e *Enum) Attrs() []string {
	attrs := maps.Keys(e.attrs)
	sort.Strings(attrs)
	return attrs
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
				values: make([]*option, 0),
				attrs:  make(map[string]struct{}),
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

func (es *Enums) Resolve(t types.Type) (*Enum, bool) {
	if _, ok := t.(*types.Named); !ok {
		return nil, false
	}
	e, ok := es.e[t]
	return e, ok
}
