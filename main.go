package main

import (
	"fmt"
	"go-serv/app"
)

func main() {
	a := app.NewApp()
	err := a.Run()
	if err != nil {
		fmt.Println("服务 run 错误:", err)
	} else {
		fmt.Println("服务 run 起来了")
	}
}
