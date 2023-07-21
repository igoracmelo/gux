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

type Mux struct {
	routes []route
}

var _ http.Handler = (*Mux)(nil)

// TODO: automatic OPTIONS
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound

	for _, route := range m.routes {
		vars, ok := parseUrlVars(r.URL.Path, route.path)
		if !ok {
			continue
		}
		if r.Method != route.method {
			status = http.StatusMethodNotAllowed
			continue
		}
		ctx := &Ctx{
			W:    w,
			R:    r,
			Vars: vars,
			Data: map[string]any{},
		}
		route.fn(ctx)
		return
	}

	w.WriteHeader(status)
}

func parseUrlVars(got string, pattern string) (map[string]string, bool) {
	vars := make(map[string]string)
	if got == pattern {
		return vars, true
	}

	gotChunks := strings.Split(got, "/")
	patternChunks := strings.Split(pattern, "/")

	if len(gotChunks) != len(patternChunks) {
		return nil, false
	}

	for i := 0; i < len(gotChunks); i++ {
		value := gotChunks[i]
		patternChunk := patternChunks[i]

		if value == patternChunk {
			continue
		}
		if patternChunk[0] != ':' {
			return nil, false
		}

		varName := strings.TrimPrefix(patternChunk, ":")
		vars[varName] = value
	}

	return vars, true
}

func (m *Mux) Handle(method, path string, fn HandlerFunc) {
	m.routes = append(m.routes, route{
		method,
		path,
		fn,
	})
}

// Convenient functions

func (m *Mux) Head(path string, fn HandlerFunc) {
	m.Handle("HEAD", path, fn)
}

func (m *Mux) Get(path string, fn HandlerFunc) {
	m.Handle("GET", path, fn)
}

func (m *Mux) Post(path string, fn HandlerFunc) {
	m.Handle("POST", path, fn)
}

func (m *Mux) Patch(path string, fn HandlerFunc) {
	m.Handle("PATCH", path, fn)
}

func (m *Mux) Put(path string, fn HandlerFunc) {
	m.Handle("PUT", path, fn)
}

func (m *Mux) Delete(path string, fn HandlerFunc) {
	m.Handle("DELETE", path, fn)
}
