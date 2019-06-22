package rpc

import (
	"context"
	"net/http"

	"bitbucket.org/mutongx/go-utils/log"

	"bitbucket.org/mutze5/wxfetcher/article"
	"bitbucket.org/mutze5/wxfetcher/db"
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
func (s *Server) FetchURL(ctx context.Context, req *FetchURLRequest) (resp *FetchURLResponse, err error) {
	log.Info("RPCServer", "New FetchURL request: %s", req.OriginalUrl)
	defer func() {
		if err != nil {
			log.Error("RPCServer", "Error in FetchURL(%s): %v", req.OriginalUrl, err)
		}
	}()
	// Fetch article body
	httpResp, err := s.http.Get(req.OriginalUrl)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	// Parse article body (Currently is WeChat only)
	atc, err := article.NewFromWxStream(httpResp.Body)
	if err != nil {
		return
	}
	key, err := db.GetWxArticleKey(ctx, atc)
	if err != nil {
		return
	}
	return &FetchURLResponse{ShortenedKey: key}, nil
}
