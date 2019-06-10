package main

import (
	"flag"
	"net"
	"os"

	"bitbucket.org/mutongx/go-utils/log"
	"google.golang.org/grpc"

	"bitbucket.org/mutze5/wxfetcher/rpc"
)

func main() {

	// Setup Logger Level
	log.Level(log.Lnotice)

	// Parse flags
	rpcListen := flag.String("rpc-listen", ":9967", "Listen address and port for RPC server")
	webListen := flag.String("web-listen", ":9968", "Listen address and port for web server")
	flag.Parse()

	// Start RPC Server
	log.Notice("Main", "Starting WxFetcher RPC server at %s...", *rpcListen)
	if rpcSocket, err := net.Listen("tcp", *rpcListen); err == nil {
		grpcServer, rpcServer := grpc.NewServer(), rpc.NewServer()
		rpc.RegisterWxFetcherServer(grpcServer, rpcServer)
		go grpcServer.Serve(rpcSocket)
	} else {
		log.Critical("Main", "Error creating RPC socket: %v", err)
		os.Exit(1)
	}

	// Start the Web Server
	log.Notice("Main", "Starting WxFetcher web server at %s...", *webListen)

	select {}

}
