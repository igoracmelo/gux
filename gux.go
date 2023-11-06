package gux

import (
	"net/http"
	"strings"
)

type Ctx struct {
	W    http.ResponseWriter
	R    *http.Request
	Vars map[string]string
	Data map[string]any
}

type HandlerFunc func(c *Ctx)

type route struct {
	method string
	path   string
	fn     HandlerFunc
}

type mux struct {
	routes   []route
	NotFound HandlerFunc
}

func New() *mux {
	return &mux{
		NotFound: func(c *Ctx) {
			http.Error(c.W, c.R.URL.Path+" not found", http.StatusNotFound)
		},
	}
}

var _ http.Handler = (*mux)(nil)

func (m *mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defaultHandler := m.NotFound
	ctx := &Ctx{
		W:    w,
		R:    req,
		Data: map[string]any{},
	}

	for _, route := range m.routes {
		vars, ok := parseUrlVars(req.URL.Path, route.path)
		if !ok {
			continue
		}
		ctx.Vars = vars

		if req.Method == route.method || req.Method == "HEAD" && route.method == "GET" {
			route.fn(ctx)
			return
		} else {
			defaultHandler = func(c *Ctx) {
				c.W.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	}

	defaultHandler(ctx)
}

func parseUrlVars(got string, pattern string) (map[string]string, bool) {
	vars := make(map[string]string)
	if got == pattern {
		return vars, true
	}

	gotChunks := strings.Split(got, "/")
	patternChunks := strings.Split(pattern, "/")

	if len(gotChunks) != len(patternChunks) {
		return vars, false
	}

	for i := 0; i < len(gotChunks); i++ {
		value := gotChunks[i]
		patternChunk := patternChunks[i]

		if value == patternChunk {
			continue
		}
		if patternChunk[0] != ':' {
			return vars, false
		}

		varName := strings.TrimPrefix(patternChunk, ":")
		vars[varName] = value
	}

	return vars, true
}

func (m *mux) Handle(method, path string, fn HandlerFunc) {
	m.routes = append(m.routes, route{
		method,
		path,
		fn,
	})
}

// Convenient functions

func (m *mux) Head(path string, fn HandlerFunc) {
	m.Handle("HEAD", path, fn)
}

func (m *mux) Get(path string, fn HandlerFunc) {
	m.Handle("GET", path, fn)
}

func (m *mux) Post(path string, fn HandlerFunc) {
	m.Handle("POST", path, fn)
}

func (m *mux) Patch(path string, fn HandlerFunc) {
	m.Handle("PATCH", path, fn)
}

func (m *mux) Put(path string, fn HandlerFunc) {
	m.Handle("PUT", path, fn)
}

func (m *mux) Delete(path string, fn HandlerFunc) {
	m.Handle("DELETE", path, fn)
}
