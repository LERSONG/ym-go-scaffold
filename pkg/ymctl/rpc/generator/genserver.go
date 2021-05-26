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

const (
	serverTemplate = `{{.head}}

package server

import (
	"context"
	"github.com/LERSONG/beetle/log"
	"github.com/LERSONG/beetle/registry"
	"github.com/LERSONG/beetle/registry/etcd"
	"google.golang.org/grpc"

	{{.imports}}

	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type {{.server}}Server struct {
	Name string
	Registry registry.Registry
	svcCtx *svc.ServiceContext
}

func New{{.server}}Server(name string, svcCtx *svc.ServiceContext) *{{.server}}Server {
	return &{{.server}}Server{
		Name: name,
		svcCtx: svcCtx,
	}
}

func (s *{{.server}}Server) GrpcRun() error {
	addr := s.svcCtx.Config.Addr
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	defer server.GracefulStop()
	pb.Register{{.server}}Server(server,s)
	go server.Serve(listen)

	log.Infof("listen rpc (%s)\n", addr)
	if err = s.register(); err != nil{
		return err
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-ch
	s.Registry.Deregister()
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
		registry.ServiceAddr(s.svcCtx.Config.Addr),
		registry.Addrs(strings.Split(s.svcCtx.Config.RegistryAddr, ",")),
	)
	if err != nil {
		return err
	}
	s.Registry = registry
	err = s.Registry.Register()
	return err
}

{{.funcs}}
`
	functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in {{.request}}) ({{.response}}, error) {
	l := logic.New{{.logicName}}(ctx,s.svcCtx)
	return l.{{.method}}(in)
}
`
)

func (g *DefaultGenerator) GenServer(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetServer()
	imports := make([]string, 0)
	logicImport := fmt.Sprintf(`"%v"`, ctx.GetLogic().Package)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

	imports = append(imports, svcImport, logicImport, pbImport)

	head := util.GetHead(proto.Name)
	service := proto.Service
	serverFilename := strings.ToLower(service.Name) + "server"
	serverFile := filepath.Join(dir.Filename, serverFilename+".go")
	funcList, err := g.genFunctions(proto.PbPackage, service)
	if err != nil {
		return err
	}

	text, err := util.LoadTemplate(category, serverTemplateFile, serverTemplate)
	if err != nil {
		return err
	}

	err = util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head":    head,
		"server":  stringx.From(service.Name).ToCamel(),
		"imports": strings.Join(imports, util.NL),
		"funcs":   strings.Join(funcList, util.NL),
	}, serverFile, true)
	return err
}

func (g *DefaultGenerator) genFunctions(goPackage string, service parser.Service) ([]string, error) {
	var functionList []string
	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, serverFuncTemplateFile, functionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		buffer, err := util.With("func").Parse(text).Execute(map[string]interface{}{
			"server":     stringx.From(service.Name).ToCamel(),
			"logicName":  fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel()),
			"method":     parser.CamelCase(rpc.Name),
			"request":    fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
			"response":   fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
			"hasComment": len(comment) > 0,
			"comment":    comment,
		})
		if err != nil {
			return nil, err
		}

		functionList = append(functionList, buffer.String())
	}
	return functionList, nil
}
