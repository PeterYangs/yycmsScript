package yycmsScript

import (
	"errors"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

type YyCmsScript struct {
	version   string
	startFunc func(app *App) (string, error)
	list      []*Item
	app       *kingpin.Application
}

func NewYyCmsScript() *YyCmsScript {

	return &YyCmsScript{version: "0.0.1"}
}

// SetVersion 设置版本号
func (y *YyCmsScript) SetVersion(version string) {

	y.version = version

}

// StartFunc 启动函数
func (y *YyCmsScript) StartFunc(f func(app *App) (string, error)) *Item {

	y.startFunc = func(app *App) (string, error) {

		return f(app)
	}

	i := NewItem("start", y.startFunc, "启动服务.")

	y.list = append(y.list, i)

	return i
}

func (y *YyCmsScript) Command(name string, help string, f func(app *App) (string, error)) *Item {

	i := NewItem(name, f, help)

	y.list = append(y.list, i)

	return i

}

func (y *YyCmsScript) Run(appName string, appDesc string) error {

	app := kingpin.New(appName, appDesc)

	y.app = app

	app.Version(y.version)

	for _, item := range y.list {

		com := app.Command(item.name, item.help)

		if len(item.flags) > 0 {

			for _, flag := range item.flags {

				f := com.Flag(flag.name, flag.name)

				if flag.required {

					f.Required()
				}

				if flag.isBool {

					f.Bool()

				} else {

					f.String()
				}

			}

		}

		if len(item.args) > 0 {

			for _, arg := range item.args {

				com.Arg(arg.name, arg.help)

			}

		}

	}

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch cmd {

	default:

		item, iErr := y.getItemByName(cmd)

		if iErr != nil {

			return iErr
		}

		str, SErr := item.fun(NewApp(y.version, y.getRequest(cmd), appName))

		if SErr != nil {

			return SErr
		}

		print(str)

	}

	return nil

}

func (y *YyCmsScript) getItemByName(name string) (*Item, error) {

	for _, item := range y.list {

		if item.name == name {

			return item, nil
		}

	}

	return nil, errors.New("未找到命令")

}

func (y *YyCmsScript) getRequest(name string) *Request {

	flags := make(map[string]string)
	args := make(map[string]string)

	startItem, _ := y.getItemByName(name)

	//参数绑定
	for i2, flag := range startItem.flags {

		flags[flag.name] = y.app.GetCommand(name).Model().Flags[i2].String()

	}

	for i2, arg := range startItem.args {

		args[arg.name] = y.app.GetCommand(name).Model().Args[i2].String()

	}

	return NewRequest(name, flags, args)

}
