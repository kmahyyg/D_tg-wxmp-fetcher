package rpc

import (
	"context"
	"fmt"
	"net/http"

	"bitbucket.org/mutze5/wxfetcher/article"
)

// Server is a implementation of WxFetcherService
type Server struct {
	http *http.Client
}

// NewServer creates a new Server
func NewServer() *Server {
	return &Server{
		http: &http.Client{},
	}
}

// FetchURL fetches article information from remote and return the URL key
func (s *Server) FetchURL(ctx context.Context, req *FetchURLRequest) (*FetchURLResponse, error) {
	resp, err := s.http.Get(req.OriginalUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(article.NewFromWxStream(resp.Body))
	// TODO: Database Operations
	return &FetchURLResponse{ShortenedKey: "TestKey"}, nil
}
