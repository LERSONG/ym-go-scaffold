package gogen

import (
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	contextFilename = "servicecontext"
	contextTemplate = `package svc

import (
	{{.configImport}}
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
`
)

func genServiceContext(dir string, api *spec.ApiSpec) error {

	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}
	var configImport = "\"" + ctlutil.JoinPackages(parentPkg, configDir) + "\""

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          contextDir,
		filename:        contextFilename + ".go",
		templateName:    "contextTemplate",
		category:        category,
		templateFile:    contextTemplateFile,
		builtinTemplate: contextTemplate,
		data: map[string]string{
			"configImport": configImport,
		},
	})
}
