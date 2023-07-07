package yycmsScript

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type App struct {
	version string
	Request *Request
	appName string
	Data    *Data
	yy      *YyCmsScript
}

func NewApp(yy *YyCmsScript, version string, request *Request, appName string) *App {

	return &App{version: version, Request: request, appName: appName, Data: NewData(), yy: yy}
}

func (app *App) accept(unixListener *net.UnixListener, callback func(message string) string) {

	defer func() {

		unixListener.Close()

		app.yy.wait.Done()

	}()

	go func(u *net.UnixListener) {

		select {

		case <-app.yy.cxt.Done():

			u.Close()

		}

	}(unixListener)

	for {

		select {

		case <-app.yy.cxt.Done():

			return

		default:

		}

		unixConn, aErr := unixListener.AcceptUnix()

		if aErr != nil {

			select {

			case <-app.yy.cxt.Done():

				return

			default:

			}

			fmt.Println(aErr)

			continue
		}

		go func(u *net.UnixConn, call func(message string) string) {

			defer u.Close()

			reader := bufio.NewReader(u)

			for {

				message, rErr := reader.ReadString('\n')

				if rErr != nil {

					select {

					case <-app.yy.cxt.Done():

						return

					default:

					}

					fmt.Println(rErr)

					return
				}

				message = strings.Replace(message, "\n", "", -1)

				switch message {

				case "__ping":

					u.Write([]byte("__ping" + "\n"))

					break

				default:

					u.Write([]byte(call(message) + "\n"))
				}

			}

		}(unixConn, callback)

	}

}

func (app *App) StartDefaultServer(callback func(message string) string) error {

	return app.StarCustomServer(app.appName, callback)

}

func (app *App) StarCustomServer(sockName string, callback func(message string) string) error {

	sockPath := "storage/app/public/" + sockName + ".sock"

	ok, _ := PathExists(sockPath)

	if ok {

		m, pErr := app.SendCustomServer(sockName, "__ping")

		if pErr != nil {

			os.Remove(sockPath)

		}

		if m == "__ping" {

			return errors.New("sock文件已存在：" + sockPath)

		}

	}

	unixAddr, err := net.ResolveUnixAddr("unix", sockPath)

	if err != nil {

		return err
	}

	unixListener, lErr := net.ListenUnix("unix", unixAddr)

	if lErr != nil {

		return lErr
	}

	app.yy.wait.Add(1)

	go app.accept(unixListener, callback)

	return nil

}

func (app *App) SendDefaultServer(message string) (string, error) {

	return app.SendCustomServer(app.appName, message)
}

func (app *App) SendCustomServer(sockName string, message string) (string, error) {

	sockPath := "storage/app/public/" + sockName + ".sock"

	unixAddr, err := net.ResolveUnixAddr("unix", sockPath)

	if err != nil {

		return "", err
	}

	conn, cErr := net.DialUnix("unix", nil, unixAddr)

	if cErr != nil {

		return "", cErr
	}

	defer func() {

		conn.Close()

	}()

	conn.Write([]byte(message + "\n"))

	reader := bufio.NewReader(conn)

	res, rErr := reader.ReadString('\n')

	if rErr != nil {

		fmt.Println(rErr)

		return "", rErr
	}

	res = strings.Replace(res, "\n", "", -1)

	return res, nil

}

func (app *App) GetCxt() context.Context {

	return app.yy.GetCxt()
}

func (app *App) Cancel() {

	app.yy.cancel()
}
