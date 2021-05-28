package generator

import (
	"bytes"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/yametech/ym-go-scaffold/pkg/conf"

	"path/filepath"
	"strings"
)

func (g *DefaultGenerator) GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetPb()
	cw := new(bytes.Buffer)
	base := filepath.Dir(proto.Src)
	cw.WriteString("protoc ")
	for _, ip := range protoImportPath {
		cw.WriteString(" -I=" + ip)
	}
	cw.WriteString(" -I=" + base)
	cw.WriteString(" " + proto.Name)
	if strings.Contains(proto.GoPackage, "/") {
		cw.WriteString(" --go_out=plugins=grpc:" + ctx.GetWd().Filename)
	} else {
		cw.WriteString(" --go_out=plugins=grpc:" + dir.Filename)
	}
	command := cw.String()
	g.log.Debug(command)
	_, err := execx.Run(command, "")
	return err
}
