package service

import (
	v1 "backend/api/helloworld/v1"
	"backend/internal/biz"
	"backend/internal/helper/log"
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.

func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	tr := otel.Tracer("full-stack-engineering-backend-api")
	spCtx, span := tr.Start(ctx, "func-a")
	span.SetAttributes(attribute.String("name", "funA"))
	type _LogStruct struct {
		CurrentTime time.Time `json:"currentTime"`
		PassWho     string    `json:"passWho"`
		Name        string    `json:"name"`
	}
	logTest := _LogStruct{
		CurrentTime: time.Time{},
		PassWho:     "jzin",
		Name:        "func-a",
	}
	log.InfofC(spCtx, "is logs")
	b, _ := json.Marshal(logTest)
	log.InfofC(spCtx, string(b))
	span.SetAttributes(attribute.Key("测试key").String(string(b)))
	time.Sleep(time.Second)

	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	span.End()
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
