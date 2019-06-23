package rpc

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"bitbucket.org/mutongx/go-utils/log"

	"bitbucket.org/mutze5/wxfetcher/article"
	"bitbucket.org/mutze5/wxfetcher/db"
	"bitbucket.org/mutze5/wxfetcher/proto"
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
func (s *Server) FetchURL(ctx context.Context, req *proto.FetchURLRequest) (resp *proto.FetchURLResponse, err error) {
	url := req.Url
	// Logs and error handling
	log.Info("RPCServer", "New FetchURL request: %s", url)
	defer func() {
		if err != nil {
			log.Error("RPCServer", "Error in FetchURL(%s): %v", url, err)
		}
	}()
	// Parse article body
	switch {
	case strings.HasPrefix(url, "http://mp.weixin.qq.com/") || strings.HasPrefix(url, "https://mp.weixin.qq.com/"):
		var httpResp *http.Response
		// Fetch article body
		httpResp, err = s.http.Get(url)
		if err != nil {
			return
		}
		defer httpResp.Body.Close()
		// Parse article body
		var atc *article.WxArticle
		atc, err = article.NewFromWxStream(httpResp.Body)
		if err != nil {
			return
		}
		// Fetch short URL key
		var meta *proto.ArticleMeta
		var key string
		meta, key, err = db.ProcessWxArticle(ctx, atc)
		if err != nil {
			return
		}
		return &proto.FetchURLResponse{Key: key, Meta: meta}, nil
	default:
		return nil, errUnsupportedURL
	}
}
