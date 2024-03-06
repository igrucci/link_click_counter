package pkg

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func LogRequest(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		logrus.WithFields(logrus.Fields{
			"method":  string(ctx.Method()),
			"uri":     string(ctx.RequestURI()),
			"headers": string(ctx.Request.Header.String()),
			"body":    string(ctx.Request.Body()),
		}).Info("request received")

		next(ctx)

		logrus.WithFields(logrus.Fields{
			"method":  string(ctx.Method()),
			"uri":     string(ctx.RequestURI()),
			"status":  ctx.Response.StatusCode(),
			"headers": string(ctx.Response.Header.String()),
			"body":    string(ctx.Response.Body()),
		}).Info("request completed")
	}
}
