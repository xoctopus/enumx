package enumx

import (
	"bytes"
	"database/sql/driver"
	_ "embed"
	"go/types"
	"strings"

	"github.com/xoctopus/genx"
	s "github.com/xoctopus/genx/snippet"
	"github.com/xoctopus/x/stringsx"

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
		// @def AssertEnumType
		s.ArgExpose(ctx, pkgid, "Enum", ident).WithName("AssertEnumType"),
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
		s.Arg(
			ctx,
			"UnknownValue",
			s.ExposeObject(ctx, e.unknown.Exposer()),
		),
		// @def Values
		s.Arg(
			ctx,
			"Values",
			s.Snippets(
				s.NewLine(1),
				func() []s.Snippet {
					ss := make([]s.Snippet, 0)
					for _, v := range e.values {
						ss = append(ss, s.Compose(
							s.Indent(2),
							s.ExposeObject(ctx, v.Exposer()),
							s.Block(","),
						))
					}
					return ss
				}()...,
			),
		),
		// @def NameToValueCases
		s.Arg(
			ctx,
			"NameToValueCases",
			s.Snippets(
				s.NewLine(1),
				func() []s.Snippet {
					ss := make([]s.Snippet, 0)
					for _, v := range e.values {
						name := strings.TrimPrefix(
							v.Name(),
							stringsx.UpperSnakeCase(v.TypeName())+"__",
						)
						ss = append(ss, s.Compose(
							s.Indent(1),
							s.BlockF("case %q:", name),
						))
						ss = append(ss, s.Compose(
							s.Indent(2),
							s.Block("return "),
							s.ExposeObject(ctx, v.Exposer()),
							s.Block(", nil"),
						))
					}
					return ss
				}()...,
			),
		),
		// @def ValueToDescCases
		s.Arg(
			ctx,
			"ValueToDescCases",
			s.Snippets(
				s.NewLine(1),
				func() []s.Snippet {
					ss := make([]s.Snippet, 0)
					for _, v := range e.values {
						text := strings.Join(v.Doc().Desc(), " ")
						ss = append(
							ss,
							s.Compose(
								s.Indent(1),
								s.Block("case "),
								s.ExposeObject(ctx, v.Exposer()),
								s.Block(":"),
							),
							s.Compose(
								s.Indent(2),
								s.BlockF("return %q", text),
							),
						)
					}
					return ss
				}()...,
			),
		),
		s.Arg(
			ctx,
			"ValueToNameCases",
			s.Snippets(
				s.NewLine(1),
				func() []s.Snippet {
					ss := make([]s.Snippet, 0)
					for _, v := range e.values {
						name := strings.TrimPrefix(
							v.Name(),
							stringsx.UpperSnakeCase(v.TypeName())+"__",
						)
						ss = append(
							ss,
							s.Compose(
								s.Indent(1),
								s.Block("case "),
								s.ExposeObject(ctx, v.Exposer()),
								s.Block(":"),
							),
							s.Compose(
								s.Indent(2),
								s.BlockF("return %q", name),
							),
						)
					}
					return ss
				}()...,
			),
		),
	}

	c.Render(s.Template(bytes.NewReader(template), args...))
}
