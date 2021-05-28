package generator

import (
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/yametech/ym-go-scaffold/pkg/conf"

	"path/filepath"
	"strings"
)

const (
	callTemplateText = `{{.head}}

//go:generate mockgen -destination ./{{.name}}_mock.go -package {{.filePackage}} -source $GOFILE

package {{.filePackage}}

import (
	"context"
	"google.golang.org/grpc"

	{{.package}}

	
)

type (
	{{.alias}}

	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		conn *grpc.ClientConn
	}
)

func New{{.serviceName}}(conn *grpc.ClientConn) {{.serviceName}} {
	return &default{{.serviceName}}{
		conn: conn,
	}
}

{{.functions}}
`

	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}},error)`

	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.serviceName}}) {{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}}, error) {
	client := {{.package}}.New{{.rpcServiceName}}Client(m.conn)
	return client.{{.method}}(ctx, in)
}
`
)

// GenCall generates the rpc client code, which is the entry point for the rpc service call.
// It is a layer of encapsulation for the rpc client and shields the details in the pb.
func (g *DefaultGenerator) GenCall(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetCall()
	service := proto.Service
	head := util.GetHead(proto.Name)

	serviceName := strings.ToLower(proto.Service.Name)

	filename := filepath.Join(dir.Filename, "client.go")
	functions, err := g.genFunction(proto.PbPackage, service)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(service)
	if err != nil {
		return err
	}

	text, err := util.LoadTemplate(category, callTemplateFile, callTemplateText)
	if err != nil {
		return err
	}
	var alias = make([]string, 0)
	for _, item := range proto.Message {

		alias = append(alias, fmt.Sprintf("%s = %s", parser.CamelCase(item.Name), fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(item.Name))))
	}

	err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"name":        serviceName,
		"alias":       strings.Join(alias, util.NL),
		"head":        head,
		"filePackage": dir.Base,
		"package":     fmt.Sprintf(`"%s"`, ctx.GetPb().Package),
		"serviceName": stringx.From(service.Name).ToCamel(),
		"functions":   strings.Join(functions, util.NL),
		"interface":   strings.Join(iFunctions, util.NL),
	}, filename, true)
	return err
}

func (g *DefaultGenerator) genFunction(goPackage string, service parser.Service) ([]string, error) {
	functions := make([]string, 0)
	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, callFunctionTemplateFile, callFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]interface{}{
			"serviceName":    stringx.From(service.Name).ToCamel(),
			"rpcServiceName": parser.CamelCase(service.Name),
			"method":         parser.CamelCase(rpc.Name),
			"package":        goPackage,
			"pbRequest":      parser.CamelCase(rpc.RequestType),
			"pbResponse":     parser.CamelCase(rpc.ReturnsType),
			"hasComment":     len(comment) > 0,
			"comment":        comment,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}
	return functions, nil
}

func (g *DefaultGenerator) getInterfaceFuncs(service parser.Service) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, callInterfaceFunctionTemplateFile, callInterfaceFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]interface{}{
				"hasComment": len(comment) > 0,
				"comment":    comment,
				"method":     parser.CamelCase(rpc.Name),
				"pbRequest":  parser.CamelCase(rpc.RequestType),
				"pbResponse": parser.CamelCase(rpc.ReturnsType),
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}
