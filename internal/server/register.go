package server

import (
	"backend/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

func NewRegister(conf *conf.Registry) (registry.Registrar, error) {
	c := &api.Config{
		Address:    conf.Consul.Address,
		Scheme:     conf.Consul.Schema,
		TLSConfig:  api.TLSConfig{},
	}
	cli, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli,consul.WithHealthCheck(conf.Consul.HealthCheck))
	return r, err
}
