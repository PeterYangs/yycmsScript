package main

import (
	"fmt"
	"yycmsScript"
)

func main() {

	t := yycmsScript.NewYyCmsScript()

	t.StartFunc(func(app *yycmsScript.App) (string, error) {

		//fmt.Println(app.Request.GetFlag("file"), "----")

		err := app.StartDefaultServer(func(message string) string {

			fmt.Println(message, "xxx")

			switch message {

			case "yes":

				return "yes啊"

			default:

				return "no"

			}

		})

		if err != nil {

			fmt.Println(err)
		}

		for {

			select {}
		}

		return "", nil
	})

	//s.Flag("file", "文件路径").Required()

	c := t.Command("check", "检查", func(app *yycmsScript.App) (string, error) {

		res, err := app.SendDefaultServer("yes")

		if err != nil {

			fmt.Println(err)
		}

		fmt.Println("收到消息:", res)

		return "", nil
	})

	c.Flag("file", "文件路径")

	t.Run("test", "测试用例")
}
