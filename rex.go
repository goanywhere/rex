package rex

import (
	"path"

	"github.com/goanywhere/env"
	"github.com/goanywhere/fs"
)

func init() {
	var basedir = fs.Getcd(2)
	env.Set("basedir", basedir)
	env.Load(path.Join(basedir, ".env"))
}
