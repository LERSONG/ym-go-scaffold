package generator

import (
	"github.com/yametech/ym-go-scaffold/pkg/conf"
	"github.com/yametech/ym-go-scaffold/pkg/util"
	"github.com/yametech/ym-go-scaffold/pkg/ymctl/rpc/parser"
	"io/ioutil"
	"os"
	"path/filepath"
)

const configTemplate = `package config

type Config struct {
	Addr string
	RegistryAddr string
}

func NewConfig(opts ...Option) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

type Option func(c *Config)

func Addr(addr string) Option {
	return func(c *Config) {
		c.Addr = addr
	}
}

func RegistryAddr(registryAddr string) Option {
	return func(c *Config) {
		c.RegistryAddr = registryAddr
	}
}
`

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
