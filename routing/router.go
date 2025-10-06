package routing

import (
	"github.com/gofiber/fiber/v2"
)

// Router provides Laravel-inspired routing
type Router struct {
	app *fiber.App
}

// New creates a new Router instance
func New(app *fiber.App) *Router {
	return &Router{app: app}
}

// Get registers a GET route
func (r *Router) Get(path string, handler fiber.Handler) *Router {
	r.app.Get(path, handler)
	return r
}

// Post registers a POST route
func (r *Router) Post(path string, handler fiber.Handler) *Router {
	r.app.Post(path, handler)
	return r
}

// Put registers a PUT route
func (r *Router) Put(path string, handler fiber.Handler) *Router {
	r.app.Put(path, handler)
	return r
}

// Delete registers a DELETE route
func (r *Router) Delete(path string, handler fiber.Handler) *Router {
	r.app.Delete(path, handler)
	return r
}

// Patch registers a PATCH route
func (r *Router) Patch(path string, handler fiber.Handler) *Router {
	r.app.Patch(path, handler)
	return r
}

// Group creates a route group with a prefix
func (r *Router) Group(prefix string, handler func(*RouterGroup)) {
	group := &RouterGroup{
		app:    r.app,
		prefix: prefix,
	}
	handler(group)
}

// RouterGroup represents a group of routes with a common prefix
type RouterGroup struct {
	app    *fiber.App
	prefix string
}

// Get registers a GET route in the group
func (g *RouterGroup) Get(path string, handler fiber.Handler) *RouterGroup {
	g.app.Get(g.prefix+path, handler)
	return g
}

// Post registers a POST route in the group
func (g *RouterGroup) Post(path string, handler fiber.Handler) *RouterGroup {
	g.app.Post(g.prefix+path, handler)
	return g
}

// Put registers a PUT route in the group
func (g *RouterGroup) Put(path string, handler fiber.Handler) *RouterGroup {
	g.app.Put(g.prefix+path, handler)
	return g
}

// Delete registers a DELETE route in the group
func (g *RouterGroup) Delete(path string, handler fiber.Handler) *RouterGroup {
	g.app.Delete(g.prefix+path, handler)
	return g
}

// Patch registers a PATCH route in the group
func (g *RouterGroup) Patch(path string, handler fiber.Handler) *RouterGroup {
	g.app.Patch(g.prefix+path, handler)
	return g
}

// Middleware adds middleware to the group
func (g *RouterGroup) Middleware(middleware ...fiber.Handler) *RouterGroup {
	for _, m := range middleware {
		g.app.Use(g.prefix, m)
	}
	return g
}
