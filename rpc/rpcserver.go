package rpc

import (
	"context"
	"errors"
	"net/http"

	"bitbucket.org/mutongx/go-utils/log"

	"bitbucket.org/mutze5/wxfetcher/article"
	"bitbucket.org/mutze5/wxfetcher/db"
)

var (
	errUnsupportedURL = errors.New("unsupported url")
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
	url := req.OriginalUrl
	// Logs and error handling
	log.Info("RPCServer", "New FetchURL request: %s", url)
	defer func() {
		if err != nil {
			log.Error("RPCServer", "Error in FetchURL(%s): %v", url, err)
		}
	}()
	// Fetch article body
	httpResp, err := s.http.Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	// Parse article body
	switch {
	case strings.HasPrefix(url, "http://mp.weixin.qq.com") || strings.HasPrefix(url, "https://mp.weixin.qq.com"):
		atc, err := article.NewFromWxStream(httpResp.Body)
		if err != nil {
			return
		}
		key, err := db.GetWxArticleKey(ctx, atc)
		if err != nil {
			return
		}
		return &FetchURLResponse{ShortenedKey: key}, nil
	default:
		return nil, errUnsupportedURL
	}
}
