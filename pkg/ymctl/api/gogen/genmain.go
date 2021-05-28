package gogen

import (
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"path"
	"strings"
)

const mainTemplate = `package main

import (
	"flag"
	"github.com/yametech/beetle/log"
	
	{{.importPackages}}
)

const serviceName = "greet"

var (
	serverAddress string
	registryAddress string
)

func main() {
	flag.StringVar(&serverAddress,"server_address",":19999","grpc server endpoint")
	flag.StringVar(&registryAddress,"registry_address","10.200.100.200:42379","registry address")
	flag.Parse()
	log.InitLogger()
	defer log.Sync()
	config := config.NewConfig(config.Addr(serverAddress), config.RegistryAddr(registryAddress))
	ctx := svc.NewServiceContext(*config)
	server := server.NewGreetServer(serviceName, ctx)
	handler.RegisterHandlers(server)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
`

func genMain(dir string, api *spec.ApiSpec) error {
	serviceName := strings.ToLower(api.Service.Name)
	if strings.HasSuffix(serviceName, "-api") {
		serviceName = strings.ReplaceAll(serviceName, "-api", "")
	}
	subDir := path.Join(cmdDir, serviceName)
	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          subDir,
		filename:        "main.go",
		templateName:    "mainTemplate",
		category:        category,
		templateFile:    mainTemplateFile,
		builtinTemplate: mainTemplate,
		data: map[string]string{
			"importPackages": genMainImports(parentPkg),
			"serviceName":    api.Service.Name,
		},
	})
}

func genMainImports(parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"", ctlutil.JoinPackages(parentPkg, configDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", ctlutil.JoinPackages(parentPkg, handlerDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", ctlutil.JoinPackages(parentPkg, contextDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", ctlutil.JoinPackages(parentPkg, serverDir)))
	return strings.Join(imports, "\n\t")
}
