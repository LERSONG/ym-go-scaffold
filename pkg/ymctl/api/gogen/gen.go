package gogen

import (
	"errors"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
	"github.com/yametech/ym-go-scaffold/pkg/util"
)

// GoCommand gen go project files from command line
func GoCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")

	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	return DoGenProject(apiFile, dir)
}

func DoGenProject(apiFile, dir string) error {
	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	if err = util.MkdirIfNotExist(dir); err != nil {
		return err
	}

	if err = genMain(dir, api); err != nil {
		return err
	}

	if err = genConfig(dir, api); err != nil {
		return err
	}

	if err = genServer(dir, api); err != nil {
		return err
	}

	if err = genServiceContext(dir, api); err != nil {
		return err
	}

	if err = genTypes(dir, api); err != nil {
		return err
	}

	if err = genRoutes(dir, api); err != nil {
		return err
	}

	if err = genHandlers(dir, api); err != nil {
		return err
	}

	if err = genLogic(dir, api); err != nil {
		return err
	}

	return nil
}
