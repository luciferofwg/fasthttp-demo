package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"net/http"
)

func main() {
	addr := ":8080"
	certFile := "../key/server.crt"
	keyFile := "../key/server.key"

	router := fasthttprouter.New()
	router.GET("/", getting)
	srv := &fasthttp.Server{
		Handler: router.Handler,
	}
	if err := srv.ListenAndServeTLS(addr, certFile, keyFile); err != nil {
		panic(err)
	}
}

func getting(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("sucess")
}
