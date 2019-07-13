package main

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	CONTENT_VIID_JSON string = "application/VIID+json; charset=utf-8"
	CONTENT_PLAIN     string = "text/plain"
	CONTENT_APP_JSON  string = "application/json"
	CONTENT_JS        string = "application/javascript"
	CONTENT_APP_XML   string = "application/xml"
	CONTENT_TEXT_XML  string = "text/xml"
	CONTENT_TEXT_HTML string = "text/html"
	CONTENT_FORMDATA  string = "multipart/form-data"
	CONTENT_FORM      string = "application/x-www-form-urlencoded"
)

const (
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
)

const (
	URI_POST_MULTIFILE  = "http://localhost:20000/postmultifile"
	URI_POST_SINGLEFILE = "http://localhost:20000/postfile"
	URI_POST            = "http://localhost:20000/post?id=hahaha"
	URI_GET             = "http://localhost:20000/get?username=admin&id=123"
)

var (
	c = fasthttp.Client{}
)

func get() {
	//请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetContentType(CONTENT_PLAIN)
	req.Header.SetMethod(METHOD_GET)
	req.SetRequestURI(URI_GET)

	//回复
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.DoTimeout(req, resp, time.Duration(time.Millisecond*50)); err != nil {
		fmt.Printf("发送数据失败，错误：%v", err)
	}
}

func postForm() {
	//请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetContentType(CONTENT_FORM)
	req.Header.SetMethod(METHOD_POST)

	req.SetRequestURI(URI_POST)

	//POST参数
	args := &fasthttp.Args{}
	args.Add("username", "admin")
	args.Add("password", "456")
	args.WriteTo(req.BodyWriter())
	//回复
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.DoTimeout(req, resp, time.Duration(time.Millisecond*50)); err != nil {
		fmt.Printf("发送数据失败，错误：%v", err)
	}
}

func postFile() {
	//上传的文件
	filename := `D:/0_XLServers/SeemmoSPJGHServer/MotorVehicle.txt`
	//创建缓冲区，用于存放文件内容
	bodyBuffer := &bytes.Buffer{}
	//创建一个multipart文件写入器，方便按照http规定格式写入内容
	bodyWrite := multipart.NewWriter(bodyBuffer)
	//从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
	f := path.Base(filename)
	fmt.Printf("=======%v\n", f)
	fileWriter, err := bodyWrite.CreateFormFile("uploadfile", f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//打开上传文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Printf("copy file filed,%v", err)
		return
	}
	//关闭bodyWriter停止写入数据
	bodyWrite.Close()
	contentType := bodyWrite.FormDataContentType()

	//请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetContentType(contentType)
	req.Header.SetMethod(METHOD_POST)
	req.SetRequestURI(URI_POST_SINGLEFILE)
	req.SetBody(bodyBuffer.Bytes())

	//回复
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.DoTimeout(req, resp, time.Duration(time.Millisecond*2000)); err != nil {
		fmt.Printf("发送数据失败，错误：%v", err)
		return
	}
	fmt.Printf("收到回复信息：%v\n", string(resp.Body()))
}

func postMultiFile() {
	//上传的文件
	fatherDir := "C:/test"
	var filenames = []string{}
	files, _ := ioutil.ReadDir(fatherDir)
	for _, f := range files {
		filenames = append(filenames, filepath.Join(fatherDir, f.Name()))
	}

	//创建缓冲区，用于存放文件内容
	bodyBuffer := &bytes.Buffer{}
	//创建一个multipart文件写入器，方便按照http规定格式写入内容
	bodyWrite := multipart.NewWriter(bodyBuffer)
	//从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次

	for _, filename := range filenames {
		f := path.Base(strings.Replace(filename, "\\", "/", -1))
		fmt.Printf("filename:%v\n", f)
		fileWriter, err := bodyWrite.CreateFormFile("files", f)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//打开上传文件
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			fmt.Printf("copy file filed,%v", err)
			return
		}
	}

	//关闭bodyWriter停止写入数据
	bodyWrite.Close()
	contentType := bodyWrite.FormDataContentType()

	//请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetContentType(contentType)
	req.Header.SetMethod(METHOD_POST)
	req.SetRequestURI(URI_POST_MULTIFILE)
	req.SetBody(bodyBuffer.Bytes())

	//回复
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.DoTimeout(req, resp, time.Duration(time.Millisecond*2000)); err != nil {
		fmt.Printf("发送数据失败，错误：%v", err)
		return
	}
	fmt.Printf("收到回复信息：%v\n", string(resp.Body()))

}

func main() {

	get()
	postForm()
	postFile()
	postMultiFile()

}
