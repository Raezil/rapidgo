package rapidgo

func SampleMiddleware1() MiddlewareFunc {
	return func(c *Context) {
		c.Set("hello", "world")
		c.Response.Header().Set("Custom-Header-1", "SampleMiddleware 1")
		c.Next()
	}
}

func SampleMiddleware2() MiddlewareFunc {
	return func(c *Context) {
		c.Response.Header().Set("Custom-Header-2", "SampleMiddleware 2")

		c.Next()
	}
}

func SampleMiddleware3() MiddlewareFunc {
	return func(c *Context) {
		c.Response.Header().Set("Custom-Header-3", "SampleMiddleware 3")

		c.Next()
	}
}
