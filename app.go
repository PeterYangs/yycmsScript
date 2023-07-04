package yycmsScript

import (
	"bufio"
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
}

func NewApp(version string, request *Request, appName string) *App {

	return &App{version: version, Request: request, appName: appName, Data: NewData()}
}

func (app *App) accept(unixListener *net.UnixListener, callback func(message string) string) {

	defer unixListener.Close()

	for {

		unixConn, aErr := unixListener.AcceptUnix()

		if aErr != nil {

			fmt.Println(aErr)

			continue
		}

		go func(u *net.UnixConn, call func(message string) string) {

			defer u.Close()

			reader := bufio.NewReader(u)

			for {

				message, rErr := reader.ReadString('\n')

				if rErr != nil {

					fmt.Println(rErr)

					return
				}

				message = strings.Replace(message, "\n", "", -1)

				u.Write([]byte(call(message) + "\n"))

			}

		}(unixConn, callback)

	}

}

func (app *App) StartDefaultServer(callback func(message string) string) error {

	sockPath := "storage/app/public/" + app.appName + ".sock"

	os.Remove(sockPath)

	unixAddr, err := net.ResolveUnixAddr("unix", sockPath)

	if err != nil {

		return err
	}

	unixListener, lErr := net.ListenUnix("unix", unixAddr)

	if lErr != nil {

		return lErr
	}

	go app.accept(unixListener, callback)

	return nil

}

func (app *App) SendDefaultServer(message string) (string, error) {

	sockPath := "storage/app/public/" + app.appName + ".sock"

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

	return res, nil
}
