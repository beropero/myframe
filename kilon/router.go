package kilon

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	Handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		make(map[string]*node),
		make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.Handlers[key] = handler
}

func (r *router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	ctx.Params = params
	if n != nil {
		key := ctx.Method + "-" + n.pattern
		ctx.handlers = append(ctx.handlers, r.Handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	ctx.Next()
}

func parsePattern(pattern string) []string {
	temp := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range temp {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}
