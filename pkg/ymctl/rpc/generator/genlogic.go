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

const (
	logicTemplate = `package logic

import (
	"context"

	{{.imports}}

)

type {{.logicName}} struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{.logicName}}(ctx context.Context,svcCtx *svc.ServiceContext) *{{.logicName}} {
	return &{{.logicName}}{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
{{.functions}}
`
	logicFunctionTemplate = `{{if .hasComment}}{{.comment}}{{end}}
func (l *{{.logicName}}) {{.method}} (in {{.request}}) ({{.response}}, error) {
	// todo: add your logic here and delete this line
	
	return &{{.responseType}}{}, nil
}
`
)

// GenLogic generates the logic file of the rpc service, which corresponds to the RPC definition items in proto.
func (g *DefaultGenerator) GenLogic(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetLogic()
	for _, rpc := range proto.Service.RPC {
		logicFilename := strings.ToLower(rpc.Name) + "logic"

		filename := filepath.Join(dir.Filename, logicFilename+".go")
		functions, err := g.genLogicFunction(proto.PbPackage, rpc)
		if err != nil {
			return err
		}

		imports := make([]string, 0)
		imports = append(imports, fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
		imports = append(imports, fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
		text, err := util.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
		if err != nil {
			return err
		}
		err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"logicName": fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel()),
			"functions": functions,
			"imports":   strings.Join(imports, util.NL),
		}, filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *DefaultGenerator) genLogicFunction(goPackage string, rpc *parser.RPC) (string, error) {
	var functions = make([]string, 0)
	text, err := util.LoadTemplate(category, logicFuncTemplateFileFile, logicFunctionTemplate)
	if err != nil {
		return "", err
	}

	logicName := stringx.From(rpc.Name + "_logic").ToCamel()
	comment := parser.GetComment(rpc.Doc())
	buffer, err := util.With("fun").Parse(text).Execute(map[string]interface{}{
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"request":      fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
		"response":     fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"responseType": fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"hasComment":   len(comment) > 0,
		"comment":      comment,
	})
	if err != nil {
		return "", err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, util.NL), nil
}
