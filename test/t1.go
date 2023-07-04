package main

import (
	"context"
	"fmt"
	"github.com/PeterYangs/yycmsScript"
	"strconv"
	"time"
)

func main() {

	t := yycmsScript.NewYyCmsScript(context.Background())

	s := t.StartFunc(func(app *yycmsScript.App) (string, error) {

		err := app.StartDefaultServer(func(message string) string {

			switch message {

			case "num":

				return app.Data.Get("num", "")

			case "stop":

				app.Cancel()

				return ""

			default:

				return "no"

			}

		})

		if err != nil {

			return "", err
		}

		num := 0

		for {

			select {

			case <-app.GetCxt().Done():

				return "", nil

			default:

				time.Sleep(1 * time.Second)

				num++

				fmt.Println(num)

				app.Data.Set("num", strconv.Itoa(num))

			}

		}

		//return "", nil
	})

	//s.Flag("file", "文件路径").Required()

	s.SmoothExit()

	c := t.Command("check", "检查", func(app *yycmsScript.App) (string, error) {

		res, err := app.SendDefaultServer("num")

		if err != nil {

			fmt.Println(err)
		}

		fmt.Println("收到消息:", res)

		return "", nil
	})

	c.Flag("file", "文件路径")

	t.Command("stop", "停止", func(app *yycmsScript.App) (string, error) {

		//app.Cancel()

		app.SendDefaultServer("stop")

		return "", nil
	})

	err := t.Run("test", "测试用例")

	if err != nil {

		fmt.Println(err)

	}

}
