package gogen

import (
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/yametech/ym-go-scaffold/pkg/util/stringx"
	"strings"
)

const serverTemplate = `package server

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/beetle/log"
	"github.com/yametech/beetle/registry"
	"github.com/yametech/beetle/registry/etcd"
	
	{{.imports}}
	
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type {{.server}}Server struct {
	Name string
	*gin.Engine
	registry registry.Registry
	SvcCtx *svc.ServiceContext
}

func New{{.server}}Server(name string, svcCtx *svc.ServiceContext) *GreetServer {
	engine := gin.New()
	engine.Use([]gin.HandlerFunc{gin.Logger(), gin.Recovery()}...)
	return &GreetServer{
		Name: name,
		Engine: engine,
		SvcCtx: svcCtx,
	}
}

func (s *{{.server}}Server) Run() error {
	go s.Engine.Run(s.SvcCtx.Config.Addr)
	log.Infof("listen http on %s\n", s.SvcCtx.Config.Addr)
	if err := s.register(); err != nil{
		return err
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-ch
	s.registry.Deregister()
	if i, ok := sig.(syscall.Signal); ok {
		os.Exit(int(i))
	} else {
		os.Exit(0)
	}
	return nil
}

func (s *{{.server}}Server) register() error {
	registry, err := etcd.NewRegistry(
		registry.ServiceName(s.Name),
		registry.ServiceAddr(s.SvcCtx.Config.Addr),
		registry.Addrs(strings.Split(s.SvcCtx.Config.RegistryAddr, ",")),
	)
	if err != nil {
		return err
	}
	s.registry = registry
	err = s.registry.Register()
	return err
}
`

func genServer(dir string, api *spec.ApiSpec) error {
	serviceName := strings.ToLower(api.Service.Name)
	if strings.HasSuffix(serviceName, "-api") {
		serviceName = strings.ReplaceAll(serviceName, "-api", "")
	}
	server := stringx.From(serviceName).ToCamel()

	parentPackage, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          serverDir,
		filename:        serviceName + "server.go",
		templateName:    "serverTemplate",
		category:        category,
		templateFile:    serverTemplateFile,
		builtinTemplate: serverTemplate,
		data: map[string]string{
			"server":  server,
			"imports": genServerImports(parentPackage),
		},
	})
}

func genServerImports(parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"\n", ctlutil.JoinPackages(parentPkg, contextDir)))
	return strings.Join(imports, "\n\t")
}
