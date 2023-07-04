package main

import (
	"fmt"
	"time"
	"yycmsScript"
)

func main() {

	t := yycmsScript.NewYyCmsScript()

	t.StartFunc(func(app *yycmsScript.App) (string, error) {

		err := app.StartDefaultServer(func(message string) string {

			switch message {

			case "num":

				return "yes啊"

			default:

				return "no"

			}

		})

		if err != nil {

			fmt.Println(err)
		}

		time.Sleep(100000 * time.Second)

		return "", nil
	})

	//s.Flag("file", "文件路径").Required()

	c := t.Command("check", "检查", func(app *yycmsScript.App) (string, error) {

		res, err := app.SendDefaultServer("num")

		if err != nil {

			fmt.Println(err)
		}

		fmt.Println("收到消息:", res)

		return "", nil
	})

	c.Flag("file", "文件路径")

	t.Run("test", "测试用例")
}
