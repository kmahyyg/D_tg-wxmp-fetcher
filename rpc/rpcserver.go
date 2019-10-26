package rpc

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"mutong.moe/go/utils/log"

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
	resp = &proto.FetchURLResponse{}
	// Logs and error handling
	log.Info("RPCServer", "New FetchURL request: %s", req.Url)
	defer func() {
		if err != nil {
			log.Error("RPCServer", "Error in FetchURL(%s): %v", req.Url, err)
			resp.Msg = err.Error()
			err = nil
		}
	}()
	// Parse article body
	switch {
	case strings.HasPrefix(req.Url, "http://mp.weixin.qq.com/") || strings.HasPrefix(req.Url, "https://mp.weixin.qq.com/"):
		var httpResp *http.Response
		// Fetch article body
		if httpResp, err = s.http.Get(req.Url); err != nil {
			resp.Error = proto.FetchURLError_NETWORK
			return
		}
		defer httpResp.Body.Close()
		// Parse article body
		var atc *article.WxArticle
		if atc, err = article.NewFromWxStream(httpResp.Body); err != nil {
			resp.Error = proto.FetchURLError_PARSE
			return
		}
		// Fetch short URL key
		if resp.Meta, resp.Key, err = db.ProcessWxArticle(ctx, atc); err != nil {
			resp.Error = proto.FetchURLError_INTERNAL
			return
		}
	default:
		resp.Error, err = proto.FetchURLError_UNSUPPORTED, errUnsupportedURL
		return
	}
	return
}
