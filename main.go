package main

import (
	"flag"
	"fmt"
	"ins/client"
	"ins/handlers"
	"ins/task"
)

var username string
var password string
var user string

func main() {

	flag.StringVar(&username, "u", "", "用户名")
	flag.StringVar(&password, "p", "", "密码")
	flag.StringVar(&user, "user", "", "目标采集用户")
	flag.Parse()

	c := client.NewClient(false, "socks5://127.0.0.1:1086")

	c.Login(username, password)

	task := task.NewTask(c, fmt.Sprintf("https://www.instagram.com/%s/", user), handlers.NewUserPage())

	task.Run()

	select {}
}
