package jwt

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	http2 "net/http"

	"log"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type server struct {
	UnimplementedGreeterServer

	hc GreeterClient
}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "hello from service"}, nil
}

func TestJwtIntegration(t *testing.T) {
	// Set up server
	var app *kratos.App
	go func() {
		testKey := "testKey"
		httpSrv := http.NewServer(
			http.Address(":8000"),
			http.Middleware(
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(testKey), nil
				}),
			),
		)
		grpcSrv := grpc.NewServer(
			grpc.Address(":9000"),
			grpc.Middleware(
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(testKey), nil
				}),
			),
		)
		serviceTestKey := "serviceTestKey"
		con, _ := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint("dns:///127.0.0.1:9001"),
			grpc.WithMiddleware(
				jwt.Client(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(serviceTestKey), nil
				}),
			),
		)
		s := &server{
			hc: NewGreeterClient(con),
		}
		RegisterGreeterServer(grpcSrv, s)
		RegisterGreeterHTTPServer(httpSrv, s)
		app = kratos.New(
			kratos.Name("helloworld"),
			kratos.Server(
				httpSrv,
				grpcSrv,
			),
		)
		// 这里插入之前服务端的启动代码
		if err := app.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	// 确保服务有足够时间启动
	time.Sleep(time.Second * 5)

	// 生成 JWT token 示例
	testToken := createTestToken("testKey", jwtv5.MapClaims{"user": "testUser"})

	// 测试 HTTP 和 gRPC 客户端
	go testHTTPClient(t, "localhost:8000", testToken)
	// go testGRPCClient(t, "localhost:9000", testToken)

	// 设置超时以防万一
	select {
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out")
	case <-func() chan struct{} {
		// 当测试成功运行后，关闭 app
		defer app.Stop()
		ch := make(chan struct{})
		go func() {
			t.Log("测试成功")
			close(ch)
		}()
		return ch
	}():
	}
}

// Helper function to create a JWT token
func createTestToken(secretKey string, claims jwtv5.Claims) string {
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}
	return signedToken
}

func testHTTPClient(t *testing.T, address, token string) {
	req, err := http2.NewRequest("GET", "http://"+address+"/helloworld/YOUR_USER_NAME", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http2.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http2.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
