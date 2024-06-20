package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"

	"os"

	"backend/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name 编译后的二进制名称
	Name string = "backend"
	// Version 已编译软件的版本。
	Version string = "v1.0.0"
	// flag标记
	flagConf string = "dev"
	// id, _ = os.Hostname()
	id = "dev"
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(
	logger log.Logger,
	rr registry.Registrar, // 注册发现的接口
	gs *grpc.Server,
	hs *http.Server,
	) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(rr),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	var ac conf.Auth
	if err := c.Scan(&ac); err != nil {
		panic(err)
	}

	// 服务注册发现
	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	// trace
	var tc conf.Trace
	if err := c.Scan(&tc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, &ac, &rc, &tc, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
