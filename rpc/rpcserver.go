package rpc

import (
	"context"
	"net/http"

	"bitbucket.org/mutze5/wxfetcher/parse"
)

type RPCServer struct {
	http *http.Client
}

func NewRPCServer() *RPCServer {
	return &RPCServer{
		http: &http.Client{},
	}
}

func (s *RPCServer) FetchURL(ctx context.Context, req *FetchURLRequest) (*FetchURLResponse, error) {
	resp, err := s.http.Get(req.OriginalUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	parse.Consume(resp.Body)
	// TODO: Database Operations
	return &FetchURLResponse{ShortenedKey: "TestKey"}, nil
}
