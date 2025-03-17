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
	trees           map[string]*node // method -> tree root node
	Middlewares     []MiddlewareFunc
	NotFoundMessage *string
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
	debug  bool
}

// New creates a new Engine instance
func New() *Engine {
	router := &Router{
		trees:       make(map[string]*node),
		Middlewares: []MiddlewareFunc{},
	}
	engine := &Engine{
		Router: router,
		groups: []*RouterGroup{
			{
				Router:     router,
				BasePath:   "/",
				Middleware: []MiddlewareFunc{},
			},
		},
		debug: true,
	}
	return engine
}

func (e *Engine) SetDebug(debug bool) {
	e.debug = debug
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

	finalHandler := func(c *Context) {
		// Apply global middlewares
		for _, middleware := range r.Router.Middlewares {
			c.handlers = append(c.handlers, middleware)
		}
		// Apply group-specific middlewares
		for _, middleware := range r.Middleware {
			c.handlers = append(c.handlers, middleware)
		}
		// Append the route handler
		c.handlers = append(c.handlers, handler)
		c.handlerIdx = -1
		c.Next()
	}

	r.Router.addRoute(method, fullPath, finalHandler)
}

func (r *Router) addRoute(method, path string, handler func(*Context)) {
	root := r.trees[method]
	if root == nil {
		root = &node{}
		r.trees[method] = root
	}
	root.insert(path, handler)
}

// Find and execute the appropriate route
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if root := r.trees[req.Method]; root != nil {
		params := make(map[string]string)
		if handler := root.search(req.URL.Path, params); handler != nil {
			c := NewContext(w, req)
			c.params = params
			handler(c)
			return
		}
	}

	if r.NotFoundMessage != nil {
		http.Error(w, *r.NotFoundMessage, http.StatusNotFound)
	} else {
		http.NotFound(w, req)
	}
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
func (e *Engine) Get(path string, handler func(c *Context))     { e.groups[0].Get(path, handler) }
func (e *Engine) Post(path string, handler func(c *Context))    { e.groups[0].Post(path, handler) }
func (e *Engine) Put(path string, handler func(c *Context))     { e.groups[0].Put(path, handler) }
func (e *Engine) Delete(path string, handler func(c *Context))  { e.groups[0].Delete(path, handler) }
func (e *Engine) Patch(path string, handler func(c *Context))   { e.groups[0].Patch(path, handler) }
func (e *Engine) Options(path string, handler func(c *Context)) { e.groups[0].Options(path, handler) }
func (e *Engine) Head(path string, handler func(c *Context))    { e.groups[0].Head(path, handler) }
func (e *Engine) SetNotFoundMessage(message string)             { e.Router.NotFoundMessage = &message }

// RouterGroup HTTP Methods
func (r *RouterGroup) Get(path string, handler func(c *Context))  { r.handle("GET", path, handler) }
func (r *RouterGroup) Post(path string, handler func(c *Context)) { r.handle("POST", path, handler) }
func (r *RouterGroup) Put(path string, handler func(c *Context))  { r.handle("PUT", path, handler) }
func (r *RouterGroup) Delete(path string, handler func(c *Context)) {
	r.handle("DELETE", path, handler)
}
func (r *RouterGroup) Patch(path string, handler func(c *Context)) { r.handle("PATCH", path, handler) }
func (r *RouterGroup) Options(path string, handler func(c *Context)) {
	r.handle("OPTIONS", path, handler)
}
func (r *RouterGroup) Head(path string, handler func(c *Context)) { r.handle("HEAD", path, handler) }

// Engine method to listen on a custom port
func (e *Engine) Listen(port ...string) error {
	address := ResolvePort(port)
	if e.debug {
		e.PrintRoutes(address)
	}
	return http.ListenAndServe(address, e.Router)
}
