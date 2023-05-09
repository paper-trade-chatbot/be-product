package main

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/paper-trade-chatbot/be-common/cache"
	"github.com/paper-trade-chatbot/be-common/database"
	"github.com/paper-trade-chatbot/be-product/service/product"
	productGrpc "github.com/paper-trade-chatbot/be-proto/product"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/paper-trade-chatbot/be-common/config"
	"github.com/paper-trade-chatbot/be-common/global"
	"github.com/paper-trade-chatbot/be-common/logging"
	"github.com/paper-trade-chatbot/be-common/server"
)

func main() {

	global.Alive = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging.Initialize(ctx)
	defer logging.Finalize()

	cache.Initialize(ctx)
	defer cache.Finalize()

	database.Initialize(ctx)
	defer database.Finalize()

	initConfig()

	grpcAddress := fmt.Sprintf("%s:%s",
		config.GetString("GRPC_SERVER_LISTEN_ADDRESS"),
		config.GetString("GRPC_SERVER_LISTEN_PORT"))
	listener := server.CreateGRpcServer(ctx, grpcAddress)

	recoveryOpt := grpc_recovery.WithRecoveryHandlerContext(
		func(ctx context.Context, p interface{}) error {
			logging.Error(ctx, "[PANIC] %s\n\n%s", p, string(debug.Stack()))
			return status.Errorf(codes.Internal, "%s", p)
		},
	)
	grpc := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(recoveryOpt),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recoveryOpt),
		)),
	)
	reflection.Register(grpc)

	productInstance := product.New()
	productGrpc.RegisterProductServiceServer(grpc, productInstance)

	address := fmt.Sprintf("%s:%s",
		config.GetString("SERVER_LISTEN_ADDRESS"),
		config.GetString("SERVER_LISTEN_PORT"))
	httpServer := server.CreateHttpServer(ctx, address)

	go func() {
		logging.Info(ctx, "grpc serving")
		if err := grpc.Serve(*listener); err != nil {
			panic(err)
		}
	}()

	global.Ready = true

	logging.Info(ctx, "Initialization complete, listening on %s...", address)
	if err := httpServer.ListenAndServe(); err != nil {
		logging.Info(ctx, err.Error())
	}

}

func initConfig() {
	global.Initialize()

}
