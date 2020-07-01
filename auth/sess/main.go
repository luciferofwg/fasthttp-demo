package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(1 * time.Minute),
		Path:   "/",
	})
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/pre", preSession)
	r.GET("/do", DoSomethine)

	//防止session会话攻击，每次登陆session+1，下次校验用户时用的是session+1
	r.GET("/incr", incSession)
	r.GET("/", check)
	r.Run(":8080")
}

//添加一条session记录
func preSession(c *gin.Context) {
	session := sessions.Default(c)
	strSession := "test.access.token"
	session.Set("test@mail.com", strSession)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"session": strSession, "result": "sucess"})
}

//获取session
func DoSomethine(c *gin.Context) {
	userEmail := c.Query("user_email")
	if userEmail == "" {
		panic("can not get user email")
	}
	session := sessions.Default(c)
	userAccessToken := session.Get(userEmail)
	if userAccessToken != nil {
		c.JSON(http.StatusOK, gin.H{"session": userAccessToken, "result": "sucess"})
	} else {
		c.JSON(http.StatusOK, gin.H{"session": "", "result": "failed"})
	}

}

func incSession(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	c.JSON(200, gin.H{"count": count})
}

func check(c *gin.Context) {
	val, err := c.Cookie("mysession")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"cookie": "", "result": "failed"})
	} else {
		c.JSON(http.StatusOK, gin.H{"cookie": val, "result": "sucess"})
	}

}
