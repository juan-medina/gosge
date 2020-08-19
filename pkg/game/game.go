package game

import (
	engineImp "github.com/juan-medina/gosge/internal/engine"
	"github.com/juan-medina/gosge/pkg/engine"
	"github.com/juan-medina/gosge/pkg/options"
)

func Run(opt options.Options, init engine.InitFunc) error {
	return engineImp.New(opt, init).Run()
}
