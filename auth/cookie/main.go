package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	loginUri     = "/sdk_service/rest/users/login/v1.1"
	logoutUri    = "/sdk_service/rest/management/application-logout"
	keepaliveUri = "/sdk_service/rest/management/options"
)

type loginResp struct {
	XMLName        xml.Name `xml:"response"`
	PwdExpireDays  string   `xml:"pwdExpireDays"`
	LeftLockedTime string   `xml:"leftLockedTime"`
	Result         struct {
		XMLName xml.Name `xml:"result"`
		Code    string   `xml:"code"`
		Errmsg  string   `xml:"errmsg"`
	}
}

//登出和保活的返回结果一致
type logoutResp struct {
	XMLName xml.Name `xml:"response"`
	Result  struct {
		XMLName xml.Name `xml:"result"`
		Code    string   `xml:"code"`
		Errmsg  string   `xml:"errmsg"`
	}
}

func main() {
	r := gin.Default()
	r.GET("/token", token)
	r.POST(loginUri, signIn)
	r.POST(logoutUri, signOut)
	r.POST(keepaliveUri, keepLive)
	r.Run("172.20.36.2:9000")
}

func token(c *gin.Context) {
	c.String(http.StatusOK, "%s", "收到的token请求")
}

func signIn(c *gin.Context) {
	fmt.Println("收到登录请求")
	signInResp := &loginResp{}
	signInResp.LeftLockedTime = "10"
	signInResp.PwdExpireDays = "10"
	signInResp.Result.Code = "0"
	signInResp.Result.Errmsg = ""
	c.XML(http.StatusOK, signInResp)
}

func signOut(c *gin.Context) {
	fmt.Println("收到登出包")
	signOutResp := &logoutResp{}
	signOutResp.Result.Code = "0"
	signOutResp.Result.Errmsg = ""
	c.XML(http.StatusOK, signOutResp)
}

func keepLive(c *gin.Context) {
	fmt.Println("收到心跳保活")
	kpResp := &logoutResp{}
	kpResp.Result.Code = "0"
	kpResp.Result.Errmsg = ""
	c.XML(http.StatusOK, kpResp)
}
