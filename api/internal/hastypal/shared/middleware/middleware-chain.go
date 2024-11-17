package middleware

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/http"
)

type MiddlewareChain struct {
	middlewares []types.Middleware
}

func NewMiddlewareChain(initialMiddlewares ...types.Middleware) *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: initialMiddlewares,
	}
}

func (c *MiddlewareChain) ApplyOn(next http.Handler) http.Handler {
	for i := 0; i < len(c.middlewares); i++ {
		next = c.middlewares[i](next)
	}

	return next
}
