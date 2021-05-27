package generator

import (
	"fmt"
	"github.com/yametech/ym-go-scaffold/pkg/conf"
	"github.com/yametech/ym-go-scaffold/pkg/util"
	"github.com/yametech/ym-go-scaffold/pkg/util/stringx"
	"github.com/yametech/ym-go-scaffold/pkg/ymctl/rpc/parser"
	"path/filepath"
	"strings"
)

const mainTemplate = `package main

import (
	"flag"
	"github.com/yametech/beetle/log"
	
	{{.imports}}
)

const serviceName = "{{.serviceName}}"

var (
	serverAddress string
	registryAddr string
)

func main() {
	flag.StringVar(&serverAddress,"server_address",":18888","grpc server endpoint")
	flag.StringVar(&registryAddr,"registry_address","10.200.100.200:42379","registry address")
	flag.Parse()
	log.InitLogger()
	defer log.Sync()
	config := config.NewConfig(
		config.Addr(serverAddress),
		config.RegistryAddr(registryAddr),
		)
	ctx := svc.NewServiceContext(*config)
	srv := server.New{{.serviceNew}}Server(serviceName, ctx)

	if err := srv.GrpcRun(); err != nil {
		panic(err)
	}
}`

func (g *DefaultGenerator) GenMain(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	//serviceName := strings.ToLower(proto.Service.Name)
	fileName := filepath.Join(ctx.GetMain().Filename, "main.go")
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
		"imports":     strings.Join(imports, util.NL),
		"serviceName": stringx.From(proto.Service.Name).ToSnake(),
		"serviceNew":  stringx.From(proto.Service.Name).ToCamel(),
	}, fileName, false)
}
