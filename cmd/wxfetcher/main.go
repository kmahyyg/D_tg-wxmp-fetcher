package main

import (
	"context"
	"flag"
	"math/rand"
	"net"
	"os"
	"time"

	"mutong.moe/go/utils/log"
	"google.golang.org/grpc"

	"bitbucket.org/mutze5/wxfetcher/db"
	"bitbucket.org/mutze5/wxfetcher/proto"
	"bitbucket.org/mutze5/wxfetcher/rpc"
	"bitbucket.org/mutze5/wxfetcher/web"
)

func main() {

	// Initialize random seed
	rand.Seed(time.Now().UTC().UnixNano())

	// Parse Flags
	rpcListen := flag.String("rpc-listen", ":9967", "Listen address and port for RPC server")
	webListen := flag.String("web-listen", ":9968", "Listen address and port for web server")
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Read Configuration
	var err error
	var cfg *appConfig
	if cfg, err = readConfig(*configPath); err != nil {
		log.Critical("Main", "Error reading configuration file: %v", err)
		os.Exit(1)
	}

	// Setup Logger Level
	log.Notice("Main", "Switching log level to %s", cfg.LoggingConfig.Level)
	log.LevelByString(cfg.LoggingConfig.Level)

	// Connect to Database
	log.Notice("Main", "Connecting to database...")
	if err := db.Connect(context.Background(), cfg.DBConfig.Driver, cfg.DBConfig.Source); err != nil {
		log.Critical("Main", "Error connecting to databse : %v", err)
		os.Exit(1)
	}

	// Start RPC Server
	log.Notice("Main", "Starting WxFetcher RPC server at %s...", *rpcListen)
	if rpcSocket, err := net.Listen("tcp", *rpcListen); err == nil {
		grpcServer, rpcServer := grpc.NewServer(), rpc.NewServer()
		proto.RegisterWxFetcherServer(grpcServer, rpcServer)
		go grpcServer.Serve(rpcSocket)
	} else {
		log.Critical("Main", "Error creating RPC socket: %v", err)
		os.Exit(1)
	}

	// Start the Web Server
	log.Notice("Main", "Starting WxFetcher web server at %s...", *webListen)
	go web.Serve(*webListen)

	select {}

}
