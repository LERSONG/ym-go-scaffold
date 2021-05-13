package generator

import (
	"path/filepath"
	"strings"

	"github.com/LERSONG/ym-go-scaffold/pkg/util"
	"github.com/LERSONG/ym-go-scaffold/pkg/util/ctx"
	"github.com/LERSONG/ym-go-scaffold/pkg/util/stringx"
	"github.com/LERSONG/ym-go-scaffold/pkg/ymctl/rpc/parser"
)

const (
	wd     = "wd"
	cmd    = "cmd"
	pkg    = "pkg"
	config = "conf"
	logic  = "logic"
	server = "server"
	svc    = "svc"
	//pt    = "proto"
	pb   = "pb"
	call = "call"
)

type (
	// DirContext defines a rpc service directories context
	DirContext interface {
		GetWd() Dir
		GetCall() Dir
		GetPkg() Dir
		GetConfig() Dir
		GetLogic() Dir
		GetServer() Dir
		GetSvc() Dir
		GetPb() Dir
		//GetProto() Dir
		GetMain() Dir
		GetServiceName() stringx.String
	}

	// Dir defines a directory
	Dir struct {
		Base     string
		Filename string
		Package  string
	}

	defaultDirContext struct {
		inner       map[string]Dir
		serviceName stringx.String
	}
)

func mkdir(ctx *ctx.ProjectContext, proto parser.Proto) (DirContext, error) {
	inner := make(map[string]Dir)
	pkgDir := filepath.Join(ctx.WorkDir, "pkg")
	cmdDir := filepath.Join(ctx.WorkDir, "cmd")
	configDir := filepath.Join(pkgDir, "conf")
	logicDir := filepath.Join(pkgDir, "logic")
	serverDir := filepath.Join(pkgDir, "server")
	svcDir := filepath.Join(pkgDir, "svc")
	//protoDir := filepath.Join(ctx.WorkDir, "proto")
	pbDir := filepath.Join(ctx.WorkDir, proto.GoPackage)
	callDir := filepath.Join(ctx.WorkDir, "client")
	if strings.ToLower(proto.Service.Name) == strings.ToLower(proto.GoPackage) {
		callDir = filepath.Join(ctx.WorkDir, strings.ToLower(stringx.From(proto.Service.Name+"_client").ToCamel()))
	}

	inner[wd] = Dir{
		Filename: ctx.WorkDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(ctx.WorkDir, ctx.Dir))),
		Base:     filepath.Base(ctx.WorkDir),
	}
	inner[cmd] = Dir{
		Filename: cmdDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(ctx.WorkDir, ctx.Dir))),
		Base:     filepath.Base(cmdDir),
	}
	inner[pkg] = Dir{
		Filename: pkgDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(pkgDir, ctx.Dir))),
		Base:     filepath.Base(pkgDir),
	}
	inner[config] = Dir{
		Filename: configDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(configDir, ctx.Dir))),
		Base:     filepath.Base(configDir),
	}
	inner[logic] = Dir{
		Filename: logicDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(logicDir, ctx.Dir))),
		Base:     filepath.Base(logicDir),
	}
	inner[server] = Dir{
		Filename: serverDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(serverDir, ctx.Dir))),
		Base:     filepath.Base(serverDir),
	}
	inner[svc] = Dir{
		Filename: svcDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(svcDir, ctx.Dir))),
		Base:     filepath.Base(svcDir),
	}
	//inner[pt] = Dir{
	//	Filename: protoDir,
	//	Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(pbDir, ctx.Dir))),
	//	Base:     filepath.Base(protoDir),
	//}
	inner[pb] = Dir{
		Filename: pbDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(pbDir, ctx.Dir))),
		Base:     filepath.Base(pbDir),
	}
	inner[call] = Dir{
		Filename: callDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(callDir, ctx.Dir))),
		Base:     filepath.Base(callDir),
	}
	for _, v := range inner {
		err := util.MkdirIfNotExist(v.Filename)
		if err != nil {
			return nil, err
		}
	}
	serviceName := strings.TrimSuffix(proto.Name, filepath.Ext(proto.Name))
	return &defaultDirContext{
		inner:       inner,
		serviceName: stringx.From(strings.ReplaceAll(serviceName, "-", "")),
	}, nil
}

func (d *defaultDirContext) GetWd() Dir {
	return d.inner[wd]
}

func (d *defaultDirContext) GetCall() Dir {
	return d.inner[call]
}

func (d *defaultDirContext) GetPkg() Dir {
	return d.inner[pkg]
}

func (d *defaultDirContext) GetConfig() Dir {
	return d.inner[config]
}

func (d *defaultDirContext) GetLogic() Dir {
	return d.inner[logic]
}

func (d *defaultDirContext) GetServer() Dir {
	return d.inner[server]
}

func (d *defaultDirContext) GetSvc() Dir {
	return d.inner[svc]
}

func (d *defaultDirContext) GetPb() Dir {
	return d.inner[pb]
}

//func (d *defaultDirContext) GetProto() Dir {
//	return d.inner[pt]
//}

func (d *defaultDirContext) GetMain() Dir {
	return d.inner[cmd]
}

func (d *defaultDirContext) GetServiceName() stringx.String {
	return d.serviceName
}

// Valid returns true if the directory is valid
func (d *Dir) Valid() bool {
	return len(d.Filename) > 0 && len(d.Package) > 0
}
