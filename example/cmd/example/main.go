package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/genx/pkg/genx"
	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/enumx/devpkg/enumx"
)

func main() {
	cwd := must.NoErrorV(os.Getwd())

	entry := filepath.Join(cwd, "testdata")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}
}
