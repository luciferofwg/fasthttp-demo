package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"os"
)

func main() {
	router := fasthttprouter.New()
	//GET 方法获取参数
	router.GET("/get", getting)
	//POST 方法获取参数&body
	router.POST("/post", posting)
	//post处理单个文件
	router.POST("/postfile", postSinglefile)
	//post处理多个文件
	router.POST("/postmultifile", postMultifile)
	srv := &fasthttp.Server{
		Handler: router.Handler,
	}
	addr := "localhost:20000"
	fmt.Printf("server listen %s\n", addr)
	if err := srv.ListenAndServe(addr); err != nil {
		fmt.Printf("listen failed,err:%v", err)
		return
	}
}

func getting(ctx *fasthttp.RequestCtx) {
	fmt.Printf("收到Get请求：\n")
	uri := ctx.Request.URI().String()
	method := ctx.Method()
	qureyArgs := ctx.QueryArgs()
	username := qureyArgs.Peek("username")
	id := qureyArgs.Peek("id")

	fmt.Printf("uri:%v \nmethod:%v \nusername:%v\nid:%v\n", uri, string(method), string(username), string(id))

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("sucess")
}

func posting(ctx *fasthttp.RequestCtx) {
	fmt.Printf("收到Post请求：\n")
	uri := ctx.Request.URI().String()
	method := ctx.Method()

	qureyArgs := ctx.QueryArgs()
	id := qureyArgs.Peek("id")

	postArgs := ctx.PostArgs()
	username := postArgs.Peek("username")
	password := postArgs.Peek("password")

	fmt.Printf("uri:%v, method:%v\n", uri, string(method))
	fmt.Printf("queryArgs:id:%v\n", string(id))
	fmt.Printf("postArgs:username:%v, password:%v\n", string(username), string(password))

}

func postSinglefile(ctx *fasthttp.RequestCtx) {
	fmt.Printf("收到Post单个文件长传：\n")
	uri := ctx.Request.URI().String()
	method := ctx.Method()
	contentType := ctx.Request.Header.ContentType()

	fmt.Printf("uri:%v, method:%v, contectTyp:%v\n", uri, string(method), string(contentType))
	fileHeader, err := ctx.FormFile("uploadfile")
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("message is error")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(err.Error())
		return
	}
	defer file.Close()

	//
	fmt.Printf("filename:%v\n", fileHeader.Filename)
	savePath := "C:\\tmp\\" + fileHeader.Filename
	localFile, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		ctx.WriteString(err.Error())
		return
	}
	defer localFile.Close()

	//写文件
	_, err = io.Copy(localFile, file)
	if err != nil {
		ctx.WriteString(err.Error())
		return
	}
	ctx.WriteString("sucess")
}

func postMultifile(ctx *fasthttp.RequestCtx) {
	fmt.Printf("收到Post多个文件长传：\n")
	uri := ctx.Request.URI().String()
	method := ctx.Method()
	contentType := ctx.Request.Header.ContentType()

	fmt.Printf("uri:%v, method:%v, contectTyp:%v\n", uri, string(method), string(contentType))

	multiForm, err := ctx.MultipartForm()
	if err != nil {
		ctx.WriteString(err.Error())
		return
	}
	fileHeaders := multiForm.File["files"]
	for key, fileHeader := range fileHeaders {
		_ = key
		filename := fileHeader.Filename
		size := fileHeader.Size
		fd, err := fileHeader.Open()
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		defer fd.Close()
		fmt.Printf("filename:%v, size:%v\n", filename, size)

		savePath := "C:\\tmp\\" + fileHeader.Filename
		localFile, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		defer localFile.Close()

		//写文件
		_, err = io.Copy(localFile, fd)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
	}
	ctx.WriteString("sucess")
}
