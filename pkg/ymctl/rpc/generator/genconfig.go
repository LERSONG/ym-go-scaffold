package generator

import (
	"github.com/LERSONG/ym-go-scaffold/pkg/conf"
	"github.com/LERSONG/ym-go-scaffold/pkg/util"
	"github.com/LERSONG/ym-go-scaffold/pkg/ymctl/rpc/parser"
	"io/ioutil"
	"os"
	"path/filepath"
)

const configTemplate = `package config

type Config struct {
}
`

// GenConfig generates the configuration structure definition file of the rpc service,
// which contains the zrpc.RpcServerConf configuration item by default.
// You can specify the naming style of the target file name through config.Config. For details,
// see https://github.com/tal-tech/go-zero/tree/master/tools/goctl/config/config.go
func (g *DefaultGenerator) GenConfig(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetConfig()

	fileName := filepath.Join(dir.Filename, "config.go")
	if util.FileExists(fileName) {
		return nil
	}

	text, err := util.LoadTemplate(category, configTemplateFileFile, configTemplate)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(text), os.ModePerm)
}
