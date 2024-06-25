package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/gorilla/handlers"
	accountV1 "user/api/account/v1"
	"user/internal/conf"
	"user/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// NewWhiteListMatcher 设置白名单，不需要 token 验证的接口
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.account.v1.AccountService/GetAccount"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	ac *conf.Auth,
	c *conf.Server,
	account *service.AccountService,
	logger log.Logger,
	) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
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
			handlers.AllowedOrigins([]string{"http://localhost:3000", "http://127.0.0.1:3000"}),
			// handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowCredentials(),
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
	accountV1.RegisterAccountServiceHTTPServer(srv, account)
	return srv
}
