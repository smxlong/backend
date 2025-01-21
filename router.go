package backend

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/smxlong/kit/webserver"
)

// Router routes requests to endpoints
type Router struct {
	g  *gin.Engine
	di *Injector
}

// NewRouter creates a new Router
func NewRouter(di *Injector) (*Router, error) {
	g, _, err := di.GetInstance(reflect.TypeOf(&gin.Engine{}))
	if err != nil {
		return nil, err
	}
	return &Router{
		g:  g.(*gin.Engine),
		di: di,
	}, nil
}

// GET routes a GET request to the given path
func (r *Router) GET(path string, handler interface{}) {
	r.g.GET(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// POST routes a POST request to the given path
func (r *Router) POST(path string, handler interface{}) {
	r.g.POST(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// PUT routes a PUT request to the given path
func (r *Router) PUT(path string, handler interface{}) {
	r.g.PUT(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// DELETE routes a DELETE request to the given path
func (r *Router) DELETE(path string, handler interface{}) {
	r.g.DELETE(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// PATCH routes a PATCH request to the given path
func (r *Router) PATCH(path string, handler interface{}) {
	r.g.PATCH(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// OPTIONS routes an OPTIONS request to the given path
func (r *Router) OPTIONS(path string, handler interface{}) {
	r.g.OPTIONS(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// HEAD routes a HEAD request to the given path
func (r *Router) HEAD(path string, handler interface{}) {
	r.g.HEAD(path, func(c *gin.Context) {
		r.di.Invoke(handler, c)
	})
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.g.ServeHTTP(w, req)
}

// RunContext runs the router on the given address until the context is done
func (r *Router) RunContext(ctx context.Context, addr string) error {
	return webserver.ListenAndServe(ctx, &http.Server{
		Addr:    addr,
		Handler: r,
	})
}

// Run runs the router on the given address. Use RunContext if you need to shut
// down the server gracefully.
func (r *Router) Run(addr string) error {
	return r.RunContext(context.Background(), addr)
}
