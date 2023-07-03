package yycmsScript

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type App struct {
	version string
	Request *Request
	appName string
}

func NewApp(version string, request *Request, appName string) *App {

	return &App{version: version, Request: request, appName: appName}
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

				//fmt.Println("收到信息：", message)
				//
				//u.Write([]byte("nice\n"))

				fmt.Println(call(message+"\n"), "kkk")

				_, eee := u.Write([]byte(call(message + "\n")))

				fmt.Println(eee, "uuuu")

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
		//conn.Close()

	}()

	conn.Write([]byte(message + "\n"))

	reader := bufio.NewReader(conn)

	res, rErr := reader.ReadString('\n')

	if rErr != nil {

		fmt.Println(rErr)

		return "", rErr
	}

	fmt.Println(res, "jjjj")

	return res, nil
}
