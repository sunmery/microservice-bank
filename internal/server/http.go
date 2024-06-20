package server

import (
	v1 "backend/api/helloworld/v1"
	"backend/internal/conf"
	"backend/internal/service"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/gorilla/handlers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewWhiteListMatcher 设置白名单，不需要 token 验证的接口
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/shop.v1.ShopService/Captcha"] = struct{}{}
	whiteList["/shop.v1.ShopService/Login"] = struct{}{}
	whiteList["/shop.v1.ShopService/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// Initializes an OTLP exporter, and configures the corresponding trace provider.
func initTracerProvider(ctx context.Context, res *resource.Resource, conn string) (func(context.Context) error, error) {
	// Set up a trace exporter

	// 服务端的Jaeger支持HTTPS时使用
	// traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(conn))

	// 服务端的Jaeger不支持HTTPS时使用otlptracehttp.WithInsecure()显式声明只使用HTTP不安全的连接
	traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// 将基于父span的采样率设置
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		sdktrace.WithBatcher(traceExporter),
		// 在资源中记录有关此应用程序的信息
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	ac *conf.Auth,
	c *conf.Server,
	tr *conf.Trace,
	greeter *service.GreeterService,
	logger log.Logger,
	) *http.Server {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			// serviceName,
			semconv.ServiceNameKey.String(tr.Jaeger.ServiceName),
			attribute.String("exporter", "otlptracehttp"),
			attribute.String("environment", "dev"),
			attribute.Float64("float", 312.23),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// shutdownTracerProvider, err := initTracerProvider(ctx, res, tr.Jaeger.Http.Endpoint)
	_, err2 := initTracerProvider(ctx, res, tr.Jaeger.Http.Endpoint)
	if err2 != nil {
		log.Fatal(err)
	}

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger), // logging 日志
			tracing.Server(), // trace 链路追踪
			// jwt 身份验证
			selector.Server(
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(ac.JwtKey), nil
				}),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
		http.Filter(handlers.CORS( // 浏览器跨域
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
