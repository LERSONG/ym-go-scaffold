package cli

import (
	"errors"
	"fmt"
	"github.com/LERSONG/ym-go-scaffold/pkg/ymctl/rpc/generator"
	"github.com/urfave/cli"
	"path/filepath"
)

func RPCNew(c *cli.Context) error {
	rpcname := c.Args().First()
	ext := filepath.Ext(rpcname)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}
	protoName := rpcname + ".proto"
	filename := filepath.Join(".", rpcname, protoName)
	src, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	err = generator.ProtoTmpl(src)
	if err != nil {
		return err
	}
	rpcGenerator, err := generator.NewDefaultRPCGenerator()
	if err != nil {
		return err
	}
	return rpcGenerator.Generate(src, filepath.Dir(src), nil)
}

func RPC(c *cli.Context) error {
	src := c.String("src")
	out := c.String("dir")
	protoImportPath := c.StringSlice("proto_path")
	if len(src) == 0 {
		return errors.New("missing -src")
	}
	if len(out) == 0 {
		return errors.New("missing -dir")
	}

	g, err := generator.NewDefaultRPCGenerator()
	if err != nil {
		return err
	}

	return g.Generate(src, out, protoImportPath)

}
