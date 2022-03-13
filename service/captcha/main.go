package main

import (
	"captcha/handler"
	pb "captcha/proto"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	service "github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
)

func main() {
	consulReg := consul.NewRegistry()

	// Create service
	srv := service.NewService(
		service.Address("127.0.0.1:12341"),
		service.Name("captcha"),
		service.Version("latest"),
		service.Registry(consulReg),
	)

	// Register handler
	pb.RegisterCaptchaHandler(srv.Server(), new(handler.Captcha))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
