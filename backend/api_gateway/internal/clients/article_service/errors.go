package article_service

import "errors"

var (
	ErrArticleNotFound    = errors.New("article not found")
	ErrInternalGRPCServer = errors.New("internal gRPC server error")
)
