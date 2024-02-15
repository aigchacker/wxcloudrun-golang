package main

import (
	"fmt"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"

	"github.com/labstack/echo/v4"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	e := echo.New()

	e.POST("/user-login-qrcode", service.GetLoginQrcode)
	e.POST("/check-login/:code", service.CheckLogin)
	e.Any("/wxmp/notify", service.WxmpNotify)

	e.Logger.Fatal(e.Start(":443"))

}
