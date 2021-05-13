package generator

import (
	"fmt"
	"github.com/LERSONG/ym-go-scaffold/pkg/conf"
	"github.com/LERSONG/ym-go-scaffold/pkg/util"
	"github.com/LERSONG/ym-go-scaffold/pkg/util/stringx"
	"github.com/LERSONG/ym-go-scaffold/pkg/ymctl/rpc/parser"
	"path/filepath"
	"strings"
)

const mainTemplate = `package main

import (
	"fmt"
	"google.golang.org/grpc/grpclog"
	
	{{.imports}}
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:9999"
)

func main() {
	ctx := svc.NewServiceContext(config.Config{})
	srv := server.New{{.serviceNew}}Server(ctx)
	fmt.Printf("[INFO] listen rpc (%s)\n", Address)
	err := srv.GrpcRun(Address)
	if err != nil {
		grpclog.Fatalln(err)
	}
}`

func (g *DefaultGenerator) GenMain(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	serviceName := strings.ToLower(proto.Service.Name)
	fileName := filepath.Join(ctx.GetMain().Filename, fmt.Sprintf("%v.go", serviceName))
	imports := make([]string, 0)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	remoteImport := fmt.Sprintf(`"%v"`, ctx.GetServer().Package)
	configImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
	imports = append(imports, svcImport, remoteImport, configImport)

	text, err := util.LoadTemplate(category, mainTemplateFile, mainTemplate)
	if err != nil {
		return err
	}

	return util.With("main").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"imports":    strings.Join(imports, util.NL),
		"serviceNew": stringx.From(proto.Service.Name).ToCamel(),
	}, fileName, false)
}
