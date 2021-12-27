package main

import (
	"context"
	"net"

	"github.com/buroz/grpc-clean-example/cmd/accounts_service/handlers"
	"github.com/buroz/grpc-clean-example/cmd/accounts_service/interceptors"
	"github.com/buroz/grpc-clean-example/internal/users"
	"github.com/buroz/grpc-clean-example/pkg/config"
	"github.com/buroz/grpc-clean-example/pkg/storage"
	"github.com/buroz/grpc-clean-example/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	conf, err := config.NewConfig()
	utils.Check(err)

	// TODO
	ctx := context.Background()

	var arangoClient = storage.NewArangoClient(&conf.Db)

	err = arangoClient.Connect(ctx)
	utils.Check(err)

	usersRepository, err := users.NewUsersRepository(ctx, arangoClient)
	utils.Check(err)

	usersService := users.NewUsersService(usersRepository)

	usersHandlerImpl := new(handlers.UsersHandler)
	err = usersHandlerImpl.Inject(usersService, &conf.Amqp)
	utils.Check(err)

	usersServiceInterceptors := interceptors.NewServiceInterceptors(usersService)

	tcpServer, err := net.Listen("tcp", ":8000")
	utils.Check(err)

	grpcImpl := grpc.NewServer(
		grpc.UnaryInterceptor(usersServiceInterceptors.AuthorizationInterceptor()),
	)

	users.RegisterUsersServiceServer(grpcImpl, usersHandlerImpl)
	// hello_pb.RegisterHelloServiceServer(grpcImpl, helloHandlerImpl)

	err = grpcImpl.Serve(tcpServer)
	utils.Check(err)
}
