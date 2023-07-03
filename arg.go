package yycmsScript

type Arg struct {
	name     string
	value    string
	help     string
	required bool
	item     *Item
}

func NewArg(name string, help string, item *Item) *Arg {

	return &Arg{name: name, help: help, item: item}
}

func (arg *Arg) Required() *Arg {

	arg.required = true

	return arg

}
