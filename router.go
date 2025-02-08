package rapidgo

import (
	"net/http"
)

// Route struct to store route information
type Route struct {
	Path    string
	Method  string
	Handler func(c *Context)
}

// Middleware function signature
type MiddlewareFunc func(c *Context)

// Router to manage routes and global middlewares
type Router struct {
	Routes      []*Route
	Middlewares []MiddlewareFunc // Store global middleware
}

// RouterGroup for grouping routes with specific base paths and middleware
type RouterGroup struct {
	Router     *Router
	BasePath   string
	Middleware []MiddlewareFunc
}

// Engine is the main struct managing router groups
type Engine struct {
	Router *Router
	groups []*RouterGroup
}

// New creates a new Engine instance
func New() *Engine {
	router := &Router{Routes: []*Route{}, Middlewares: []MiddlewareFunc{}}
	engine := &Engine{
		Router: router,
		groups: []*RouterGroup{
			{
				Router:     router,
				BasePath:   "/",
				Middleware: []MiddlewareFunc{},
			},
		},
	}
	return engine
}

// Handle requests, applying middleware at the group level
func (r *RouterGroup) handle(method string, path string, handler func(c *Context)) {
	fullPath := r.BasePath + path
	if r.BasePath == "/" {
		fullPath = path
	}

	finalHandler := func(c *Context) {
		for _, middleware := range r.Router.Middlewares {
			c.handlers = append(c.handlers, middleware)
		}
		for _, middleware := range r.Middleware {
			c.handlers = append(c.handlers, middleware)
		}
		c.handlers = append(c.handlers, handler)
		c.handlerIdx = -1
		c.Next()
	}

	r.Router.Routes = append(r.Router.Routes, &Route{
		Path:    fullPath,
		Method:  method,
		Handler: finalHandler,
	})
}

// Find and execute the appropriate route
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.Routes {
		if route.Method == req.Method && route.Path == req.URL.Path {
			route.Handler(NewContext(w, req))
			return
		}
	}
	http.NotFound(w, req)
}

// Grouping Routes
func (e *Engine) Group(path string) *RouterGroup {
	group := &RouterGroup{
		Router:     e.Router, // Use the shared router instance
		BasePath:   path,
		Middleware: []MiddlewareFunc{},
	}
	e.groups = append(e.groups, group)
	return group
}

// Use middleware at group level
func (r *RouterGroup) Use(middleware ...MiddlewareFunc) {
	r.Middleware = append(r.Middleware, middleware...)
}

// Use middleware at global level
func (e *Engine) Use(middleware ...MiddlewareFunc) {
	e.Router.Middlewares = append(e.Router.Middlewares, middleware...)
}

// Engine HTTP Methods
func (e *Engine) GET(path string, handler func(c *Context))     { e.groups[0].GET(path, handler) }
func (e *Engine) POST(path string, handler func(c *Context))    { e.groups[0].POST(path, handler) }
func (e *Engine) PUT(path string, handler func(c *Context))     { e.groups[0].PUT(path, handler) }
func (e *Engine) DELETE(path string, handler func(c *Context))  { e.groups[0].DELETE(path, handler) }
func (e *Engine) PATCH(path string, handler func(c *Context))   { e.groups[0].PATCH(path, handler) }
func (e *Engine) OPTIONS(path string, handler func(c *Context)) { e.groups[0].OPTIONS(path, handler) }
func (e *Engine) HEAD(path string, handler func(c *Context))    { e.groups[0].HEAD(path, handler) }

// RouterGroup HTTP Methods
func (r *RouterGroup) GET(path string, handler func(c *Context))  { r.handle("GET", path, handler) }
func (r *RouterGroup) POST(path string, handler func(c *Context)) { r.handle("POST", path, handler) }
func (r *RouterGroup) PUT(path string, handler func(c *Context))  { r.handle("PUT", path, handler) }
func (r *RouterGroup) DELETE(path string, handler func(c *Context)) {
	r.handle("DELETE", path, handler)
}
func (r *RouterGroup) PATCH(path string, handler func(c *Context)) { r.handle("PATCH", path, handler) }
func (r *RouterGroup) OPTIONS(path string, handler func(c *Context)) {
	r.handle("OPTIONS", path, handler)
}
func (r *RouterGroup) HEAD(path string, handler func(c *Context)) { r.handle("HEAD", path, handler) }

// Engine method to listen on a custom port
func (e *Engine) Listen() error {
	return http.ListenAndServe(DEFAULT_PORT, e.Router)
}
