package app

import (
	common_log "common/log"
	"fmt"
	"go-hj-hospital/HospitalQueue/router"
	"go-hj-hospital/config"
	"go-hj-hospital/util"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

var app *iris.Application

func IrisInit() {
	app = iris.New()

	app.Use(recover.New())

	app.Logger().SetLevel("info")
	irisLogConfig := logger.DefaultConfig()
	irisLogConfig.LogFuncCtx = irisLogFunc
	app.Use(logger.New(irisLogConfig))
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost},
		AllowedHeaders:   []string{"*"},
	}))

	app.AllowMethods(iris.MethodOptions)

	app.Use(func(ctx *context.Context) {
		ctx.Next()
	})

	router.Routes(app)
}

func IrisStart() {
	listener, err := net.Listen("tcp4", config.Instance.Host)
	if err != nil {
		common_log.Error(err)
		os.Exit(1)
	}

	handleSignal(listener)
	if err := app.Run(iris.Listener(listener), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02",
		Charset:                           "UTF-8",
	})); err != nil {
		common_log.Error(err)
		os.Exit(1)
	}
}

func handleSignal(server net.Listener) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		fmt.Printf("got signal [%s], exiting now", s)
		util.CloseMasterDB()
		if err := server.Close(); nil != err {
			fmt.Println("server close failed: ", err.Error())
		}
		common_log.Info("Exited")
		os.Exit(0)
	}()
}

func irisLogFunc(ctx *context.Context, latency time.Duration) {
	var ip, method, path string

	status := ctx.GetStatusCode()
	method = ctx.Method()
	path = ctx.Request().URL.RequestURI()
	if method == "OPTIONS" {
		return
	}

	line := fmt.Sprintf("%4v %s %s %v %s", latency, ip, method, status, path)
	if context.StatusCodeNotSuccessful(status) {
		body, _ := ctx.GetBody()
		common_log.Error(line, string(body))
		return
	}
}
