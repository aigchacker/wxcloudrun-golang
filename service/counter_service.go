package service

import (
	"fmt"

	"log"

	"github.com/cutesdk/cutesdk-go/wxmp"
	"github.com/idoubi/goutils"
	"github.com/labstack/echo/v4"
)

const (
	appid            = "wx9eaf473f4a8eca98"               // wxmp appid
	secret           = "ae5bf47ebb26bcb3903fa73065f91b3b" // wxmp secret
	token            = "7VuyexreqaavPWgLHiemvBpiXme"
	encodingAESKey   = "RwCnIeQvwkFFk0deJnsrA3WgZGQE2atglBZm4vw1bjQ"
	subscribeMessage = "您好，欢迎关注瞬息照相馆。\n\n马上开始为你生成风格百变的AI写真照哦！"
)

func main() {
	e := echo.New()

	e.POST("/user-login-qrcode", GetLoginQrcode)
	e.POST("/check-login/:code", CheckLogin)
	e.Any("/wxmp/notify", WxmpNotify)

	e.Logger.Fatal(e.Start(":443"))
}

func GetLoginQrcode(c echo.Context) error {
	// login scene
	scene := "mplogin"
	// expires: 5 minutes
	var expires int64 = 300
	qrcode, err := getSceneQrcode(scene, expires)
	if err != nil {
		log.Printf("get mplogin qrcode failed: %v\n", err)
		return c.String(500, "get login qrcode failed")
	}

	loginCode := goutils.MD5(qrcode.Ticket)
	cacheKey := fmt.Sprintf("lgt:%s", loginCode)

	// todo setex cacheKey 0 expires

	if err != nil {
		log.Printf("get login qrcode falied: %v\n", err)
		return c.String(500, "get login qrcode failed")
	}

	return c.JSON(200, map[string]interface{}{
		"login_code":  loginCode,
		"login_qrurl": qrcode.Qrurl,
		"expires":     expires,
		"cache_key":   cacheKey,
	})
}

func CheckLogin(c echo.Context) error {
	loginCode := c.Param("code")

	cacheKey := fmt.Sprintf("lgt:%s", loginCode)

	// todo get user openid from cache

	jwtToken := "xxx"

	return c.JSON(200, map[string]interface{}{
		"token":     jwtToken,
		"cache_key": cacheKey,
	})
}

func WxmpNotify(c echo.Context) error {
	// 1. get wxmp notify EventKey(including login_code)
	// 2. get user openid
	// 3. set user openid to cacheKey(same with cacheKey in CheckLogin)

	return nil
}

// SceneQrcode struct
type SceneQrcode struct {
	Scene   string `json:"scene"`
	Expires int64  `json:"expires"`
	Ticket  string `json:"ticket"`
	Qrurl   string `json:"qrurl"`
}

func getSceneQrcode(scene string, expires int64) (*SceneQrcode, error) {
	cli, _ := getWxmpClient()

	uri := "/cgi-bin/qrcode/create"

	params := map[string]interface{}{
		"expire_seconds": expires,
		"action_name":    "QR_STR_SCENE",
		"action_info": map[string]interface{}{
			"scene": map[string]interface{}{
				"scene_str": scene,
			},
		},
	}

	res, err := cli.PostWithToken(uri, params)
	if err != nil {
		log.Printf("get wxmp qrcode failed: %v\n", err)
		return nil, err
	}

	ticket := res.GetString("ticket")
	if ticket == "" {
		log.Printf("get wxmp qrcode failed: %s\n", res.String())
		return nil, fmt.Errorf("get wxmp qrcode failed: %s", res.String())
	}

	qrurl := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", ticket)

	return &SceneQrcode{
		Scene:   scene,
		Expires: expires,
		Ticket:  ticket,
		Qrurl:   qrurl,
	}, nil
}

func getWxmpClient() (*wxmp.Client, error) {
	opts := &wxmp.Options{
		Appid:  appid,
		Secret: secret,
		Debug:  true,
	}

	return wxmp.NewClient(opts)
}
