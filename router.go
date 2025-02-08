package rapidgo

import (
	"net/http"
	"regexp"
	"strings"
)

// Route struct to store route information
type Route struct {
	Path       string
	PathParams []string
	PathRegex  *regexp.Regexp
	Method     string
	Handler    func(c *Context)
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

func NewRoute(method, path string, handler func(c *Context)) *Route {
	// Find all :params in the path
	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllString(path, -1)

	// Clean up the param names by removing the colon
	params := make([]string, len(matches))
	for i, match := range matches {
		params[i] = strings.TrimPrefix(match, ":")
	}

	// Replace :params with (\w+)
	pattern := path
	for _, param := range matches {
		pattern = strings.Replace(pattern, param, `(\w+)`, -1)
	}

	regex := regexp.MustCompile(pattern)

	return &Route{
		Path:       path,
		PathParams: params,
		PathRegex:  regex,
		Method:     method,
		Handler:    handler,
	}
}

// Handle requests, applying middleware at the group level
func (r *RouterGroup) handle(method string, path string, handler func(c *Context)) {
	fullPath := r.BasePath + path
	if r.BasePath == "/" {
		fullPath = path
	}

	route := NewRoute(method, fullPath, handler)

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

	route.Handler = finalHandler
	r.Router.Routes = append(r.Router.Routes, route)
}

// Find and execute the appropriate route
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.Routes {
		if route.Method != req.Method {
			continue
		}

		// if route has no params, just compare the path
		if len(route.PathParams) == 0 {
			if route.Path == req.URL.Path {
				route.Handler(NewContext(w, req))
				return
			}
		}

		// if route has params, check if the path matches the regex
		matches := route.PathRegex.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			// create a new context
			c := NewContext(w, req)

			// Add params to the context
			for i, param := range route.PathParams {
				c.params[param] = matches[i+1]
			}
			route.Handler(c)
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
