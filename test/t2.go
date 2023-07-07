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

	t.StartFunc(func(app *yycmsScript.App) (string, error) {

		err := app.StarCustomServer("t2", func(message string) string {

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

	t.Command("num", "检查", func(app *yycmsScript.App) (string, error) {

		res, err := app.SendCustomServer("t2", "num")

		if err != nil {

			fmt.Println(err)
		}

		fmt.Println("收到消息:", res)

		return "", nil
	})

	t.Command("stop", "停止", func(app *yycmsScript.App) (string, error) {

		fmt.Println(app.SendCustomServer("t2", "stop"))

		return "", nil
	})

	err := t.Run("test2", "测试用例")

	if err != nil {

		fmt.Println(err)

	}

}
