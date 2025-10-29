package enumx

import (
	"bytes"
	"database/sql/driver"
	_ "embed"
	"go/types"
	"strings"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"

	"github.com/xoctopus/enumx/pkg/enumx"
)

//go:embed enumx.go.tpl
var template []byte

func init() {
	genx.Register(&g{})
}

type g struct {
	enums *Enums
}

func (x *g) Identifier() string {
	return "enum"
}

func (x *g) New(c genx.Context) genx.Generator {
	return &g{enums: NewEnums(c)}
}

func (x *g) Generate(c genx.Context, t types.Type) error {
	if e, ok := x.enums.Resolve(t); ok {
		if e.IsValid() {
			x.generate(c, e)
			return nil
		}
	}
	return nil
}

func (x *g) generate(c genx.Context, e *Enum) {
	ctx := c.Context()

	ident := s.IdentTT(ctx, e.typ)
	pkgid := "github.com/xoctopus/enumx/pkg/enumx"

	args := []*s.TArg{
		// @def bytes.ToUpper
		s.ArgExpose(ctx, "bytes", "ToUpper"),
		// @def fmt.Sprintf
		s.ArgExpose(ctx, "fmt", "Sprintf"),
		// @def fmt.Sscanf
		s.ArgExpose(ctx, "fmt", "Sscanf"),
		// @def EnumerationType github.com/xoctopus/enumx/pkg.Enum[Type]
		s.ArgExpose(ctx, pkgid, "Enum", ident).WithName("EnumerationType"),
		// @def github.com/xoctopus/enumx/pkg.Scan
		s.ArgExpose(ctx, pkgid, "Scan"),
		// @def github.com/xoctopus/enumx.ParseErrorFor
		s.ArgExpose(ctx, pkgid, "ParseErrorFor", ident),

		// @def Type
		s.Arg(ctx, "Type", ident),
		// @def database/sql/driver.Value
		s.ArgT[driver.Value](ctx),
		// @def github.com/xoctopus/enumx.DriverValueOffset
		s.ArgT[enumx.DriverValueOffset](ctx),

		// @def UnknownValue
		s.Arg(ctx, "UnknownValue", s.ExposeObject(ctx, e.unknown.Exposer())),
		// @def Values
		s.Arg(ctx, "Values", e.Values(ctx)),
		// @def NameToValueCases
		s.Arg(ctx, "StringToValueCases", e.StringToValueCases(ctx)),
		// @def ValueToDescCases
		s.Arg(ctx, "ValueToTextCases", e.ValueToTextCases(ctx)),
		// @def ValueToNameCases
		s.Arg(ctx, "ValueToStringCases", e.ValueToStringCases(ctx)),
	}
	ss := []s.Snippet{s.Template(bytes.NewReader(template), args...)}

	for _, attr := range e.Attrs() {
		if v := strings.ToLower(attr); v != "text" && v != "string" {
			ss = append(ss, e.Attr(ctx, attr))
		}
	}

	c.Render(s.Snippets(s.NewLine(1), ss...))
}
