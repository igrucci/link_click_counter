package main

import (
	"context"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"link_click_counter/config"
	"link_click_counter/internal"
	"link_click_counter/pkg"
	"log"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("Unable to read config.yaml %v", err)
	}

	connPool, err := config.ConnectDB(cfg)
	if err != nil {
		logrus.Fatalf("Unable to connect to database: %v", err)
		return
	}
	if err := connPool.Ping(context.Background()); err != nil {
		logrus.Fatalf("Unable to ping database: %v", err)
		return
	}
	defer connPool.Close()

	router := fasthttprouter.New()

	router.POST("/counters", pkg.LogRequest(func(ctx *fasthttp.RequestCtx) {
		internal.CreateCounter(ctx, connPool)
	}))
	router.GET("/counters", pkg.LogRequest(func(ctx *fasthttp.RequestCtx) {
		internal.GetCounters(ctx, connPool)
	}))
	router.GET("/counters/:code", pkg.LogRequest(func(ctx *fasthttp.RequestCtx) {
		internal.RedirectURL(ctx, connPool)
	}))
	router.GET("/counters/:code/stats", pkg.LogRequest(func(ctx *fasthttp.RequestCtx) {
		internal.GetCounterStats(ctx, connPool)
	}))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
