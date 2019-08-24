package promptui

import (
	"io"
	"os"
)

type Ctx struct {
	stdout io.WriteCloser
}

type component interface {
	Render(io.Writer)
}

func New() *Ctx {
	return &Ctx{
		stdout: os.Stdout,
	}
}

func (ctx *Ctx) Progress() *ProgressGroup { return nil }
func (ctx *Ctx) Spinner() *SpinGroup      { return nil }
func (ctx *Ctx) SpinGroup() *SpinGroup    { return nil }
func (ctx *Ctx) InFrame() *Frame          { return nil }
func (ctx *Ctx) InColorFrame() *Frame     { return nil }
func (ctx *Ctx) Fmt()                     {}
func (ctx *Ctx) Ask()                     {}
func (ctx *Ctx) Select()                  {}
func (ctx *Ctx) Confirm()                 {}
func (ctx *Ctx) Password()                {}
