package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/Matltin/simple-bank/api"
	db "github.com/Matltin/simple-bank/db/sqlc"
	"github.com/Matltin/simple-bank/gapi"
	"github.com/Matltin/simple-bank/mail"
	"github.com/Matltin/simple-bank/pb"
	"github.com/Matltin/simple-bank/util"
	"github.com/Matltin/simple-bank/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config: %s", err)
		return
	}

	if config.ENVIROMENT == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %s ", err)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)

	runGateWayServer(store, config, taskDistributor)
	// runGrpcServer(store, config)
	// runGinServer(store, config)
}

func runGrpcServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %v", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal().Msgf("cannot start gRPC server: %s", err)
	}
}

func runGateWayServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %v", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msgf("cannot register handle server: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener")
	}

	log.Printf("start HTTP gatewaay server at %s", lis.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(lis, handler)
	if err != nil {
		log.Fatal().Msgf("cannot start HTTP gateway server")
	}
}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot start server: %s", err)
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mail := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAdderss, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mail)
	err := taskProcessor.Start()
	if err != nil {
		fmt.Println("failed to start task processor")
	}
}
