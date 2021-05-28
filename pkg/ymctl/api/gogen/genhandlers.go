package gogen

import (
	"fmt"
	"path"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const handlerTemplate = `package handler

import (
	"context"
	"github.com/gin-gonic/gin"

	{{.ImportPackages}}

	"net/http"

	
)

func {{.HandlerName}}(ctx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}{{end}}

		l := logic.New{{.LogicType}}(context.Background(), ctx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		} else {
			{{if .HasResp}}c.JSON(http.StatusOK, resp){{else}}c.JSON(http.StatusOK, nil){{end}}
		}
	}
}
`

type handlerInfo struct {
	ImportPackages string
	HandlerName    string
	RequestType    string
	LogicType      string
	Call           string
	HasResp        bool
	HasRequest     bool
}

func genHandler(dir string, group spec.Group, route spec.Route) error {
	handler := getHandlerName(route)
	if getHandlerFolderPath(group, route) != handlerDir {
		handler = strings.Title(handler)
	}
	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	return doGenToFile(dir, handler, group, route, handlerInfo{
		ImportPackages: genHandlerImports(group, route, parentPkg),
		HandlerName:    handler,
		RequestType:    util.Title(route.RequestTypeName()),
		LogicType:      strings.Title(getLogicName(route)),
		Call:           strings.Title(strings.TrimSuffix(handler, "Handler")),
		HasResp:        len(route.ResponseTypeName()) > 0,
		HasRequest:     len(route.RequestTypeName()) > 0,
	})
}

func doGenToFile(dir, handler string, group spec.Group,
	route spec.Route, handleObj handlerInfo) error {

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getHandlerFolderPath(group, route),
		filename:        strings.ToLower(handler) + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data:            handleObj,
	})
}

func genHandlers(dir string, api *spec.ApiSpec) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, group, route); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"",
		util.JoinPackages(parentPkg, getLogicFolderPath(group, route))))
	imports = append(imports, fmt.Sprintf("\"%s\"", util.JoinPackages(parentPkg, contextDir)))
	if len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", util.JoinPackages(parentPkg, typesDir)))
	}

	return strings.Join(imports, "\n\t")
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")
	return handler, nil
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return handlerDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(handlerDir, folder)
}

func getHandlerName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Handler"
}

func getLogicName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Logic"
}
