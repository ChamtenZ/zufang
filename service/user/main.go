package main

import (
	"user/handler"
	"user/model"
	pb "user/proto"

	// "github.com/micro/micro/v3/service"
	// "github.com/micro/micro/v3/service/logger"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	service "github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
)

func main() {
	//初始化连接池
	model.InitRedis()
	model.InitDb()

	consulReg := consul.NewRegistry()

	// Create service
	srv := service.NewService(
		service.Address("127.0.0.1:12342"),
		service.Registry(consulReg),
		service.Name("user"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterUserHandler(srv.Server(), new(handler.User))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
