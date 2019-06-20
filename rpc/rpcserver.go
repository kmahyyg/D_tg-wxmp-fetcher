package rpc

import (
	"context"
	"net/http"

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
func (s *Server) FetchURL(ctx context.Context, req *FetchURLRequest) (*FetchURLResponse, error) {
	resp, err := s.http.Get(req.OriginalUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	atc, err := article.NewFromWxStream(resp.Body)
	if err != nil {
		return nil, err
	}
	key, err := db.GetWxArticleKey(ctx, atc)
	if err != nil {
		return nil, err
	}
	return &FetchURLResponse{ShortenedKey: key}, nil
}
