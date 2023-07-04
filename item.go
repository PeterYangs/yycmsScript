package yycmsScript

type Item struct {
	fun        func(app *App) (string, error)
	flags      []*Flag
	args       []*Arg
	name       string
	help       string
	hide       bool
	smoothExit bool
}

func NewItem(name string, fun func(app *App) (string, error), help string) *Item {

	return &Item{fun: fun, help: help, flags: []*Flag{}, args: []*Arg{}, name: name}
}

func (i *Item) Flag(name string, help string) *Flag {

	f := NewFlag(name, help, i)

	i.flags = append(i.flags, f)

	return f

}

func (i *Item) Arg(name string, help string) *Arg {

	a := NewArg(name, help, i)

	i.args = append(i.args, a)

	return a
}

func (i *Item) Hide() *Item {

	i.hide = true

	return i
}

// SmoothExit 打开平滑退出
func (i *Item) SmoothExit() *Item {

	i.smoothExit = true

	return i

}
