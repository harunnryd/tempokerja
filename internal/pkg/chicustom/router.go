package chicustom

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router interface {
	chi.Router
	Action(...Rest)
}

type router struct {
	chi.Router
	PrepareHandler Handler
}

func NewRouter(chirouter chi.Router) Router {
	return &router{
		Router: chirouter,
	}
}

var _ chi.Router = &router{}

func (r *router) copy(chirouter chi.Router) Router {
	return &router{
		Router:         chirouter,
		PrepareHandler: r.PrepareHandler,
	}
}

// With adds inline middlewares for an endpoint handler
func (r *router) With(middlewares ...func(http.Handler) http.Handler) chi.Router {
	return r.copy(r.Router.With(middlewares...))
}

// Group adds a new inline-Router along the current routing
// path, with a fresh middleware stack for the inline-Router
func (r *router) Group(fn func(r chi.Router)) chi.Router {
	im := r.copy(r.With())
	if fn != nil {
		fn(im)
	}
	return im
}

// Route mounts a sub-Router along a `pattern`` string.
func (r *router) Route(pattern string, fn func(r chi.Router)) chi.Router {
	subRouter := r.copy(chi.NewRouter())
	if fn != nil {
		fn(subRouter)
	}
	r.Mount(pattern, subRouter)
	return subRouter
}

// Mount attaches another http.Handler along ./pattern/*
func (r *router) Mount(pattern string, handler http.Handler) {
	r.Router.Mount(pattern, handler)
}

// Handle adds routes for `pattern` that matches all HTTP methods
func (r *router) Handle(pattern string, handler http.Handler) {
	r.Router.Handle(pattern, handler)
}

// Method adds routes for `pattern` that matches the `method` HTTP method
func (r *router) Method(method, pattern string, handler http.Handler) {
	r.Router.Method(method, pattern, handler)
}

// Action adds one or more HTTPAction for `handler.Pattern()` that matches the `handler.HTTPMethod()` HTTP method
func (r *router) Action(handlers ...Rest) {
	for _, handler := range handlers {
		r.Router.Method(handler.GetHTTPMethod(), handler.GetPattern(), handler.GetHandler())
	}
}
